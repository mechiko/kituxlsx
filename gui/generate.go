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

func (g *gui) generate() error {
	model := reductor.Instance().Model("")
	if model.Order == "" {
		model.Order = "1"
		reductor.Instance().SetModel("", model)
	}
	prc, err := process.New(g)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	err = prc.ReadXlsx(model.File)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	msg := fmt.Sprintf("обработано %d марок", len(prc.Records))
	logMessage(msg, g.textRich)
	err = prc.GeneratePallet()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	msg = fmt.Sprintf("сгенерировано %d палет по %d шт", len(prc.PalletOrder), model.PerPallet)
	logMessage(msg, g.textRich)
	name := filepath.Base(model.File)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	excel := ucexcel.New(name)
	if err := excel.Open(); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := excel.ReportList(prc); err != nil {
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
	pals := make([]*domain.Pallete, 0)
	for _, pal := range prc.PalletOrder {
		records := prc.Pallet[pal]
		if len(records) > 0 {
			plt := &domain.Pallete{
				Code:   pal,
				Volume: "объем 330мл",
				Name:   "BIA LON THABREW крепостью 4,8%",
				Box:    records[0].Box,
			}
			pals = append(pals, plt)
		}
	}
	if err := g.pdf(outName, pals); err != nil {
		logErrMessage(err.Error(), g.textRich)
		return fmt.Errorf("%w", err)
	}
	return nil
}
