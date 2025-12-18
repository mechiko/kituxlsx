package process

import (
	"kituxlsx/domain"

	"github.com/mechiko/utility"
)

// const startSSCC = "1462709225" // gs1 rus id zapivkom для памяти запивком

type Process struct {
	domain.Apper
	Sscc     []string
	Cis      []*utility.CisInfo
	Pallet   map[string][]*utility.CisInfo
	warnings []string
	errors   []string
}

func New(app domain.Apper) (*Process, error) {
	k := &Process{
		Apper:    app,
		Pallet:   make(map[string][]*utility.CisInfo),
		warnings: make([]string, 0),
		errors:   make([]string, 0),
		Sscc:     make([]string, 0),
		Cis:      make([]*utility.CisInfo, 0),
	}
	return k, nil
}

func (k *Process) AddWarn(warn string) {
	k.warnings = append(k.warnings, warn)
}

func (k *Process) Warnings() []string {
	return k.warnings
}

func (k *Process) AddError(err string) {
	k.errors = append(k.errors, err)
}

func (k *Process) Errors() []string {
	return k.errors
}

func (k *Process) ResetPalletMap() {
	for key := range k.Pallet {
		delete(k.Pallet, key)
	}
}

func (k *Process) Reset() {
	k.ResetPalletMap()
	k.Cis = make([]*utility.CisInfo, 0)
	k.Sscc = make([]string, 0)
	k.errors = make([]string, 0)
	k.warnings = make([]string, 0)
}
