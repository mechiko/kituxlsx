package process

import (
	"fmt"
	"kituxlsx/domain"
	"kituxlsx/reductor"

	"kituxlsx/utility"
)

const boxPrefix = "box"

// start начальный номер SSCC палетты
// count количество в одной палетте
func (k *Process) GeneratePallet() error {
	indexPallet := 0
	model := reductor.Instance().Model("")
	startSSCC := model.StartNumberSSCC
	lastSSCC := model.StartNumberSSCC
	box := 0
	for {
		box++
		cis := nextRecords(k.Records, indexPallet, model.PerPallet)
		if len(cis) == 0 {
			// больше нет км
			// выходим без расчета номера палеты и последняя палета так и останется последней сгенерированной
			break
		}
		lastSSCC = startSSCC + indexPallet
		if len(model.PrefixSSCC) < 7 {
			model.PrefixSSCC = fmt.Sprintf("%07s", model.PrefixSSCC)
		}
		pallet, err := utility.GenerateSSCC(lastSSCC, model.PrefixSSCC)
		if err != nil {
			return fmt.Errorf("generate sscc error %w", err)
		}
		if _, ok := k.Pallet[pallet]; ok {
			return fmt.Errorf("паллета %s уже сгенерирована прежде в обработке", pallet)
		}
		k.Sscc = append(k.Sscc, pallet)
		boxStr := fmt.Sprintf("%s_%03d", boxPrefix, box)
		for _, cc := range cis {
			cc.Box = boxStr
		}
		k.Pallet[pallet] = cis
		k.PalletOrder = append(k.PalletOrder, pallet)
		if len(cis) < model.PerPallet {
			// последняя не полная палетта
			break
		}
		indexPallet++
	}
	model.LastSSCC = lastSSCC
	_ = reductor.Instance().SetModel("", model)
	return nil
}

// получить следующие км
// i номер группы по count штук
// если размер массива меньше count значит последний
// елси размер массива 0 значит больше нет
func (k *Process) nextRecords(i int, count int) (out []*domain.Record) {
	lenCis := len(k.Records)
	out = make([]*domain.Record, 0)
	first := i * count // первая км в цикле 24 шт
	for i := 0; i < count; i++ {
		index := i + first
		if (index + 1) > lenCis {
			return out
		}
		out = append(out, k.Records[index])
	}
	return out
}

// index from 0 startIndex 0
// nextRecords returns a batch of records starting from startIndex
// Returns empty slice when no more records are available
func nextRecords(arr []*domain.Record, index int, count int) []*domain.Record {
	startIndex := index * count
	if startIndex >= len(arr) {
		return []*domain.Record{}
	}
	endIndex := startIndex + count
	// если последний индекс больше длины массива укорачиваем до размера массива
	if endIndex > len(arr) {
		endIndex = len(arr)
	}
	return arr[startIndex:endIndex]
}
