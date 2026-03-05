package arena_bank

import (
	"log"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev1"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_bank/bank_mode"
	"wartank/app/lev2/arena/arena_bank/bf_bank_build"
	"wartank/app/lev2/arena/arena_bank/bf_bank_prod"
	"wartank/app/lev2/arena/arena_bank/bf_bank_take"
	"wartank/app/lev2/arena/arena_bank/bf_bank_upgrade"
	"wartank/app/lev2/arena/arena_bank/bf_bank_upgrade_fast"
	"wartank/app/lev2/arena/arena_build"
)

/*
	Предоставляет объект банка на базе
*/

// Банк -- объект банка на базе
type АренаБанк struct {
	ИАренаСтроение
	конт       ILocalCtx
	сереброБот ИСтатПарам
	режим1     *bank_mode.BankMode // 1 режим работы на выбор
	режим2     *bank_mode.BankMode // 2 режим работы на выбор
}

// НовБанк -- возвращает новый арену банка
func НовБанк(конт ILocalCtx) ИАренаБанк {

	сам := &АренаБанк{
		конт:       конт,
		сереброБот: lev1.НовСтатПарам("серебро бота"),
		режим1:     bank_mode.NewBankMode(конт),
		режим2:     bank_mode.NewBankMode(конт),
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Банк",
		СтрКонтроль_: `<span class="green2">Серебро</span><br/>`,
		СтрУрл_:      "https://wartank.ru/production/Bank",
	}
	сам.ИАренаСтроение = arena_build.НовАренаСтроение(конт, аренаКонфиг)
	конт.Set("банк", сам, "Арена банка бота")
	return сам
}

func (сам *АренаБанк) Пуск() {
	сам.ИАренаСтроение.Пуск()
	bf_bank_build.БанкПостроить(сам.конт)
	bf_bank_upgrade.БанкАпгрейд(сам.конт)
	bf_bank_upgrade_fast.БанкАпгрейдБесплатно(сам.конт)
	bf_bank_take.БанкЗабрать(сам.конт)
	bf_bank_prod.СереброПроизводить(сам.конт)
}

// Проверяет необходимость постройки полигона
func (сам *АренаБанк) построитьУлучшить() bool {
	var списБанк []string

	{ // Зайти на страницу постройки
		// https://wartank.ru/building-upgrade/Bank
		списБанк = сам.Сеть().ВебВоркер().Получ("https://wartank.ru/building-upgrade/Bank")
		стрСсылка := ""
		еслиНайти := false
		// <a class="simple-but border mb5" href="Bank?192-1.ILinkListener-upgradeLink-link">
		for _, стрСсылка = range списБанк {
			if strings.Contains(стрСсылка, `href="Bank?`) {
				еслиНайти = true
				break
			}
		}
		if !еслиНайти { // Время полигона вышло
			return false
		}
		_ссылка := strings.TrimPrefix(стрСсылка, `<a class="simple-but border mb5" href="`)
		_ссылка = strings.TrimSuffix(_ссылка, `">`)
		ссылка := "https://wartank.ru/building-upgrade/" + _ссылка
		// https://wartank.ru/building-upgrade/Bank?162-1.ILinkListener-upgradeLink-link
		списБанк = сам.Сеть().ВебВоркер().Получ(ссылка)
	}
	{ // Заказать постройку
		// https://wartank.ru/building-upgrade/Bank
		стрСсылка := ""
		еслиНайти := false
		// <a class="simple-but border mb5" href="Bank?163-1.ILinkListener-upgradeLink-link">
		for _, стрСсылка = range списБанк {
			if strings.Contains(стрСсылка, `href="Bank?`) {
				еслиНайти = true
				break
			}
		}
		if еслиНайти { // Время полигона вышло
			_ссылка := strings.TrimPrefix(стрСсылка, `<a class="simple-but border mb5" href="`)
			_ссылка = strings.TrimSuffix(_ссылка, `">`)
			ссылка := "https://wartank.ru/building-upgrade/" + _ссылка
			// https://wartank.ru/building-upgrade/Bank?162-1.ILinkListener-upgradeLink-link
			списБанк = сам.Сеть().ВебВоркер().Получ(ссылка)
		}
	}
	{ // подтверждение постройки
		// <a class="simple-but border w50 mXa mb10" w:id="confirmLink" href="../wicket/page?187-1.ILinkListener-confirmLink"><span><span>да, подтверждаю</span></span></a>
		стрСсылка := ""
		еслиНайти := false
		for _, стрСсылка = range списБанк {
			if strings.Contains(стрСсылка, `.ILinkListener-confirmLink`) {
				еслиНайти = true
				break
			}
		}
		if !еслиНайти { // Время полигона вышло
			return false
		}
		_ссылка := strings.TrimPrefix(стрСсылка, `<a class="simple-but border w50 mXa mb10" w:id="confirmLink" href="../`)
		_ссылка = strings.TrimSuffix(_ссылка, "\"><span><span>да, подтверждаю</span></span></a>")
		ссылка := "https://wartank.ru/" + _ссылка
		// https://wartank.ru/wicket/page?135-1.ILinkListener-confirmLink
		_ = сам.Сеть().ВебВоркер().Получ(ссылка)
	}
	log.Printf("Банк.построитьПровер(): построен успешно\n")
	return true
}

// РежимРаботы2 -- возвращает объект режима2
func (сам *АренаБанк) РежимРаботы2() ИБанкРежим {
	return сам.режим2
}

// РежимРаботы1 -- возвращает объект режима1
func (сам *АренаБанк) РежимРаботы1() ИБанкРежим {
	return сам.режим1
}

// СереброБот -- возвращает серебро от бота
func (сам *АренаБанк) СереброБот() ИСтатПарам {
	return сам.сереброБот
}
