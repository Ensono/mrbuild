package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
func NewFileStat(path string) *FileStat {
	fs := new(FileStat)

	// set the parent path and the filename of the fs object
	fs.Filename = filepath.Base(path)
	fs.Directory = filepath.Dir(path)
	fs.Extension = filepath.Ext(path)[1:]

	// get the name of the file without the extension
	fs.Name = fs.Filename[0 : len(fs.Filename)-len(filepath.Ext(path))]

	return fs
}

func (fs *FileStat) GetFullPath() string {
	return fmt.Sprintf("%s%s%s", fs.Directory, string(os.PathSeparator), fs.Filename)
}
