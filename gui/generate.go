package gui

import (
	"fmt"
	"kituxlsx/process"
	"kituxlsx/reductor"
	"kituxlsx/ucexcel"
	"kituxlsx/utility"
	"path/filepath"
	"strings"
)

func (g *gui) generate() error {
	model := reductor.Instance().Model("")
	if model.Order == "" {
		model.Order = "1"
		reductor.Instance().SetModel("", model)
	}
	process, err := process.New(g)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	err = process.ReadXlsx(model.File)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	msg := fmt.Sprintf("обработано %d марок", len(process.Records))
	logMessage(msg, g.textRich)
	err = process.GeneratePallet()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	msg = fmt.Sprintf("сгенерировано %d палет по %d шт", len(process.PalletOrder), model.PerPallet)
	logMessage(msg, g.textRich)
	name := filepath.Base(model.File)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	excel := ucexcel.New(name)
	if err := excel.Open(); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := excel.ReportList(process); err != nil {
		return fmt.Errorf("%w", err)
	}
	outName, err := utility.DialogSaveFile(utility.Excel, "pallet_"+name+".xlsx", ".")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := excel.ToFileName(outName); err != nil {
		logErrMessage(err.Error(), g.textRich)
		return fmt.Errorf("%w", err)
	}
	modelFinal := reductor.Instance().Model("")
	modelFinal.StartNumberSSCC = modelFinal.LastSSCC + 1
	if err := modelFinal.Sync(g); err != nil {
		logErrMessage(err.Error(), g.textRich)
		return fmt.Errorf("%w", err)
	}
	reductor.Instance().SetModel("", modelFinal)
	return nil
}
