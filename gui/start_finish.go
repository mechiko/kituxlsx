package gui

import "kituxlsx/reductor"

func (g *gui) start() {
	model := reductor.Instance().Model("")
	model.IsProcess = true
	reductor.Instance().SetModel("", model)
	g.dlg.Synchronize(func() {
		g.acceptPB.SetEnabled(false)
		g.cancelPB.SetEnabled(false)
		g.selectFile.SetEnabled(false)
	})
}

func (g *gui) finish() {
	modelFinal := reductor.Instance().Model("")
	modelFinal.IsProcess = false
	reductor.Instance().SetModel("", modelFinal)
	g.dlg.Synchronize(func() {
		g.order.SetText(modelFinal.Order)
		g.perPallet.SetValue(float64(modelFinal.PerPallet))
		g.startPallet.SetValue(float64(modelFinal.StartNumberSSCC))
		g.acceptPB.SetEnabled(true)
		g.cancelPB.SetEnabled(true)
		g.selectFile.SetEnabled(true)
	})
}
