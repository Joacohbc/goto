package tests

import (
	"goto/src/core"
	"goto/src/gpath"
	"goto/src/utils"
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePath(t *testing.T) {
	t.Run("onlyDirectory with valid directory", func(t *testing.T) {
		dir := t.TempDir()

		result, err := core.ResolvePath([]string{dir}, true, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		absDir, _ := filepath.Abs(dir)
		if result != absDir {
			t.Errorf("got %q, want %q", result, absDir)
		}
	})

	t.Run("onlyDirectory with non-existent path", func(t *testing.T) {
		_, err := core.ResolvePath([]string{"/non/existent/path"}, true, false)
		if err == nil {
			t.Error("expected error for non-existent directory")
		}
	})

	t.Run("onlyDirectory with file path", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "file.txt")
		os.WriteFile(file, []byte("test"), 0644)

		_, err := core.ResolvePath([]string{file}, true, false)
		if err == nil {
			t.Error("expected error for file path (not a directory)")
		}
	})

	t.Run("onlyDirectory with empty args", func(t *testing.T) {
		_, err := core.ResolvePath([]string{""}, true, false)
		if err == nil {
			t.Error("expected error for empty path")
		}
	})

	t.Run("onlyDirectory joins multiple args", func(t *testing.T) {
		// Create a nested directory structure
		base := t.TempDir()
		nested := filepath.Join(base, "sub", "dir")
		os.MkdirAll(nested, 0755)

		result, err := core.ResolvePath([]string{base, "sub", "dir"}, true, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		absNested, _ := filepath.Abs(nested)
		if result != absNested {
			t.Errorf("got %q, want %q", result, absNested)
		}
	})

	t.Run("resolve by abbreviation", func(t *testing.T) {
		_, cleanup := resetConfigFile(t, false)
		defer cleanup()

		// Load current gpaths and add a known entry
		gpathsList, err := utils.LoadGPaths(false)
		if err != nil {
			t.Fatalf("failed to load gpaths: %v", err)
		}

		dir := t.TempDir()
		gpathsList = append(gpathsList, gpath.GotoPath{
			Path:         dir,
			Abbreviation: "testabbr",
		})
		if err := utils.UpdateGPaths(false, gpathsList); err != nil {
			t.Fatalf("failed to update gpaths: %v", err)
		}

		result, err := core.ResolvePath([]string{"testabbr"}, false, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result != dir {
			t.Errorf("got %q, want %q", result, dir)
		}
	})

	t.Run("resolve by index", func(t *testing.T) {
		_, cleanup := resetConfigFile(t, false)
		defer cleanup()

		gpathsList, err := utils.LoadGPaths(false)
		if err != nil {
			t.Fatalf("failed to load gpaths: %v", err)
		}

		// Index 0 should resolve to the first entry's path
		result, err := core.ResolvePath([]string{"0"}, false, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result != gpathsList[0].Path {
			t.Errorf("got %q, want %q", result, gpathsList[0].Path)
		}
	})

	t.Run("resolve plain directory path (not index or abbreviation)", func(t *testing.T) {
		_, cleanup := resetConfigFile(t, false)
		defer cleanup()

		dir := t.TempDir()

		result, err := core.ResolvePath([]string{dir}, false, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		absDir, _ := filepath.Abs(dir)
		if result != absDir {
			t.Errorf("got %q, want %q", result, absDir)
		}
	})

	t.Run("resolve invalid path (not index, abbreviation, or directory)", func(t *testing.T) {
		_, cleanup := resetConfigFile(t, false)
		defer cleanup()

		_, err := core.ResolvePath([]string{"/does/not/exist"}, false, false)
		if err == nil {
			t.Error("expected error for non-existent path that is not an abbreviation or index")
		}
	})

	t.Run("resolve with temporal flag", func(t *testing.T) {
		_, cleanup := resetConfigFile(t, true)
		defer cleanup()

		// Load temporal gpaths and add a known entry
		gpathsList, err := utils.LoadGPaths(true)
		if err != nil {
			t.Fatalf("failed to load gpaths: %v", err)
		}

		dir := t.TempDir()
		gpathsList = append(gpathsList, gpath.GotoPath{
			Path:         dir,
			Abbreviation: "tmpabbr",
		})
		if err := utils.UpdateGPaths(true, gpathsList); err != nil {
			t.Fatalf("failed to update gpaths: %v", err)
		}

		result, err := core.ResolvePath([]string{"tmpabbr"}, false, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result != dir {
			t.Errorf("got %q, want %q", result, dir)
		}
	})
}
