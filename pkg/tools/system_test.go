package tools

import (
	"os"
	"os/user"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type SystemHelperTestSuite struct {
	suite.Suite
}

func TestSystemHelperTestSuite(t *testing.T) {
	suite.Run(t, &SystemHelperTestSuite{})
}

func (s *SystemHelperTestSuite) TestWrite() {
	filePath := "/tmp/testfile"
	defer os.Remove(filePath)

	s.Nil(Write(filePath, "test data", 0644))
	s.FileExists(filePath)

	content, _ := os.ReadFile(filePath)
	s.Equal("test data", string(content))
}

func (s *SystemHelperTestSuite) TestRead() {
	filePath := "/tmp/testfile"
	defer os.Remove(filePath)

	err := os.WriteFile(filePath, []byte("test data"), 0644)
	s.Nil(err)

	s.Equal("test data", Read(filePath))
}

func (s *SystemHelperTestSuite) TestRemove() {
	filePath := "/tmp/testfile"

	err := os.WriteFile(filePath, []byte("test data"), 0644)
	s.Nil(err)

	s.True(Remove(filePath))
}

func (s *SystemHelperTestSuite) TestExec() {
	s.Equal("test", Exec("echo 'test'"))
}

func (s *SystemHelperTestSuite) TestExecAsync() {
	command := "echo 'test' > /tmp/testfile"
	defer os.Remove("/tmp/testfile")

	err := ExecAsync(command)
	s.Nil(err)

	time.Sleep(time.Second)

	content, _ := os.ReadFile("/tmp/testfile")
	s.Equal("test\n", string(content))
}

func (s *SystemHelperTestSuite) TestMkdir() {
	dirPath := "/tmp/testdir"
	defer os.RemoveAll(dirPath)

	s.Nil(Mkdir(dirPath, 0755))
}

func (s *SystemHelperTestSuite) TestChmod() {
	filePath := "/tmp/testfile"
	defer os.Remove(filePath)

	err := os.WriteFile(filePath, []byte("test data"), 0644)
	s.Nil(err)

	s.True(Chmod(filePath, 0755))
}

func (s *SystemHelperTestSuite) TestChown() {
	filePath := "/tmp/testfile"
	defer os.Remove(filePath)

	err := os.WriteFile(filePath, []byte("test data"), 0644)
	s.Nil(err)

	currentUser, err := user.Current()
	s.Nil(err)
	groups, err := currentUser.GroupIds()
	s.Nil(err)

	s.True(Chown(filePath, currentUser.Username, groups[0]))
}

func (s *SystemHelperTestSuite) TestExists() {
	s.True(Exists("/tmp"))
	s.False(Exists("/tmp/123"))
}

func (s *SystemHelperTestSuite) TestEmpty() {
	s.True(Empty("/tmp/123"))
	s.False(Empty("/tmp"))
}

func (s *SystemHelperTestSuite) TestMv() {
	filePath := "/tmp/testfile"
	defer os.Remove(filePath)

	err := os.WriteFile(filePath, []byte("test data"), 0644)
	s.Nil(err)

	s.Nil(Mv(filePath, "/tmp/testfile2"))
	s.False(Exists(filePath))
}

func (s *SystemHelperTestSuite) TestCp() {
	filePath := "/tmp/testfile"
	defer os.Remove(filePath)

	err := os.WriteFile(filePath, []byte("test data"), 0644)
	s.Nil(err)

	s.Nil(Cp(filePath, "/tmp/testfile2"))
	s.True(Exists(filePath))
}

func (s *SystemHelperTestSuite) TestSize() {
	size, err := Size("/tmp/123")
	s.Equal(int64(0), size)
	s.Error(err)
}

func (s *SystemHelperTestSuite) TestFileSize() {
	size, err := FileSize("/tmp/123")
	s.Equal(int64(0), size)
	s.Error(err)
}
