package arena_convoy

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/alias"
	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"

	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_build"
	"wartank/app/lev2/arena/arena_convoy/bf_glory_find"
	"wartank/app/lev2/arena/arena_convoy/bf_glory_make"
	"wartank/app/lev2/arena/arena_convoy/bf_glory_take"
)

/*
	Объект конвоя в ангаре
*/

// АренаКонвой -- объект конвоя в ангаре
type АренаКонвой struct {
	ИАренаСтроение
	конт ILocalCtx
}

// НовКонвой -- возвращает новый *Convoy
func НовКонвой(конт ILocalCtx) *АренаКонвой {
	сам := &АренаКонвой{
		конт: конт,
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Конвой",
		СтрКонтроль_: `<title>Конвой</title>`,
		СтрУрл_:      "https://wartank.ru/convoy",
	}
	сам.ИАренаСтроение = arena_build.НовАренаСтроение(конт, аренаКонфиг)
	конт.Set("конвой", сам, "Арена конвоя бота")
	_ = ИАренаКонвой(сам)
	return сам
}

func (сам *АренаКонвой) Пуск() {
	сам.ИАренаСтроение.Пуск()
	if сам.Состояние().Получ() == cons.РежимНеСуществует {
		сам.Состояние().Уст(cons.РежимПостроено)
	}
	bf_glory_find.СлаваНайти(сам.конт)
	bf_glory_make.СлаваБой(сам.конт)
	bf_glory_take.СлаваВзять(сам.конт)
	сам.обновитьВремя()
}

// Обновляет оставшееся время конвоя
func (сам *АренаКонвой) обновитьВремя() {
	// Время подходит надо обновляться
	сам.Обновить()
	сам.ОбратВремяУст("20")
	// Найти строку с упоминанием оставшегося времени конвоя
	lstConvoy := сам.СписПолучить()
	var (
		strLastTime string
		еслиНайдено bool
		isMask      bool
	)
	for _, lastTime := range lstConvoy {
		if strings.Contains(lastTime, `До следующего конвоя: `) {
			strLastTime = lastTime
			еслиНайдено = true
			break
		}
		if strings.Contains(lastTime, `Полная маскировка через `) {
			strLastTime = lastTime
			isMask = true
			break
		}
		// <div class="bot"><a class="simple-but border red" w:id="startFight" href="convoy?7-1.ILinkListener-root-startFight"><span><span>В БОЙ!</span></span></a></div>
		if strings.Contains(lastTime, `ILinkListener-root-startFight`) {
			return
		}
		if strings.Contains(lastTime, `ILinkListener-root-findEnemy`) {
			return
		}
		// <div class="bot"><a class="simple-but border" w:id="startMasking" href="convoy?12-1.ILinkListener-root-startMasking"><span><span>В БОЙ!</span></span></a></div>
		if strings.Contains(lastTime, `ILinkListener-root-startMasking`) {
			return
		}
	}
	switch {
	case еслиНайдено: // Большая пауза между конвоями
		// Ждём окончания ожидания конвоя
		lstTime := strings.Split(strLastTime, `До следующего конвоя: `)
		strLastTime = lstTime[1]
		сам.ОбратВремяУст(АВремя(strLastTime))
	case isMask: // Если маскировка между конвоями
		// Ждём окончания ожидания конвоя
		lstTime := strings.Split(strLastTime, `Полная маскировка через `)
		strLastTime = lstTime[1]
		сам.ОбратВремяУст(АВремя(strLastTime))
	}
}
