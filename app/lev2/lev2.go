// package lev2 -- слой арен
package lev2

import (
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev2/arena/arena_angar"
	"wartank/app/lev2/arena/arena_arsenal"
	"wartank/app/lev2/arena/arena_bank"
	"wartank/app/lev2/arena/arena_base"
	"wartank/app/lev2/arena/arena_battle"
	"wartank/app/lev2/arena/arena_convoy"
	"wartank/app/lev2/arena/arena_division"
	"wartank/app/lev2/arena/arena_fuel_duel"
	"wartank/app/lev2/arena/arena_fuel_storage"
	"wartank/app/lev2/arena/arena_market"
	"wartank/app/lev2/arena/arena_masters"
	"wartank/app/lev2/arena/arena_medal"
	"wartank/app/lev2/arena/arena_mine"
	"wartank/app/lev2/arena/arena_missions"
	"wartank/app/lev2/arena/arena_polygon"
)

// НовАренаБитваМастеров -- возвращает новую арену битву мастеров PVP
func НовАренаБитваМастеров(конт ILocalCtx) ИАренаСтроение {
	битва := arena_masters.НовБитваМастеров(конт)
	return битва
}

// НовСражение -- возвращает новую арену сражения PVE
func НовСражение(конт ILocalCtx) ИАренаСтроение {
	сражение := arena_battle.НовСражение(конт)
	return сражение
}

// НовАренаРынок -- возвращает новую арену рынка
func НовАренаРынок(конт ILocalCtx) ИАрена {
	рынок := arena_market.НовРынок(конт)
	return рынок
}

// НовАренаМедаль -- возвращает новую арену медалей
func НовАренаМедали(конт ILocalCtx) ИАрена {
	медали := arena_medal.НовАренаМедали(конт)
	return медали
}

// НовБойТопливо -- возвращает новую арену боя за топливо
func НовБойТопливо(конт ILocalCtx) ИАренаСтроение {
	арена := arena_fuel_duel.НовАренаТопливоДуэль(конт)
	return арена
}

// НовКонвой -- возвращает новый конвой
func НовКонвой(конт ILocalCtx) ИАренаКонвой {
	конвой := arena_convoy.НовКонвой(конт)
	return конвой
}

// НовБанк -- возвращает новый банк бота
func НовБанк(конт ILocalCtx) ИАренаБанк {
	банк := arena_bank.НовБанк(конт)
	return банк
}

// НовМиссииПростые -- возвращает арену новых простых миссий
func НовМиссииПростые(конт ILocalCtx) ИАренаМиссииПростые {
	миссии := arena_missions.НовМиссии(конт)
	return миссии
}

// НовАренаБак -- возвращает арену база топлива
func НовАренаБак(конт ILocalCtx) ИАренаБак {
	бак := arena_fuel_storage.НовАренаБак(конт)
	return бак
}

// НовАнгар -- возвращает новый ангар
func НовАнгар(конт ILocalCtx) ИАренаАнгар {
	ангар := arena_angar.НовАнгар(конт)
	return ангар
}

// НовПолигон -- возвращает новый полигон
func НовПолигон(конт ILocalCtx) ИАренаПолигон {
	полигон := arena_polygon.НовПолигон(конт)
	return полигон
}

// НовАрсенал -- возвращает новый арсенал
func НовАрсенал(конт ILocalCtx) ИАренаАрсенал {
	арсенал := arena_arsenal.НовАрсенал(конт)
	return арсенал
}

// НовБаза -- возвращает новую базу
func НовБаза(конт ILocalCtx) ИАренаБаза {
	база := arena_base.НовБаза(конт)
	return база
}

// НовШахта -- возвращает новую шахту
func НовШахта(конт ILocalCtx) ИАренаШахта {
	шахта := arena_mine.НовШахта(конт)
	return шахта
}

// ЗапустиДивизию -- запускает бой дивизий в фоновых горутинах (без остановки -- контекст управляет жизнью)
func ЗапустиДивизию(конт ILocalCtx) {
	arena_division.НовДивизия(конт)
}
