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
	// ReadFile reads the file named by filename and returns the contents.
	ReadFile(name []byte) ([]byte, error)
}

// dir implements fileSystem using the native file system restricted to a
// specific directory tree.
//
// While the fileSystem.ReadFile method takes '/'-separated paths,
// a dir's string value is a filename on the native file system,
// not a URL, so it is separated by filepath.Separator,
// which isn't necessarily '/'.
//
// An empty dir is treated as "."
type dir string

// ReadFile reads the file named by filename and returns the contents.
func (d dir) ReadFile(name []byte) ([]byte, error) {
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	fullname := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+string(name))))
	return ioutil.ReadFile(fullname)
}

// virtualFS implements fileSystem using an in-memory map to define the fileSystem
// this is used for testing
type virtualFS map[string]string

// ReadFile returns the content for appropriate filename
// if the name does not exist, then it will return os.ErrNotExist
func (fs virtualFS) ReadFile(name []byte) ([]byte, error) {
	search := path.Clean("/" + string(name))
	content, ok := fs[search]
	if !ok {
		return nil, os.ErrNotExist
	}
	return []byte(content), nil
}

func absname(cwd string, name []byte) []byte {
	if len(name) > 0 && name[0] == '/' {
		return name
	}
	return []byte(path.Join(cwd, string(name)))
}
