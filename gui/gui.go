package gui

import (
	"fmt"
	"kituxlsx/domain"
	"kituxlsx/reductor"
	"kituxlsx/utility"
	"path/filepath"

	"github.com/mechiko/walk"
	dcl "github.com/mechiko/walk/declarative"
)

type gui struct {
	domain.Apper
	dlg                *walk.Dialog
	page               *walk.Composite
	acceptPB, cancelPB *walk.PushButton
	selectFile         *walk.PushButton
	textRich           *RichEdit
	fileName           *walk.Label
	perPallet          *walk.NumberEdit
	startPallet        *walk.NumberEdit
	order              *walk.TextEdit
	ssccPrefix         *walk.TextEdit
}

func New(app domain.Apper) *gui {
	g := &gui{
		Apper: app,
	}
	return g
}

func (g *gui) StartDialog(app domain.Apper) (out string, err error) {
	icon, err := walk.Resources.Icon("3")
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	model := reductor.Instance().Model("")
	if err := (dcl.Dialog{
		AssignTo:      &g.dlg,
		Title:         "Агрегация для Тайрос",
		Size:          dcl.Size{Width: 700, Height: 400},
		Icon:          icon,
		Layout:        dcl.VBox{Spacing: 10, Margins: dcl.Margins{Left: 5, Top: 5, Right: 5, Bottom: 5}},
		DefaultButton: &g.acceptPB,
		CancelButton:  &g.cancelPB,
		Children: []dcl.Widget{
			dcl.Composite{
				Layout:   dcl.VBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 0, Top: 0, Right: 0, Bottom: 0}},
				Border:   true,
				AssignTo: &g.page,
				MinSize:  dcl.Size{Height: 200},
				Children: []dcl.Widget{},
			},
			dcl.Composite{
				Layout:  dcl.HBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 0, Top: 0, Right: 0, Bottom: 0}},
				Border:  false,
				MinSize: dcl.Size{Width: 500},
				Children: []dcl.Widget{
					dcl.Label{
						Text: "Файл для обработки:",
					},
					dcl.Label{
						AssignTo: &g.fileName,
						Text:     "",
					},
					dcl.PushButton{
						AssignTo: &g.selectFile,
						Text:     "Выбрать",
						OnClicked: func() {
							model := reductor.Instance().Model("")
							if model.IsProcess {
								utility.MessageBox("ошибка", "запущен процесс обработки. дождитесь его завершения")
								return
							}
							fName, err := utility.DialogOpenFile([]utility.FileType{utility.Excel}, "", ".")
							if err != nil {
								model.File = ""
								logErrMessage(err.Error(), g.textRich)
								return
							}
							model.File = fName
							reductor.Instance().SetModel("", model)
							g.fileName.SetText(filepath.Base(fName))
							logMessage(filepath.Base(fName), g.textRich)
						},
					},
					dcl.HSpacer{},
				},
			},
			dcl.Composite{
				Layout:  dcl.HBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 0, Top: 0, Right: 0, Bottom: 0}},
				Border:  false,
				MinSize: dcl.Size{Width: 500},
				Children: []dcl.Widget{
					dcl.Label{
						Text: "В палете:",
					},
					dcl.NumberEdit{
						AssignTo: &g.perPallet,
						Value:    0,
					},
					dcl.Label{
						Text: "SSCC Префикс:",
					},
					dcl.TextEdit{
						AssignTo: &g.ssccPrefix,
						Text:     model.PrefixSSCC,
					},
					dcl.HSpacer{},
				},
			},
			dcl.Composite{
				Layout:  dcl.HBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 0, Top: 0, Right: 0, Bottom: 0}},
				Border:  false,
				MinSize: dcl.Size{Width: 500},
				Children: []dcl.Widget{
					dcl.Label{
						Text: "Номер палеты :",
					},
					dcl.NumberEdit{
						AssignTo: &g.startPallet,
						Value:    0,
					},
					dcl.Label{
						Text: "Заказ :",
					},
					dcl.TextEdit{
						AssignTo: &g.order,
						Text:     model.Order,
					},
					dcl.HSpacer{},
				},
			},
			dcl.Composite{
				Border: false,
				Layout: dcl.HBox{MarginsZero: true, SpacingZero: true, Margins: dcl.Margins{Left: 5, Top: 5, Right: 5, Bottom: 0}},
				Children: []dcl.Widget{
					dcl.PushButton{
						AssignTo: &g.acceptPB,
						Text:     "Обработать",
						OnClicked: func() {
							model := reductor.Instance().Model("")
							if model.IsProcess {
								utility.MessageBox("ошибка", "запущен процесс обработки. дождитесь его завершения")
								return
							}
							if model.File == "" {
								utility.MessageBox("ошибка", "выберите файл")
								return
							}
							if !utility.PathOrFileExists(model.File) {
								utility.MessageBox("ошибка", "файл не найден "+model.File)
								return
							}
							model.Order = g.order.Text()
							model.PerPallet = int(g.perPallet.Value())
							model.StartNumberSSCC = int(g.startPallet.Value())
							model.PrefixSSCC = g.ssccPrefix.Text()
							reductor.Instance().SetModel("", model)
							g.start()
							go func() {
								if err := g.generate(); err != nil {
									logErrMessage(err.Error(), g.textRich)
								}
								g.finish()
							}()
							// dlg.Accept()
						},
					},
					dcl.PushButton{
						AssignTo: &g.cancelPB,
						Text:     "Выход",
						OnClicked: func() {
							model := reductor.Instance().Model("")
							if model.IsProcess {
								utility.MessageBox("ошибка", "запущен процесс обработки. дождитесь его завершения")
								return
							}
							g.dlg.Cancel()
						},
					},
					dcl.HSpacer{},
				},
			},
			dcl.VSpacer{},
		},
	}).Create(nil); err != nil {
		return "", fmt.Errorf("%w", err)
	}
	g.dlg.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		model := reductor.Instance().Model("")
		if model.IsProcess {
			*canceled = true
			utility.MessageBox("ошибка", "запущен процесс обработки. дождитесь его завершения")
			return
		}
		*canceled = false
	})
	g.textRich, err = NewRichEdit(g.page)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	if g.textRich == nil {
		return "", fmt.Errorf("nil richtext")
	}
	g.dlg.SetBounds(walk.Rectangle{
		X:     300,
		Y:     300,
		Width: 400,
		// Height: 200,
	})
	g.order.SetText(model.Order)
	g.perPallet.SetValue(float64(model.PerPallet))
	g.startPallet.SetValue(float64(model.StartNumberSSCC))
	logMessage("выберите файл", g.textRich)
	if ret := g.dlg.Run(); ret != 1 {
		return "", fmt.Errorf("dialog return %d", ret)
	}
	return out, nil
}
