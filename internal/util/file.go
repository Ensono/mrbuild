package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FileStat struct {
	Name      string
	Filename  string
	Extension string
	Directory string
}

// Exists returns a boolean stating if a file exists or not
func Exists(filePath string) bool {
	_, err := os.Stat(filePath)

	if errors.Is(err, os.ErrNotExist) || err != nil {
		return false
	}

	return true
}

// IsInputFromPipe determines if the command has been run and the
// data is being passed via a pipe
func IsInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

// NewFileStat takes the path to a file and returns an object containing
// information about the file, such as the filename, its extension.
// It also contains the name of the file without the extension, which can be
// useful in identifying project files
//
// This method does not use the filepath pacakge as it is not cross platform
// For example on Linux, it does not understand a Windows path with backslashes
func NewFileStat(path string) *FileStat {
	fs := new(FileStat)

	// set the pattern to split up the path into consitutent parts
	re := regexp.MustCompile(`(?m)\\|/`)
	splitPath := re.Split(path, -1)

	// set the parent path and the filename of the fs object
	fs.Filename = splitPath[len(splitPath)-1]                                           // filepath.Base(path)
	fs.Directory = strings.Join(splitPath[:len(splitPath)-1], string(os.PathSeparator)) // filepath.Dir(path)

	// split the filename to get the extension
	splitFilename := strings.Split(fs.Filename, ".")
	fs.Extension = splitFilename[len(splitFilename)-1]

	// get the name of the file without the extension
	fs.Name = fs.Filename[0 : len(fs.Filename)-len(filepath.Ext(path))]

	return fs
}

func (fs *FileStat) GetFullPath() string {
	return fmt.Sprintf("%s%s%s", fs.Directory, string(os.PathSeparator), fs.Filename)
}
