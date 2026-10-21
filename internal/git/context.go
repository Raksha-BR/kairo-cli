// internal/git/context.go
package git

import (
	"os"
	gogit "github.com/go-git/go-git/v5"
)

// Context holds the current Git environment
type Context struct {
	RepoPath   string
	BranchName string
}

// GetCurrentContext returns the git repo root and active branch.
func GetCurrentContext() Context {
	// 1. Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return Context{}
	}

	// 2. Look for a .git folder in this directory or parents
	repo, err := gogit.PlainOpenWithOptions(dir, &gogit.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return Context{} // Not in a git repository
	}

	// 3. Get the active branch name
	head, err := repo.Head()
	branchName := ""
	if err == nil {
		branchName = head.Name().Short()
	}

	// 4. Get the absolute path of the repository root
	wt, err := repo.Worktree()
	repoPath := dir
	if err == nil {
		repoPath = wt.Filesystem.Root()
	}

	return Context{
		RepoPath:   repoPath,
		BranchName: branchName,
	}
}