package process

import (
	"kituxlsx/domain"
)

// const startSSCC = "1462709225" // gs1 rus id zapivkom для памяти запивком

type Process struct {
	domain.Apper
	Sscc        []string
	Pallet      map[string][]*domain.Record
	Records     []*domain.Record
	PalletOrder []string
}

func New(app domain.Apper) (*Process, error) {
	k := &Process{
		Apper:       app,
		Pallet:      make(map[string][]*domain.Record),
		Sscc:        make([]string, 0),
		PalletOrder: make([]string, 0),
	}
	return k, nil
}

func (k *Process) ResetPalletMap() {
	for key := range k.Pallet {
		delete(k.Pallet, key)
	}
}

func (k *Process) Reset() {
	k.ResetPalletMap()
	k.Records = make([]*domain.Record, 0)
	k.Sscc = make([]string, 0)
}
