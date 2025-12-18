package guitk9

import (
	"kituxlsx/guitk9/dconfig"
	"kituxlsx/reductor"
)

func (a *GuiApp) onConfig() {
	model := reductor.Instance().Model("")
	data := dconfig.ConfigDialogData{
		PrefixSSCC: model.PrefixSSCC,
	}
	dlg := dconfig.NewConfigDialog(&data)
	dlg.ShowModal()
	if data.Ok {
		model.PrefixSSCC = data.PrefixSSCC
		if err := model.Sync(a); err != nil {
			a.Logger().Errorf("диалог onConfig синхронизация модели %v", err)
		}
	}
	reductor.Instance().SetModel("", model)
}
