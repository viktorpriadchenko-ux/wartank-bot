// package battle_register -- бизнес-функция регистрации танк в битве мастеров
package bf_masters_register

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// СражениеРегистрация -- регистрирует бота на битву мастеров
func СражениеРегистрация(конт ILocalCtx) {
	битва := конт.Get("pvp").Val().(ИАренаСтроение)
	/*
		Здесь процесс хитрый.
		Битвы следуют с некоторым интервалом.

		При запуске бота будет состояние не создано.
		После регистрации -- платный апгрейд
		При обратном отсчёте --- ожидание
		При идущей битве -- работа.
		После окончания -- забрать и переход в ожидание.

		Целевое состояние здесь -- платный апгрейд (из не существует, построено, ожидание -- запрещено)
	*/
	if битва.Состояние().Получ() == cons.РежимНеСуществует { // БотоФерма только запущена
		битва.Состояние().Уст(cons.РежимПостроено)
	}
	еслиПостроено := битва.Состояние().Получ() == cons.РежимПостроено
	еслиОжидание := битва.Состояние().Получ() == cons.РежимОжидание // Уже были битвы
	if !(еслиОжидание || еслиПостроено) {
		return
	}
	регистрация(конт)
}

func регистрация(конт ILocalCtx) {
	битва := конт.Get("pvp").Val().(ИАренаСтроение)
	var (
		лстБитва  = битва.СписПолучить()
		стрСсылка string
	)
	defer func() {
		битва.СтрОбновить(лстБитва)
	}()
	// Найдено приглашение на участие
	if len(лстБитва) == 0 {
		битва.Обновить()
		лстБитва = битва.СписПолучить()
	}
	if len(лстБитва) < 113 { // Уже обратный отсчёт
		битва.Состояние().Уст(cons.РежимОжидание)
		return
	}
	// <a w:id="joinLink" href="pvp?45-5.ILinkListener-joinLink" class="simple-but border"><span><span>Участвовать в битве</span></span></a>
	for _, стрСсылка = range лстБитва {
		if strings.Contains(стрСсылка, `.ILinkListener-joinLink" class="simple-but border"><span><span>Участвовать в битве`) {
			break
		}
	}
	if !strings.Contains(стрСсылка, "ILinkListener-joinLink") {
		битва.Состояние().Уст(cons.РежимОжидание)
		return
	}

	// <a w:id="joinLink" href="pvp?45-5.ILinkListener-joinLink" class="simple-but border"><span><span>Участвовать в битве</span></span></a>
	_стрСсылка := strings.TrimPrefix(стрСсылка, `<a w:id="joinLink" href="`)
	_стрСсылка = strings.TrimSuffix(_стрСсылка, `" class="simple-but border"><span><span>Участвовать в битве</span></span></a>`)
	_стрСсылка = "https://wartank.ru/" + _стрСсылка
	// https://wartank.ru/pvp?45-5.ILinkListener-joinLink
	res := битва.Сеть().Get(_стрСсылка)
	res.Hassert("регистрация(): при регистрации на сражение")
	битва.Состояние().Уст(cons.РежимАпгрейдПлатный)
}
