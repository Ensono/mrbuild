package util

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestFileExist(t *testing.T) {

	// create a new file
	filename := filepath.Join(t.TempDir(), "dummtest.txt")
	err := os.WriteFile(filename, []byte("data"), 0644)

	if err != nil {
		t.Error("Unable to create temporary file")
	} else {

		status := Exists(filename)

		if !status {
			t.Error("File should exist")
		}
	}

}

func TestFileDoesNotExist(t *testing.T) {

	status := Exists("path/to/non/existent/file")
	if status {
		t.Error("File should not exist")
	}
}

func TestFileStat(t *testing.T) {

	// build up the test tables
	tables := []struct {
		path     string
		dir      string
		name     string
		ext      string
		filename string
	}{
		{
			"./mrbuild.yaml",
			".",
			"mrbuild",
			"yaml",
			"mrbuild.yaml",
		},
		{
			"c:\\users\\testing\\tools\\mrbuild.yaml",
			"c:\\users\\testing\\tools",
			"mrbuild",
			"yaml",
			"mrbuild.yaml",
		},
		{
			"/home/tester/bin/config.toml",
			"/home/tester/bin",
			"config",
			"toml",
			"config.toml",
		},
	}

	for _, table := range tables {

		// get the file stat object and compare against the expected values
		fs := NewFileStat(table.path)

		// As the paths are going to be different on different platforms, the
		// tests will split up the values so that that they can be compared with
		// each other as slices
		re := regexp.MustCompile(`(?m)\\|/`)
		tableDirectoryParts := re.Split(filepath.Dir(table.path), -1)
		fsDirectoryParts := re.Split(fs.Directory, -1)

		if !reflect.DeepEqual(tableDirectoryParts, fsDirectoryParts) {
			t.Errorf("Directory expected '%s', actual '%s'",
				strings.Join(tableDirectoryParts, string(os.PathSeparator)),
				strings.Join(fsDirectoryParts, string(os.PathSeparator)),
			)
		}

		if fs.Name != table.name {
			t.Errorf("Name expected '%s', actual '%s'", table.name, fs.Name)
		}

		if fs.Extension != table.ext {
			t.Errorf("Extension expected '%s', actual '%s'", table.ext, fs.Extension)
		}

		if fs.Filename != table.filename {
			t.Errorf("Filename expected '%s', actual '%s'", table.filename, fs.Filename)
		}
	}
}
