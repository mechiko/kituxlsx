package gui

import (
	"fmt"
	"kituxlsx/domain"
	"kituxlsx/process"
	"kituxlsx/reductor"
	"kituxlsx/ucexcel"
	"kituxlsx/utility"
	"path/filepath"
	"strings"
)

func generate(app domain.Apper, tr *RichEdit) {
	model := reductor.Instance().Model("")
	if model.Order == "" {
		model.Order = "1"
		reductor.Instance().SetModel("", model)
	}
	fileName, err := utility.DialogOpenFile([]utility.FileType{utility.Excel}, "", ".")
	if err != nil {
		logErrMessage(err.Error(), tr)
		// utility.MessageBox32("ошибка", err.Error())
		return
	}
	app.Logger().Info(fileName)
	process, err := process.New(app)
	if err != nil {
		logErrMessage(err.Error(), tr)
		// utility.MessageBox32("ошибка", err.Error())
		return
	}
	err = process.ReadXlsx(fileName)
	if err != nil {
		logErrMessage(err.Error(), tr)
		return
	}
	msg := fmt.Sprintf("обработано %d марок", len(process.Records))
	logMessage(msg, tr)
	err = process.GeneratePallet()
	if err != nil {
		logErrMessage(err.Error(), tr)
		return
	}
	msg = fmt.Sprintf("сгенерировано %d палет по %d шт", len(process.PalletOrder), model.PerPallet)
	logMessage(msg, tr)
	name := filepath.Base(fileName)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	excel := ucexcel.New(name)
	if err := excel.Open(); err != nil {
		logErrMessage(err.Error(), tr)
		return
	}
	if err := excel.ReportList(process); err != nil {
		logErrMessage(err.Error(), tr)
		return
	}
	outName, err := utility.DialogSaveFile(utility.Excel, "pallet_"+name+".xlsx", ".")
	if err != nil {
		logErrMessage(err.Error(), tr)
		// utility.MessageBox32("ошибка", err.Error())
		return
	}
	if err := excel.ToFileName(outName); err != nil {
		logErrMessage(err.Error(), tr)
		return
	}

}
