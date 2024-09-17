package io

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-rat/utils/env"
	"github.com/stretchr/testify/suite"
)

type IOTestSuite struct {
	suite.Suite
}

func TestIOTestSuite(t *testing.T) {
	suite.Run(t, &IOTestSuite{})
}

func (s *IOTestSuite) SetupTest() {
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		s.NoError(os.MkdirAll("testdata", 0755))
	}
}

func (s *IOTestSuite) TearDownTest() {
	s.NoError(os.RemoveAll("testdata"))
}

func (s *IOTestSuite) TestWriteCreatesFileWithCorrectContent() {
	path := "testdata/write_test.txt"
	data := "Hello, World!"
	permission := os.FileMode(0644)

	s.NoError(Write(path, data, permission))

	content, err := Read(path)
	s.NoError(err)
	s.Equal(data, content)
}

func (s *IOTestSuite) TestWriteAppendAppendsToFile() {
	path := "testdata/append_test.txt"
	initialData := "Hello"
	appendData := ", World!"

	s.NoError(Write(path, initialData, 0644))
	s.NoError(WriteAppend(path, appendData))

	content, err := Read(path)
	s.NoError(err)
	s.Equal("Hello, World!", content)
}

func (s *IOTestSuite) TestCompress() {
	src := []string{"testdata/compress_test1.txt", "testdata/compress_test2.txt"}
	err := Write(src[0], "File 1", 0644)
	s.NoError(err)
	err = Write(src[1], "File 2", 0644)
	s.NoError(err)

	err = Compress(src, "testdata/compress_test.zip", Zip)
	s.NoError(err)
}

func (s *IOTestSuite) TestUnCompress() {
	src := []string{"testdata/uncompress_test1.txt", "testdata/uncompress_test2.txt"}
	err := Write(src[0], "File 1", 0644)
	s.NoError(err)
	err = Write(src[1], "File 2", 0644)
	s.NoError(err)

	err = Compress(src, "testdata/uncompress_test.zip", Zip)
	s.NoError(err)

	err = UnCompress("testdata/uncompress_test.zip", "testdata/uncompressed", Zip)
	s.NoError(err)

	data, err := Read("testdata/uncompressed/uncompress_test1.txt")
	s.NoError(err)
	s.Equal("File 1", data)

	data, err = Read("testdata/uncompressed/uncompress_test2.txt")
	s.NoError(err)
	s.Equal("File 2", data)
}

func (s *IOTestSuite) TestRemoveDeletesFileOrDirectory() {
	path := "testdata/remove_test"
	s.NoError(Mkdir(path, 0755))
	s.DirExists(path)

	s.NoError(Remove(path))
	s.NoDirExists(path)
}

func (s *IOTestSuite) TestMkdirCreatesDirectory() {
	path := "testdata/mkdir_test"
	s.NoError(Mkdir(path, 0755))
	s.DirExists(path)
}

func (s *IOTestSuite) TestChmodChangesPermissions() {
	if env.IsWindows() {
		s.T().Skip("Skipping on Windows")
	}
	path := "testdata/chmod_test.txt"
	s.NoError(Write(path, "test", 0644))

	s.NoError(Chmod(path, 0755))
	info, err := os.Stat(path)
	s.NoError(err)
	s.Equal(os.FileMode(0755), info.Mode().Perm())
}

func (s *IOTestSuite) TestChownChangesOwner() {
	if env.IsWindows() {
		s.T().Skip("Skipping on Windows")
	}
	path := "testdata/chown_test.txt"
	s.NoError(Write(path, "test", 0644))

	s.NoError(Chown(path, "root", "root"))
}

func (s *IOTestSuite) TestExistsReturnsTrueForExistingPath() {
	path := "testdata/exists_test.txt"
	s.NoError(Write(path, "test", 0644))
	s.True(Exists(path))
}

func (s *IOTestSuite) TestExistsReturnsFalseForNonExistingPath() {
	path := "testdata/nonexistent.txt"
	s.False(Exists(path))
}

func (s *IOTestSuite) TestEmptyReturnsTrueForEmptyDirectory() {
	path := "testdata/empty_test"
	s.NoError(Mkdir(path, 0755))
	s.True(Empty(path))
}

func (s *IOTestSuite) TestEmptyReturnsFalseForNonEmptyDirectory() {
	path := "testdata/nonempty_test"
	s.NoError(Mkdir(path, 0755))
	s.NoError(Write(filepath.Join(path, "file.txt"), "test", 0644))
	s.False(Empty(path))
}

func (s *IOTestSuite) TestMvMovesFile() {
	src := "testdata/mv_src.txt"
	dst := "testdata/mv_dst.txt"
	s.NoError(Write(src, "test", 0644))

	s.NoError(Mv(src, dst))
	s.FileExists(dst)
	s.NoFileExists(src)
}

func (s *IOTestSuite) TestCpCopiesFile() {
	src := "testdata/cp_src.txt"
	dst := "testdata/cp_dst.txt"
	s.NoError(Write(src, "test", 0644))

	s.NoError(Cp(src, dst))
	s.FileExists(dst)
	s.FileExists(src)
}

func (s *IOTestSuite) TestSizeReturnsCorrectSize() {
	path := "testdata/size_test.txt"
	data := "12345"
	s.NoError(Write(path, data, 0644))

	size, err := Size(path)
	s.NoError(err)
	s.Equal(int64(len(data)), size)
}

func (s *IOTestSuite) TestTempDirCreatesTemporaryDirectory() {
	dir, err := TempDir("tempdir_test")
	s.NoError(err)
	s.DirExists(dir)
	s.NoError(Remove(dir))
}

func (s *IOTestSuite) TestReadDirReturnsDirectoryEntries() {
	path := "testdata/readdir_test"
	s.NoError(Mkdir(path, 0755))
	s.NoError(Write(filepath.Join(path, "file1.txt"), "test", 0644))
	s.NoError(Write(filepath.Join(path, "file2.txt"), "test", 0644))

	entries, err := ReadDir(path)
	s.NoError(err)
	s.Len(entries, 2)
}

func (s *IOTestSuite) TestIsDirReturnsTrueForDirectory() {
	path := "testdata/isdir_test"
	s.NoError(Mkdir(path, 0755))
	s.True(IsDir(path))
}

func (s *IOTestSuite) TestIsDirReturnsFalseForFile() {
	path := "testdata/isfile_test.txt"
	s.NoError(Write(path, "test", 0644))
	s.False(IsDir(path))
}
