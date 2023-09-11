package tests

import (
	"github.com/goravel/framework/testing"

	"panel/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
