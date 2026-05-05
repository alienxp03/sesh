package namer

import (
	"path/filepath"
	"strings"
)

func getGitTopLevelPath(n *RealNamer, path string) (bool, string, error) {
	isGit, topLevel, err := n.git.ShowTopLevel(path)
	if err != nil {
		return false, "", nil
	}
	if !isGit || topLevel == "" {
		return false, "", nil
	}
	return true, topLevel, nil
}

func isOutsideRoot(relativePath string) bool {
	if relativePath == ".." {
		return true
	}
	if strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
		return true
	}
	return filepath.IsAbs(relativePath)
}

func gitRelativeName(n *RealNamer, rootPath string, targetPath string) (string, bool) {
	relativePath, err := filepath.Rel(rootPath, targetPath)
	if err != nil {
		return "", false
	}
	if isOutsideRoot(relativePath) {
		return "", false
	}
	repoName := n.pathwrap.Base(rootPath)
	if relativePath == "." {
		return repoName, true
	}
	return repoName + "/" + relativePath, true
}

// Names a session as <worktreeName>/<relativePath> where worktreeName is the
// basename of the current worktree and relativePath is the input path's offset
// from it.
func gitName(n *RealNamer, path string) (string, error) {
	isGit, topLevel, err := getGitTopLevelPath(n, path)
	if err != nil {
		return "", err
	}
	if !isGit {
		return "", nil
	}
	if name, ok := gitRelativeName(n, topLevel, path); ok {
		return name, nil
	}
	return dirName(n, path)
}

// Names the session root as the basename of the current worktree. This collapses
// nested subdirectories to their containing worktree.
func gitRootName(n *RealNamer, path string) (string, error) {
	isGit, topLevel, err := getGitTopLevelPath(n, path)
	if err != nil {
		return "", err
	}
	if !isGit {
		return "", nil
	}
	if name, ok := gitRelativeName(n, topLevel, topLevel); ok {
		return name, nil
	}
	return dirName(n, topLevel)
}
