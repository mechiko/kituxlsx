package process

import (
	"kituxlsx/reductor"
)

func (k *Process) ReadOrder() (err error) {

	model := reductor.Instance().Model("")
	_ = reductor.Instance().SetModel("", model)
	return nil
}
