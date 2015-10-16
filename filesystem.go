package mmark

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// fileSystem implements access to files. The elements in a file path are
// separated by slash ('/', U+002F) characters, regardless of host
// operating system convention.
type fileSystem interface {
	// readFile reads the file named by filename and returns the contents.
	readFile(name string) ([]byte, error)
}

// dir implements fileSystem using the native file system restricted to a
// specific directory tree.
//
// While the fileSystem.readFile method takes '/'-separated paths,
// a Dir's string value is a filename on the native file system,
// not a URL, so it is separated by filepath.Separator,
// which isn't necessarily '/'.
//
// An empty Dir is treated as "."
type dir string

// readFile reads the file named by filename and returns the contents.
func (d dir) readFile(name string) ([]byte, error) {
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	fullname := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	return ioutil.ReadFile(fullname)
}

// virtualFS implements fileSystem using an in-memory map to define the filesystem
type virtualFS map[string]string

// readFile returns the content for appropriate filename
// if the name does not exist, then it will return os.ErrNotExist
func (fs virtualFS) readFile(name string) ([]byte, error) {
	content, ok := fs[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return []byte(content), nil
}
