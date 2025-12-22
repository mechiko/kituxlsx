package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"kituxlsx/app"
	"kituxlsx/config"
	"kituxlsx/gui"
	"kituxlsx/licenser"
	"kituxlsx/reductor"
	"kituxlsx/zaplog"

	"kituxlsx/utility"

	"go.uber.org/zap"
)

// если local true то папка создается локально
var local = flag.Bool("local", false, "")

var fileExe string
var dir string

const StartNumber = 7
const CountByPallet = 24

func init() {
	flag.Parse()
	fileExe = os.Args[0]
	var err error
	dir, err = filepath.Abs(filepath.Dir(fileExe))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get absolute path: %v\n", err)
		os.Exit(1)
	}
	if err := os.Chdir(dir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to change directory: %v\n", err)
		os.Exit(1)
	}
}

func errMessageExit(loger *zap.SugaredLogger, title string, err error) {
	if loger != nil {
		loger.Errorf("%s %v", title, err)
	}
	utility.MessageBox(title, err.Error())
	os.Exit(-1)
}

func main() {
	cfg, err := config.New("", !*local)
	if err != nil {
		errMessageExit(nil, "ошибка конфигурации", err)
	}

	var logsOutConfig = map[string][]string{
		"logger": {"stdout", filepath.Join(cfg.LogPath(), config.Name)},
	}
	zl, err := zaplog.New(logsOutConfig, true)
	if err != nil {
		errMessageExit(nil, "ошибка создания логера", err)
	}

	lg, err := zl.GetLogger("logger")
	if err != nil {
		errMessageExit(nil, "ошибка получения логера", err)
	}
	loger := lg.Sugar()
	loger.Debug("zaplog started")
	loger.Infof("mode = %s", config.Mode)
	if cfg.Warning() != "" {
		loger.Infof("pkg:config warning %s", cfg.Warning())
	}
	_, err = licenser.New(licenser.MAC, "")
	if err != nil {
		errMessageExit(loger, "ошибка лицензии", err)
	}

	// создаем приложение с опциями из конфига и логером основным
	app := app.New(cfg, loger, dir)

	model := reductor.Model{}
	model.Read(app)
	if model.StartNumberSSCC < 0 {
		model.StartNumberSSCC = 0
		if err := model.Sync(app); err != nil {
			errMessageExit(loger, "Ошибки записи файла конфигурации", err)
		}
	}
	if model.PerPallet < 0 {
		model.PerPallet = 1
		if err := model.Sync(app); err != nil {
			errMessageExit(loger, "Ошибки записи файла конфигурации", err)
		}
	}

	// создаем редуктор с новой моделью
	reductor.New(model, app.Logger())
	_, _ = gui.New(app).StartDialog(app)
	// guitk9.New(k, app).Run()
}
