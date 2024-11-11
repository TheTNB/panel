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
	s.NoError(WriteAppend(path, appendData, 0644))

	content, err := Read(path)
	s.NoError(err)
	s.Equal("Hello, World!", content)
}

func (s *IOTestSuite) TestCompress() {
	abs, err := filepath.Abs("testdata")
	s.NoError(err)
	src := []string{"compress_test1.txt", "compress_test2.txt"}
	err = Write(filepath.Join(abs, src[0]), "File 1", 0644)
	s.NoError(err)
	err = Write(filepath.Join(abs, src[1]), "File 2", 0644)
	s.NoError(err)

	err = Compress(abs, src, filepath.Join(abs, "compress_test.zip"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "compress_test.bz2"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "compress_test.tar"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "compress_test.gz"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "compress_test.tar.gz"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "compress_test.tgz"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "compress_test.xz"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "compress_test.7z"))
	s.NoError(err)

	s.NoError(Remove("testdata"))
}

func (s *IOTestSuite) TestUnCompress() {
	abs, err := filepath.Abs("testdata")
	s.NoError(err)
	src := []string{"uncompress_test1.txt", "uncompress_test2.txt"}
	err = Write(filepath.Join(abs, src[0]), "File 1", 0644)
	s.NoError(err)
	err = Write(filepath.Join(abs, src[1]), "File 2", 0644)
	s.NoError(err)

	err = Compress(abs, src, filepath.Join(abs, "uncompress_test.zip"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "uncompress_test.bz2"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "uncompress_test.tar"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "uncompress_test.gz"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "uncompress_test.tar.gz"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "uncompress_test.tgz"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "uncompress_test.xz"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "uncompress_test.7z"))
	s.NoError(err)

	err = UnCompress(filepath.Join(abs, "uncompress_test.zip"), filepath.Join(abs, "uncompressed_zip"))
	s.NoError(err)
	data, err := Read("testdata/uncompressed_zip/uncompress_test1.txt")
	s.NoError(err)
	s.Equal("File 1", data)
	data, err = Read("testdata/uncompressed_zip/uncompress_test2.txt")
	s.NoError(err)
	s.Equal("File 2", data)
	err = UnCompress(filepath.Join(abs, "uncompress_test.bz2"), filepath.Join(abs, "uncompressed_bz2"))
	s.NoError(err)
	data, err = Read("testdata/uncompressed_bz2/uncompress_test1.txt")
	s.NoError(err)
	s.Equal("File 1", data)
	data, err = Read("testdata/uncompressed_bz2/uncompress_test2.txt")
	s.NoError(err)
	s.Equal("File 2", data)
	err = UnCompress(filepath.Join(abs, "uncompress_test.tar"), filepath.Join(abs, "uncompressed_tar"))
	s.NoError(err)
	data, err = Read("testdata/uncompressed_tar/uncompress_test1.txt")
	s.NoError(err)
	s.Equal("File 1", data)
	data, err = Read("testdata/uncompressed_tar/uncompress_test2.txt")
	s.NoError(err)
	s.Equal("File 2", data)
	err = UnCompress(filepath.Join(abs, "uncompress_test.gz"), filepath.Join(abs, "uncompressed_gz"))
	s.NoError(err)
	data, err = Read("testdata/uncompressed_gz/uncompress_test1.txt")
	s.NoError(err)
	s.Equal("File 1", data)
	data, err = Read("testdata/uncompressed_gz/uncompress_test2.txt")
	s.NoError(err)
	s.Equal("File 2", data)
	err = UnCompress(filepath.Join(abs, "uncompress_test.tar.gz"), filepath.Join(abs, "uncompressed_tar_gz"))
	s.NoError(err)
	data, err = Read("testdata/uncompressed_tar_gz/uncompress_test1.txt")
	s.NoError(err)
	s.Equal("File 1", data)
	data, err = Read("testdata/uncompressed_tar_gz/uncompress_test2.txt")
	s.NoError(err)
	s.Equal("File 2", data)
	err = UnCompress(filepath.Join(abs, "uncompress_test.tgz"), filepath.Join(abs, "uncompressed_tgz"))
	s.NoError(err)
	data, err = Read("testdata/uncompressed_tgz/uncompress_test1.txt")
	s.NoError(err)
	s.Equal("File 1", data)
	data, err = Read("testdata/uncompressed_tgz/uncompress_test2.txt")
	s.NoError(err)
	s.Equal("File 2", data)
	err = UnCompress(filepath.Join(abs, "uncompress_test.xz"), filepath.Join(abs, "uncompressed_xz"))
	s.NoError(err)
	data, err = Read("testdata/uncompressed_xz/uncompress_test1.txt")
	s.NoError(err)
	s.Equal("File 1", data)
	data, err = Read("testdata/uncompressed_xz/uncompress_test2.txt")
	s.NoError(err)
	s.Equal("File 2", data)
	err = UnCompress(filepath.Join(abs, "uncompress_test.7z"), filepath.Join(abs, "uncompressed_7z"))
	s.NoError(err)
	data, err = Read("testdata/uncompressed_7z/uncompress_test1.txt")
	s.NoError(err)
	s.Equal("File 1", data)
	data, err = Read("testdata/uncompressed_7z/uncompress_test2.txt")
	s.NoError(err)
	s.Equal("File 2", data)

	s.NoError(Remove("testdata"))
}

func (s *IOTestSuite) TestListCompress() {
	abs, err := filepath.Abs("testdata")
	s.NoError(err)
	src := []string{"list_archive_test1.txt", "list_archive_test2.txt"}
	err = Write(filepath.Join(abs, src[0]), "File 1", 0644)
	s.NoError(err)
	err = Write(filepath.Join(abs, src[1]), "File 2", 0644)
	s.NoError(err)

	err = Compress(abs, src, filepath.Join(abs, "list_archive_test.zip"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "list_archive_test.bz2"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "list_archive_test.tar"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "list_archive_test.gz"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "list_archive_test.tar.gz"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "list_archive_test.tgz"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "list_archive_test.xz"))
	s.NoError(err)
	err = Compress(abs, src, filepath.Join(abs, "list_archive_test.7z"))
	s.NoError(err)

	list, err := ListCompress(filepath.Join(abs, "list_archive_test.zip"))
	s.NoError(err)
	s.Len(list, 2)
	list, err = ListCompress(filepath.Join(abs, "list_archive_test.bz2"))
	s.NoError(err)
	s.Len(list, 2)
	list, err = ListCompress(filepath.Join(abs, "list_archive_test.tar"))
	s.NoError(err)
	s.Len(list, 2)
	list, err = ListCompress(filepath.Join(abs, "list_archive_test.gz"))
	s.NoError(err)
	s.Len(list, 2)
	list, err = ListCompress(filepath.Join(abs, "list_archive_test.tar.gz"))
	s.NoError(err)
	s.Len(list, 2)
	list, err = ListCompress(filepath.Join(abs, "list_archive_test.tgz"))
	s.NoError(err)
	s.Len(list, 2)
	list, err = ListCompress(filepath.Join(abs, "list_archive_test.xz"))
	s.NoError(err)
	s.Len(list, 2)
	list, err = ListCompress(filepath.Join(abs, "list_archive_test.7z"))
	s.NoError(err)
	s.Len(list, 2)

	s.NoError(Remove("testdata"))
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
