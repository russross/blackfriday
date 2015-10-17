package mmark

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// FileSystem implements access to files. The elements in a file path are
// separated by slash ('/', U+002F) characters, regardless of host
// operating system convention.
type FileSystem interface {
	// ReadFile reads the file named by filename and returns the contents.
	ReadFile(name string) ([]byte, error)
}

// dir implements FileSystem using the native file system restricted to a
// specific directory tree.
//
// While the FileSystem.ReadFile method takes '/'-separated paths,
// a Dir's string value is a filename on the native file system,
// not a URL, so it is separated by filepath.Separator,
// which isn't necessarily '/'.
//
// An empty Dir is treated as "."
type Dir string

// ReadFile reads the file named by filename and returns the contents.
func (d Dir) ReadFile(name string) ([]byte, error) {
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	fullname := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	return ioutil.ReadFile(fullname)
}

// virtualFS implements FileSystem using an in-memory map to define the FileSystem
// this is used for testing
type virtualFS map[string]string

// ReadFile returns the content for appropriate filename
// if the name does not exist, then it will return os.ErrNotExist
func (fs virtualFS) ReadFile(name string) ([]byte, error) {
	search := path.Clean("/" + name)
	content, ok := fs[search]
	if !ok {
		return nil, os.ErrNotExist
	}
	return []byte(content), nil
}

func absname(cwd string, name string) string {
	if len(name) > 0 && name[0] == '/' {
		return name
	}
	return path.Join(cwd, name)
}
