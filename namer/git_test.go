package namer

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/git"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
)

func newNamer(t *testing.T) (*RealNamer, *pathwrap.MockPath, *git.MockGit) {
	t.Helper()
	mockPathwrap := new(pathwrap.MockPath)
	mockGit := new(git.MockGit)
	mockHome := new(home.MockHome)
	config := model.Config{DirLength: 1}
	return &RealNamer{pathwrap: mockPathwrap, git: mockGit, home: mockHome, config: config}, mockPathwrap, mockGit
}

func TestGitName(t *testing.T) {
	t.Run("regular clone at main tree root", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/nu"
		mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/nu", nil)
		mp.On("Base", "/Users/hansolo/code/project/nu").Return("nu")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "nu", name)
	})

	t.Run("regular clone, nested subdir", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/nu/server"
		mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/nu", nil)
		mp.On("Base", "/Users/hansolo/code/project/nu").Return("nu")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "nu/server", name)
	})

	t.Run("regular clone, linked worktree", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/nu/.wk/5969"
		mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/nu/.wk/5969", nil)
		mp.On("Base", "/Users/hansolo/code/project/nu/.wk/5969").Return("5969")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "5969", name)
	})

	t.Run("bare repo, .bare suffix", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/sesh/main"
		mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/sesh/main", nil)
		mp.On("Base", "/Users/hansolo/code/project/sesh/main").Return("main")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "main", name)
	})

	t.Run("bare repo, no suffix", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/sesh/main"
		mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/sesh/main", nil)
		mp.On("Base", "/Users/hansolo/code/project/sesh/main").Return("main")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "main", name)
	})

	t.Run("path with spaces", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/alice/My Projects/cool repo/src"
		mg.On("ShowTopLevel", path).Return(true, "/Users/alice/My Projects/cool repo", nil)
		mp.On("Base", "/Users/alice/My Projects/cool repo").Return("cool repo")

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "cool repo/src", name)
	})

	t.Run("non-git directory, ShowTopLevel returns false", func(t *testing.T) {
		n, _, mg := newNamer(t)
		path := "/Users/hansolo/.config/nvim"
		mg.On("ShowTopLevel", path).Return(false, "", nil)

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "", name)
	})

	t.Run("non-git directory, ShowTopLevel errors", func(t *testing.T) {
		n, _, mg := newNamer(t)
		path := "/Users/hansolo/.config/nvim"
		mg.On("ShowTopLevel", path).Return(false, "", fmt.Errorf("not a git repository"))

		name, err := gitName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "", name)
	})
}

func TestGitNameLinkedWorktreeOutsideMainRootUsesCurrentWorktree(t *testing.T) {
	n, mp, mg := newNamer(t)
	path := "/Users/hansolo/code/worktrees/nu__5969"
	mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/worktrees/nu__5969", nil)
	mp.On("Base", "/Users/hansolo/code/worktrees/nu__5969").Return("nu__5969")

	name, err := gitName(n, path)
	assert.NoError(t, err)
	assert.Equal(t, "nu__5969", name)
}

func TestGitRootName(t *testing.T) {
	t.Run("regular clone, nested subdir collapses to repo root", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/nu/server/subdir"
		mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/nu", nil)
		mp.On("Base", "/Users/hansolo/code/project/nu").Return("nu")

		name, err := gitRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "nu", name)
	})

	t.Run("regular clone, nested subdir inside linked worktree", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/nu/.wk/5969/src"
		mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/nu/.wk/5969", nil)
		mp.On("Base", "/Users/hansolo/code/project/nu/.wk/5969").Return("5969")

		name, err := gitRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "5969", name)
	})

	t.Run("bare repo, nested subdir in worktree", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/project/sesh/main/namer"
		mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/project/sesh/main", nil)
		mp.On("Base", "/Users/hansolo/code/project/sesh/main").Return("main")

		name, err := gitRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "main", name)
	})

	t.Run("bare repo .git suffix, feature worktree", func(t *testing.T) {
		n, mp, mg := newNamer(t)
		path := "/Users/hansolo/code/myrepo/develop"
		mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/myrepo/develop", nil)
		mp.On("Base", "/Users/hansolo/code/myrepo/develop").Return("develop")

		name, err := gitRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "develop", name)
	})

	t.Run("returns empty when ShowTopLevel returns empty", func(t *testing.T) {
		n, _, mg := newNamer(t)
		path := "/Users/hansolo/code/project/sesh/main"
		mg.On("ShowTopLevel", path).Return(false, "", nil)

		name, err := gitRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "", name)
	})

	t.Run("non-git directory", func(t *testing.T) {
		n, _, mg := newNamer(t)
		path := "/Users/hansolo/.config/nvim"
		mg.On("ShowTopLevel", path).Return(false, "", nil)

		name, err := gitRootName(n, path)
		assert.NoError(t, err)
		assert.Equal(t, "", name)
	})
}

func TestGitRootNameLinkedWorktreeOutsideMainRootUsesCurrentWorktree(t *testing.T) {
	n, mp, mg := newNamer(t)
	path := "/Users/hansolo/code/worktrees/nu__5969/src"
	mg.On("ShowTopLevel", path).Return(true, "/Users/hansolo/code/worktrees/nu__5969", nil)
	mp.On("Base", "/Users/hansolo/code/worktrees/nu__5969").Return("nu__5969")

	name, err := gitRootName(n, path)
	assert.NoError(t, err)
	assert.Equal(t, "nu__5969", name)
}
