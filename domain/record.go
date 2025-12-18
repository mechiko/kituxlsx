package domain

import (
	"fmt"
	"strings"

	"kituxlsx/utility"
)

type Record struct {
	Cis    *utility.CisInfo
	Gtin   string
	Name   string
	Serial string
	Box    string
}

// вставляем GS в файлах от криницы, там он опущен...
// это только для пива!!!
func NewRecord(row []string) (*Record, error) {
	if len(row) < 3 {
		return nil, fmt.Errorf("записей меньше 3")
	}
	s := row[0]
	cis, err := utility.ParseCisInfo(s)
	if err != nil {
		return nil, fmt.Errorf("получение КМ %w", err)
	}
	gtin := row[1]
	if gtin != cis.Gtin {
		return nil, fmt.Errorf("ошибка gtin таблицы %s не равен gtin в марке %s", row[1], cis.Gtin)
	}
	name := row[2]

	r := &Record{
		Cis:  cis,
		Gtin: gtin,
		Name: name,
	}
	return r, nil
}

// полная строка 11 ячеек
func IsRecord(row []string) bool {
	return len(row) >= 3 && strings.HasPrefix(row[0], "01")
}
