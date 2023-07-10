package helpers

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

func (s *SystemHelperTestSuite) TestWriteFile() {
	filePath := "/tmp/testfile"
	defer os.Remove(filePath)

	s.True(WriteFile(filePath, "test data", 0644))
	s.FileExists(filePath)

	content, _ := os.ReadFile(filePath)
	s.Equal("test data", string(content))
}

func (s *SystemHelperTestSuite) TestReadFile() {
	filePath := "/tmp/testfile"
	defer os.Remove(filePath)

	err := os.WriteFile(filePath, []byte("test data"), 0644)
	s.Nil(err)

	s.Equal("test data", ReadFile(filePath))
}

func (s *SystemHelperTestSuite) TestRemoveFile() {
	filePath := "/tmp/testfile"

	err := os.WriteFile(filePath, []byte("test data"), 0644)
	s.Nil(err)

	s.True(RemoveFile(filePath))
}

func (s *SystemHelperTestSuite) TestExecShell() {
	s.Equal("test\n", ExecShell("echo 'test'"))
}

func (s *SystemHelperTestSuite) TestExecShellAsync() {
	command := "echo 'test' > /tmp/testfile"
	defer os.Remove("/tmp/testfile")

	ExecShellAsync(command)

	time.Sleep(time.Second)

	content, _ := os.ReadFile("/tmp/testfile")
	s.Equal("test\n", string(content))
}

func (s *SystemHelperTestSuite) TestMkdir() {
	dirPath := "/tmp/testdir"
	defer os.RemoveAll(dirPath)

	s.True(Mkdir(dirPath, 0755))
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
	groups, err := currentUser.GroupIds()
	s.Nil(err)

	s.True(Chown(filePath, currentUser.Username, groups[0]))
}
