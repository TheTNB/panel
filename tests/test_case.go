package tests

import (
	"github.com/goravel/framework/testing"

	"github.com/TheTNB/panel/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
