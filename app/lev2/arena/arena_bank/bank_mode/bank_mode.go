package bank_mode

import (
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev1"
)

/*
	Объект допустимого режима банка.
*/

// BankMode -- объект допустимого режима банка
type BankMode struct {
	серебро   ИСтатПарам
	timeCount string
}

// NewBankMode -- возвращает новый *BankMode
func NewBankMode(конт ILocalCtx) *BankMode {
	return &BankMode{
		серебро: lev1.НовСтатПарам("серебро"),
	}
}

// Серебро -- возвращает объект серебра режима
func (сам *BankMode) Серебро() ИСтатПарам {
	return сам.серебро
}

// ВремяСделать -- возвращает временя производства режима
func (сам *BankMode) ВремяСделать() string {
	return сам.timeCount
}

// ВремяСделатьУст -- устанавливает время производства режима
func (сам *BankMode) ВремяСделатьУст(val string) {
	сам.timeCount = val
}
