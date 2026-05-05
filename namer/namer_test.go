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

func TestFromPath(t *testing.T) {
	t.Run("when path does not contain a symlink", func(t *testing.T) {
		mockPathwrap := new(pathwrap.MockPath)
		mockGit := new(git.MockGit)
		mockHome := new(home.MockHome)
		config := model.Config{DirLength: 1}
		n := NewNamer(mockPathwrap, mockGit, mockHome, config)

		t.Run("name for git repo", func(t *testing.T) {
			path := "/Users/josh/config/dotfiles/.config/neovim"
			mockPathwrap.On("EvalSymlinks", path).Return(path, nil)
			mockGit.On("ShowTopLevel", path).Return(true, "/Users/josh/config/dotfiles", nil)
			mockPathwrap.On("Base", "/Users/josh/config/dotfiles").Return("dotfiles")
			name, _ := n.Name(path)
			assert.Equal(t, "dotfiles/_config/neovim", name)
		})

		t.Run("returns base on non-git dir", func(t *testing.T) {
			path := "/Users/josh/.config/neovim"
			mockPathwrap.On("EvalSymlinks", path).Return(path, nil)
			mockGit.On("ShowTopLevel", path).Return(false, "", fmt.Errorf("not a git repository (or any of the parent"))
			mockPathwrap.On("Base", path).Return("neovim")
			name, _ := n.Name(path)
			assert.Equal(t, "neovim", name)
		})
	})

	t.Run("when path contains a symlink", func(t *testing.T) {
		t.Run("name for symlinked file in symlinked git repo", func(t *testing.T) {
			mockPathwrap := new(pathwrap.MockPath)
			mockGit := new(git.MockGit)
			mockHome := new(home.MockHome)
			config := model.Config{DirLength: 1}
			n := NewNamer(mockPathwrap, mockGit, mockHome, config)
			resolved := "/Users/josh/dotfiles/.config/neovim"
			mockPathwrap.On("EvalSymlinks", "/Users/josh/d/.c/neovim").Return(resolved, nil)
			mockGit.On("ShowTopLevel", resolved).Return(true, "/Users/josh/dotfiles", nil)
			mockPathwrap.On("Base", "/Users/josh/dotfiles").Return("dotfiles")
			name, _ := n.Name("/Users/josh/d/.c/neovim")
			assert.Equal(t, "dotfiles/_config/neovim", name)
		})

		t.Run("name for git bare repo", func(t *testing.T) {
			mockPathwrap := new(pathwrap.MockPath)
			mockGit := new(git.MockGit)
			mockHome := new(home.MockHome)
			config := model.Config{DirLength: 1}
			n := NewNamer(mockPathwrap, mockGit, mockHome, config)
			resolved := "/Users/josh/projects/sesh/main"
			mockPathwrap.On("EvalSymlinks", "/Users/josh/p/sesh/main").Return(resolved, nil)
			mockGit.On("ShowTopLevel", resolved).Return(true, "/Users/josh/projects/sesh/main", nil)
			mockPathwrap.On("Base", "/Users/josh/projects/sesh/main").Return("main")
			name, _ := n.Name("/Users/josh/p/sesh/main")
			assert.Equal(t, "main", name)
		})

		t.Run("returns base on non-git dir", func(t *testing.T) {
			mockPathwrap := new(pathwrap.MockPath)
			mockGit := new(git.MockGit)
			mockHome := new(home.MockHome)
			config := model.Config{DirLength: 1}
			n := NewNamer(mockPathwrap, mockGit, mockHome, config)
			resolved := "/Users/josh/.config/neovim"
			mockPathwrap.On("EvalSymlinks", "/Users/josh/c/neovim").Return(resolved, nil)
			mockGit.On("ShowTopLevel", resolved).Return(false, "", fmt.Errorf("not a git repository"))
			mockPathwrap.On("Base", resolved).Return("neovim")
			name, _ := n.Name("/Users/josh/c/neovim")
			assert.Equal(t, "neovim", name)
		})
	})
}
