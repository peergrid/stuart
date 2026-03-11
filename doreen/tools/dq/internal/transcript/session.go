package transcript

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const claudeProjectsDir = ".claude/projects"

// ProjectDir returns the path to a project's transcript directory.
func ProjectDir(projectName string) string {
	home, _ := os.UserHomeDir()
	slug := "-home-" + os.Getenv("USER") + "-" + projectName
	return filepath.Join(home, claudeProjectsDir, slug)
}

// FindSessionFiles returns top-level session JSONL files sorted by mtime.
func FindSessionFiles(projDir string) ([]string, error) {
	entries, err := os.ReadDir(projDir)
	if err != nil {
		return nil, err
	}

	type fileInfo struct {
		path  string
		mtime int64
	}

	var files []fileInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		path := filepath.Join(projDir, e.Name())
		info, err := e.Info()
		if err != nil {
			continue
		}
		files = append(files, fileInfo{path, info.ModTime().Unix()})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].mtime < files[j].mtime
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

// LatestSession returns the most recently modified session file in a project.
func LatestSession(projDir string) (string, error) {
	files, err := FindSessionFiles(projDir)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", os.ErrNotExist
	}
	return files[len(files)-1], nil
}
