package transcript

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const claudeProjectsDir = ".claude/projects"

// ProjectDirFromCWD detects the project transcript directory from the current
// working directory. Claude Code stores transcripts under ~/.claude/projects/
// using a slug derived from the absolute project path (slashes replaced with
// dashes, leading dash).
func ProjectDirFromCWD() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot determine CWD: %w", err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}

	projectsBase := filepath.Join(home, claudeProjectsDir)

	// Try CWD and each parent directory until we find a matching project slug.
	dir := cwd
	for {
		slug := pathToSlug(dir)
		candidate := filepath.Join(projectsBase, slug)
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("no Claude project found for CWD %s", cwd)
}

// ProjectDirFromName returns the transcript directory for a named project.
func ProjectDirFromName(name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	projectsBase := filepath.Join(home, claudeProjectsDir)
	entries, err := os.ReadDir(projectsBase)
	if err != nil {
		return "", err
	}
	// Match by suffix — the slug ends with the project name
	for _, e := range entries {
		if e.IsDir() && strings.HasSuffix(e.Name(), "-"+name) {
			return filepath.Join(projectsBase, e.Name()), nil
		}
	}
	// Try exact match
	candidate := filepath.Join(projectsBase, name)
	if info, err := os.Stat(candidate); err == nil && info.IsDir() {
		return candidate, nil
	}
	return "", fmt.Errorf("project %q not found in %s", name, projectsBase)
}

// pathToSlug converts an absolute path to the Claude projects slug format.
// /home/user/myproject -> -home-user-myproject
func pathToSlug(absPath string) string {
	return strings.ReplaceAll(absPath, "/", "-")
}

// FindSessionsInWindow returns all session JSONL files modified within the
// given time window, sorted by modification time (oldest first).
func FindSessionsInWindow(projDir string, since, until time.Time) ([]string, error) {
	entries, err := os.ReadDir(projDir)
	if err != nil {
		return nil, err
	}

	type fileInfo struct {
		path  string
		mtime time.Time
	}

	var files []fileInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		mtime := info.ModTime()
		if mtime.Before(since) || mtime.After(until) {
			continue
		}
		files = append(files, fileInfo{
			path:  filepath.Join(projDir, e.Name()),
			mtime: mtime,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].mtime.Before(files[j].mtime)
	})

	result := make([]string, len(files))
	for i, f := range files {
		result[i] = f.path
	}
	return result, nil
}

// FindSubagentFiles returns subagent JSONL files for a given session.
func FindSubagentFiles(projDir, sessionID string) ([]string, error) {
	subagentDir := filepath.Join(projDir, sessionID, "subagents")
	entries, err := os.ReadDir(subagentDir)
	if err != nil {
		// Try prefix match
		parentEntries, err2 := os.ReadDir(projDir)
		if err2 != nil {
			return nil, err
		}
		for _, e := range parentEntries {
			if e.IsDir() && strings.HasPrefix(e.Name(), sessionID) {
				candidate := filepath.Join(projDir, e.Name(), "subagents")
				entries, err = os.ReadDir(candidate)
				if err == nil {
					subagentDir = candidate
					break
				}
			}
		}
		if err != nil {
			return nil, nil // No subagents found
		}
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".jsonl") {
			files = append(files, filepath.Join(subagentDir, e.Name()))
		}
	}
	sort.Strings(files)
	return files, nil
}

// ParseDuration parses human-friendly duration strings: 30m, 2h, 12h, 1d, 3d, 1w.
// Also accepts ISO dates (2026-03-10) and ISO timestamps.
func ParseDuration(s string) (time.Time, error) {
	// Try as duration shorthand: Nd, Nh, Nm, Nw
	re := regexp.MustCompile(`^(\d+)([mhdw])$`)
	if m := re.FindStringSubmatch(s); m != nil {
		n, _ := strconv.Atoi(m[1])
		now := time.Now()
		switch m[2] {
		case "m":
			return now.Add(-time.Duration(n) * time.Minute), nil
		case "h":
			return now.Add(-time.Duration(n) * time.Hour), nil
		case "d":
			return now.AddDate(0, 0, -n), nil
		case "w":
			return now.AddDate(0, 0, -n*7), nil
		}
	}

	// Try as ISO date
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}

	// Try as ISO timestamp
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("cannot parse time %q (use: 30m, 2h, 3d, 1w, or ISO date)", s)
}
