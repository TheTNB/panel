package tools

import (
	"os/user"
	"path/filepath"
	"testing"
	"time"

	"github.com/goravel/framework/support/env"
	"github.com/stretchr/testify/suite"
)

type SystemHelperTestSuite struct {
	suite.Suite
}

func TestSystemHelperTestSuite(t *testing.T) {
	suite.Run(t, &SystemHelperTestSuite{})
}

func (s *SystemHelperTestSuite) WriteCreatesFileWithCorrectContent() {
	filePath, _ := TempFile("testfile")

	err := Write(filePath.Name(), "test data", 0644)
	s.Nil(err)

	content, _ := Read(filePath.Name())
	s.Equal("test data", content)
	s.Nil(filePath.Close())
	s.Nil(Remove(filePath.Name()))
}

func (s *SystemHelperTestSuite) WriteCreatesDirectoriesIfNeeded() {
	filePath, _ := TempFile("testdir/testfile")

	err := Write(filePath.Name(), "test data", 0644)
	s.Nil(err)

	content, _ := Read(filePath.Name())
	s.Equal("test data", content)
	s.Nil(filePath.Close())
	s.Nil(Remove(filePath.Name()))
}

func (s *SystemHelperTestSuite) WriteFailsIfDirectoryCannotBeCreated() {
	filePath := "/nonexistent/testfile"

	err := Write(filePath, "test data", 0644)
	s.NotNil(err)
}

func (s *SystemHelperTestSuite) WriteFailsIfFileCannotBeWritten() {
	filePath, _ := TempFile("testfile")
	s.Nil(filePath.Close())
	s.Nil(Chmod(filePath.Name(), 0400))

	err := Write(filePath.Name(), "test data", 0644)
	s.NotNil(err)

	s.Nil(Chmod(filePath.Name(), 0644))
	s.Nil(Remove(filePath.Name()))
}

func (s *SystemHelperTestSuite) WriteAppendSuccessfullyAppendsDataToFile() {
	filePath, _ := TempFile("testfile")

	err := Write(filePath.Name(), "initial data", 0644)
	s.Nil(err)

	err = WriteAppend(filePath.Name(), " appended data")
	s.Nil(err)

	content, _ := Read(filePath.Name())
	s.Equal("initial data appended data", content)
	s.Nil(filePath.Close())
	s.Nil(Remove(filePath.Name()))
}

func (s *SystemHelperTestSuite) WriteAppendCreatesFileIfNotExists() {
	filePath, _ := TempFile("testfile")
	s.Nil(filePath.Close())
	s.Nil(Remove(filePath.Name()))

	err := WriteAppend(filePath.Name(), "test data")
	s.Nil(err)

	content, _ := Read(filePath.Name())
	s.Equal("test data", content)
	s.Nil(Remove(filePath.Name()))
}

func (s *SystemHelperTestSuite) WriteAppendReturnsErrorIfPathIsADirectory() {
	dirPath, _ := TempDir("testdir")

	err := WriteAppend(dirPath, "test data")
	s.NotNil(err)

	s.Nil(Remove(dirPath))
}

func (s *SystemHelperTestSuite) ReadSuccessfullyReadsFileContent() {
	filePath, _ := TempFile("testfile")

	err := Write(filePath.Name(), "test data", 0644)
	s.Nil(err)

	content, err := Read(filePath.Name())
	s.Nil(err)
	s.Equal("test data", content)

	s.Nil(filePath.Close())
	s.Nil(Remove(filePath.Name()))
}

func (s *SystemHelperTestSuite) ReadReturnsErrorForNonExistentFile() {
	_, err := Read("/nonexistent/testfile")
	s.NotNil(err)
}

func (s *SystemHelperTestSuite) RemoveSuccessfullyRemovesFile() {
	filePath, _ := TempFile("testfile")

	err := Write(filePath.Name(), "test data", 0644)
	s.Nil(err)

	err = Remove(filePath.Name())
	s.Nil(err)

	s.False(Exists(filePath.Name()))
}

func (s *SystemHelperTestSuite) RemoveReturnsErrorForNonExistentFile() {
	err := Remove("/nonexistent/testfile")
	s.NotNil(err)
}

func (s *SystemHelperTestSuite) TestExec() {
	output, err := Exec("echo test")
	s.Equal("test", output)
	s.Nil(err)
}

func (s *SystemHelperTestSuite) TestExecAsync() {
	command := "echo test > test.txt"
	if env.IsWindows() {
		command = "echo test> test.txt"
	}

	err := ExecAsync(command)
	s.Nil(err)

	time.Sleep(time.Second)

	content, err := Read("test.txt")
	s.Nil(err)

	condition := "test\n"
	if env.IsWindows() {
		condition = "test\r\n"
	}
	s.Equal(condition, content)
	s.Nil(Remove("test.txt"))
}

func (s *SystemHelperTestSuite) TestMkdir() {
	dirPath, _ := TempDir("testdir")

	s.Nil(Mkdir(dirPath, 0755))
	s.Nil(Remove(dirPath))
}

func (s *SystemHelperTestSuite) TestChmod() {
	filePath, _ := TempFile("testfile")

	err := Write(filePath.Name(), "test data", 0644)
	s.Nil(err)

	s.Nil(Chmod(filePath.Name(), 0755))
	s.Nil(filePath.Close())
	s.Nil(Remove(filePath.Name()))
}

func (s *SystemHelperTestSuite) TestChown() {
	filePath, _ := TempFile("testfile")

	err := Write(filePath.Name(), "test data", 0644)
	s.Nil(err)

	currentUser, err := user.Current()
	s.Nil(err)
	groups, err := currentUser.GroupIds()
	s.Nil(err)

	err = Chown(filePath.Name(), currentUser.Username, groups[0])
	if env.IsWindows() {
		s.NotNil(err)
	} else {
		s.Nil(err)
	}
	s.Nil(filePath.Close())
	s.Nil(Remove(filePath.Name()))
}

func (s *SystemHelperTestSuite) TestExists() {
	filePath, _ := TempFile("testfile")

	s.True(Exists(filePath.Name()))
	s.False(Exists("123"))
	s.Nil(filePath.Close())
	s.Nil(Remove(filePath.Name()))
}

func (s *SystemHelperTestSuite) TestEmpty() {
	filePath, _ := TempFile("testfile")

	s.True(Empty(filePath.Name()))
	if env.IsWindows() {
		s.True(Empty("C:\\Windows\\System32\\drivers\\etc\\hosts"))
	} else {
		s.True(Empty("/etc/hosts"))
	}
	s.Nil(filePath.Close())
	s.Nil(Remove(filePath.Name()))
}

func (s *SystemHelperTestSuite) TestMv() {
	filePath, _ := TempFile("testfile")

	err := Write(filePath.Name(), "test data", 0644)
	s.Nil(err)

	newFilePath, _ := TempFile("testfile2")

	s.Nil(newFilePath.Close())
	s.Nil(filePath.Close())

	s.Nil(Mv(filePath.Name(), newFilePath.Name()))
	s.False(Exists(filePath.Name()))
	s.Nil(Remove(newFilePath.Name()))
	s.Nil(Remove(filePath.Name()))
}

func (s *SystemHelperTestSuite) TestCp() {
	tempDir, _ := TempDir("testdir")

	err := Write(filepath.Join(tempDir, "testfile"), "test data", 0644)
	s.Nil(err)

	s.Nil(Cp(filepath.Join(tempDir, "testfile"), filepath.Join(tempDir, "testfile2")))
	s.True(Exists(filepath.Join(tempDir, "testfile2")))
	s.Nil(Remove(tempDir))
}

func (s *SystemHelperTestSuite) TestSize() {
	filePath, _ := TempFile("testfile")

	err := Write(filePath.Name(), "test data", 0644)
	s.Nil(err)

	size, err := Size(filePath.Name())
	s.Nil(err)
	s.Equal(int64(len("test data")), size)
	s.Nil(filePath.Close())
	s.Nil(Remove(filePath.Name()))
}

func (s *SystemHelperTestSuite) TestFileInfo() {
	filePath, _ := TempFile("testfile")

	err := Write(filePath.Name(), "test data", 0644)
	s.Nil(err)

	info, err := FileInfo(filePath.Name())
	s.Nil(err)
	s.Equal(filepath.Base(filePath.Name()), info.Name())
	s.Nil(filePath.Close())
	s.Nil(Remove(filePath.Name()))
}

func (s *SystemHelperTestSuite) TestUnArchiveSuccessfullyUnarchivesFile() {
	file, _ := TempFile("test")
	dstDir, _ := TempDir("archive")

	err := Write(file.Name(), "test data", 0644)
	s.Nil(err)

	err = Archive([]string{file.Name()}, filepath.Join(dstDir, "test.zip"))
	s.Nil(err)
	s.FileExists(filepath.Join(dstDir, "test.zip"))

	err = UnArchive(filepath.Join(dstDir, "test.zip"), dstDir)
	s.Nil(err)
	s.FileExists(filepath.Join(dstDir, filepath.Base(file.Name())))
	s.Nil(file.Close())
	s.Nil(Remove(file.Name()))
	s.Nil(Remove(dstDir))
}

func (s *SystemHelperTestSuite) TestUnArchiveFailsForNonExistentFile() {
	srcFile := "nonexistent.zip"
	dstDir, _ := TempDir("unarchived")

	err := UnArchive(srcFile, dstDir)
	s.NotNil(err)
	s.Nil(Remove(dstDir))
}

func (s *SystemHelperTestSuite) TestArchiveSuccessfullyArchivesFiles() {
	srcFile, _ := TempFile("test")
	dstDir, _ := TempDir("archive")

	err := Write(srcFile.Name(), "test data", 0644)
	s.Nil(err)

	err = Archive([]string{srcFile.Name()}, filepath.Join(dstDir, "test.zip"))
	s.Nil(err)
	s.FileExists(filepath.Join(dstDir, "test.zip"))
	s.Nil(srcFile.Close())
	s.Nil(Remove(srcFile.Name()))
	s.Nil(Remove(dstDir))
}

func (s *SystemHelperTestSuite) TestArchiveFailsForNonExistentFiles() {
	srcFile := "nonexistent"
	dstDir, _ := TempDir("archive")

	err := Archive([]string{srcFile}, filepath.Join(dstDir, "test.zip"))
	s.NotNil(err)
	s.Nil(Remove(dstDir))
}
