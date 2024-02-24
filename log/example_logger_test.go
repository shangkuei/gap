package log_test

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"time"

	"github.com/shangkuei/gap/log"
	"github.com/spf13/afero"
)

//go:embed testdata/console.toml
var consoleTOML embed.FS

//go:embed testdata/file.toml
var fileTOML embed.FS

type embedFS embed.FS

func (e *embedFS) Create(name string) (afero.File, error) {
	return nil, errors.New("not supported")
}

func (e *embedFS) Mkdir(name string, perm os.FileMode) error {
	return errors.New("not supported")
}

func (e *embedFS) MkdirAll(path string, perm os.FileMode) error {
	return errors.New("not supported")
}

func (e *embedFS) Open(path string) (afero.File, error) {
	f, err := (*embed.FS)(e).Open(path)
	if err != nil {
		return nil, err
	}
	return file{File: f}, nil
}

func (e *embedFS) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := (*embed.FS)(e).Open(name)
	if err != nil {
		return nil, err
	}
	return file{File: f}, nil
}

func (e *embedFS) Remove(name string) error {
	return errors.New("not supported")
}

func (e *embedFS) RemoveAll(path string) error {
	return errors.New("not supported")
}

func (e *embedFS) Rename(oldname, newname string) error {
	return errors.New("not supported")
}

func (e *embedFS) Stat(name string) (os.FileInfo, error) {
	f, err := (*embed.FS)(e).Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (e *embedFS) Name() string {
	return "embedFS"
}

func (e *embedFS) Chmod(name string, mode os.FileMode) error {
	return errors.New("not supported")
}

func (e *embedFS) Chown(name string, uid, gid int) error {
	return errors.New("not supported")
}

func (e *embedFS) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return errors.New("not supported")
}

type file struct {
	fs.File
}

func (f file) Close() error {
	return f.File.Close()
}

func (f file) Read(p []byte) (n int, err error) {
	return f.File.Read(p)
}

func (f file) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errors.New("not supported")
}

func (f file) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("not supported")
}

func (f file) Write(p []byte) (n int, err error) {
	return 0, errors.New("not supported")
}

func (f file) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, errors.New("not supported")
}

func (f file) Name() string {
	stat, _ := f.File.Stat()
	return stat.Name()
}
func (f file) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errors.New("not supported")
}

func (f file) Readdirnames(n int) ([]string, error) {
	return nil, errors.New("not supported")
}

func (f file) Stat() (os.FileInfo, error) {
	return f.File.Stat()
}

func (f file) Sync() error {
	return nil
}

func (f file) Truncate(size int64) error {
	return errors.New("not supported")
}

func (f file) WriteString(s string) (ret int, err error) {
	return 0, errors.New("not supported")
}

type File interface {
	Readdirnames(n int) ([]string, error)
	Stat() (os.FileInfo, error)
	Sync() error
	Truncate(size int64) error
	WriteString(s string) (ret int, err error)
}

func ExampleLogger_console() {
	config, err := log.ConfigurationFromViper(log.ViperConfiguration{
		ConfigFile: "testdata/console.toml",
		FileSystem: (*embedFS)(&consoleTOML),
	})
	if err != nil {
		panic(err)
	}
	config.ReplaceAttr = func(groups []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.TimeKey && len(groups) == 0 {
			return slog.Attr{}
		}
		return attr
	}
	logger := log.Logger(config)

	logger.Info("Hello, World!")

	// Output:
	// INF Hello, World!
}

func ExampleLogger_file() {
	config, err := log.ConfigurationFromViper(log.ViperConfiguration{
		ConfigFile: "testdata/file.toml",
		FileSystem: (*embedFS)(&fileTOML),
	})
	if err != nil {
		panic(err)
	}
	config.ReplaceAttr = func(groups []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.TimeKey && len(groups) == 0 {
			return slog.Attr{}
		}
		return attr
	}
	logger := log.Logger(config)
	defer os.Remove("file.log")

	logger.Info("Hello, World!")

	data, err := os.ReadFile("file.log")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	// Output:
	// level=INFO msg="Hello, World!"
}
