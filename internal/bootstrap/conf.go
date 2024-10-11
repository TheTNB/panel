package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/TheTNB/panel/internal/app"
)

func initConf() {
	executable, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("failed to get executable: %v", err))
	}
	res, err := filepath.EvalSymlinks(filepath.Dir(executable))
	if err != nil {
		panic(fmt.Sprintf("failed to get executable path: %v", err))
	}
	if isTesting() || isAir() || isDirectlyRun() {
		res, err = os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("failed to get working directory: %v", err))
		}
	}

	app.Conf = koanf.New(".")
	if err = app.Conf.Load(file.Provider(filepath.Join(res, "config/config.yml")), yaml.Parser()); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
}

func initGlobal() {
	app.Root = app.Conf.MustString("app.root")
	app.Version = app.Conf.MustString("app.version")
	app.Locale = app.Conf.MustString("app.locale")

	// 初始化时区
	loc, err := time.LoadLocation(app.Conf.MustString("app.timezone"))
	if err != nil {
		panic(fmt.Sprintf("failed to load timezone: %v", err))
	}
	time.Local = loc
}

// isTesting checks if the application is running in testing mode.
func isTesting() bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, "-test.") {
			return true
		}
	}

	return false
}

// isAir checks if the application is running using Air.
func isAir() bool {
	for _, arg := range os.Args {
		if strings.Contains(filepath.ToSlash(arg), "/storage/temp") {
			return true
		}
	}

	return false
}

// isDirectlyRun checks if the application is running using go run.
func isDirectlyRun() bool {
	executable, _ := os.Executable()
	return strings.Contains(filepath.Base(executable), os.TempDir()) ||
		(strings.Contains(filepath.ToSlash(executable), "/var/folders") && strings.Contains(filepath.ToSlash(executable), "/T/go-build")) // macOS
}
