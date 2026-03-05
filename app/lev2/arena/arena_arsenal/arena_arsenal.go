// package arena_arsenal -- объект оружейной на базе
package arena_arsenal

import (
	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev1"
	"wartank/app/lev1/web_log"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_arsenal/bf_ammo_make"
	"wartank/app/lev2/arena/arena_arsenal/bf_ammo_stat"
	"wartank/app/lev2/arena/arena_arsenal/bf_arsenal_build"
	"wartank/app/lev2/arena/arena_arsenal/bf_arsenal_take"
	"wartank/app/lev2/arena/arena_arsenal/bf_arsenal_upgrade"
	"wartank/app/lev2/arena/arena_build"
)

const (
	стрКумулятивы = "кумулятивы"
	стрБронебойки = "бронебойки"
	стрФугасы     = "фугасы"
	стрРемки      = "ремки"
)

// Арсенал -- объект оружейной на базе
type АренаАрсенал struct {
	ИАренаСтроение
	вЛог       ИВебЛог
	лог        ILogBuf
	база       ИАренаБаза
	фугас      ИСтатПарам
	бронебойка ИСтатПарам
	кумулятив  ИСтатПарам
	ремка      ИСтатПарам
	конт       ILocalCtx
}

// НовАрсенал -- возвращает новый *Arsenal
func НовАрсенал(конт ILocalCtx) *АренаАрсенал {
	лог := NewLogBuf()
	лог.Info("НовАрсенал()\n")

	сам := &АренаАрсенал{
		конт:       конт,
		база:       конт.Get("база").Val().(ИАренаБаза),
		фугас:      lev1.НовСтатПарам(стрФугасы),
		бронебойка: lev1.НовСтатПарам(стрБронебойки),
		кумулятив:  lev1.НовСтатПарам(стрКумулятивы),
		ремка:      lev1.НовСтатПарам(стрРемки),
		лог:        лог,
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Арсенал",
		СтрКонтроль_: `<span class="green2">Ремкомплект</span><br/>`,
		СтрУрл_:      "https://wartank.ru/production/Armory",
	}
	сам.ИАренаСтроение = arena_build.НовАренаСтроение(конт, аренаКонфиг)
	сам.вЛог = web_log.НовВебЛог(true)
	// go сам.пуск()
	сам.вЛог.Добавить("НовАрсенал(): Арсенал создан")
	конт.Set("арсенал", сам, "Арсенал бота")
	_ = ИАренаАрсенал(сам)
	return сам
}

func (сам *АренаАрсенал) Пуск() {
	сам.ИАренаСтроение.Пуск()
	bf_arsenal_build.АрсеналПостроить(сам.конт)
	bf_arsenal_upgrade.АрсеналАпгрейд(сам.конт)
	bf_ammo_stat.СнарядыСтат(сам.конт)
	bf_ammo_make.СнарядыСделать(сам.конт)
	bf_arsenal_take.АрсеналЗабрать(сам.конт)
}

//=============================
/*
// Проверяет необходимость постройки
func (сам *АренаАрсенал) проверитьПостроить() bool {
	сам.вЛог.Добавить("Арсенал.проверитьПостроить()\n")
	_ = сам.Сеть().ВебВоркер().Получ("https://wartank.ru/building-upgrade/Armory")
	return true
}

// Обновляет текущий уровень арсенала (может быть не построена)
func (сам *АренаАрсенал) уровеньОбновить() bool {
	сам.вЛог.Добавить("Арсенал.уровеньОбновить()\n")
	списСтр := сам.Сеть().ВебВоркер().Получ("http://wartank.ru/buildings")
	// <span class="green2">Оружейная - 0</span><br/>
	var (
		еслиНайти = false
		стр       = ""
	)
	for _, стр = range списСтр {
		if strings.Contains(стр, `<span class="green2">Оружейная -`) {
			еслиНайти = true
			break
		}
	}
	if !еслиНайти {
		сам.вЛог.Добавить("Арсенал.уровеньОбновить(): не надо\n")
		return false
	}
	_стр := strings.TrimPrefix(стр, `<span class="green2">Оружейная - `)
	_стр = strings.TrimSuffix(_стр, `</span><br/>`)
	иУровень, ош := strconv.Atoi(_стр)
	if ош != nil {
		сам.лог.Err("уровеньОбновить(): строка уровня сбойная, стр=%q, ош=\n\t%v\n", стр, ош)
		сам.вЛог.Добавить("ОШИБКА Арсенал.уровеньОбновить(): строка уровня сбойная, стр=%q, ош=\n\t%v\n", стр, ош)
		return false
	}
	сам.Уровень().Уст(иУровень)
	сам.лог.Info("уровеньОбновить(): уровень=%d\n", иУровень)
	сам.вЛог.Добавить("Арсенал.уровеньОбновить(): уровень=%d\n", иУровень)
	return true
}

// Строит арсенал при нулевом уровне
func (сам *АренаАрсенал) построить() (bool, error) {
	сам.вЛог.Добавить("Арсенал.построить()\n")
	списСтр := сам.Сеть().ВебВоркер().Получ("https://wartank.ru/building-upgrade/Armory")
	// <span class="green2">Арсенал - 0</span><br/>
	var (
		еслиНайти = false
		стр       = ""
	)
	for _, стр = range списСтр {
		if strings.Contains(стр, `ILinkListener-upgradeLink-link`) {
			еслиНайти = true
			break
		}
	}
	if !еслиНайти {
		сам.вЛог.Добавить("Арсенал.построить(): не надо\n")
		return true, nil
	}
	// <a class="simple-but border mb5" href="Armory?30-1.ILinkListener-upgradeLink-link">
	// Пробуем построить арсенал
	_стр := strings.TrimPrefix(стр, `<a class="simple-but border mb5" href="`)
	_стр = strings.TrimSuffix(_стр, `">`)
	ссылка := "https://wartank.ru/building-upgrade/" + _стр
	// https://wartank.ru/building-upgrade/Armory?35-1.ILinkListener-upgradeLink-link
	списСтр = сам.Сеть().ВебВоркер().Получ(ссылка)
	еслиНайти = false
	// "<a class=\"simple-but border mb5\" href=\"Armory?14-1.ILinkListener-upgradeLink-link\">"
	for _, стр = range списСтр {
		if strings.Contains(стр, `ILinkListener-upgradeLink-link`) {
			еслиНайти = true
			break
		}
	}
	if !еслиНайти {
		сам.вЛог.Добавить("Арсенал.построить(): не надо\n")
		return true, nil
	}
	сам.лог.Info("построить(): ок\n")
	сам.вЛог.Добавить("Арсенал.построить(): ок\n")
	return true, nil
}

// Пытается проапгрейдить арсенал
func (сам *АренаАрсенал) проапгрейдить() bool {
	сам.вЛог.Добавить("Арсенал.проапгрейдить()\n")
	var (
		еслиНайти = false
		списСтр   []string
		стр       = ""
	)
	фнКупить := func() bool {
		defer time.Sleep(time.Millisecond * 1000)
		списСтр = сам.Сеть().ВебВоркер().Получ("https://wartank.ru/building-upgrade/Armory")
		for _, стр = range списСтр {
			// <a class="simple-but border mb5" href="Armory?5-1.ILinkListener-upgradeLink-link">
			if strings.Contains(стр, `ILinkListener-upgradeLink-link`) {
				еслиНайти = true
				break
			}
		}
		if !еслиНайти {
			сам.вЛог.Добавить("Арсенал.проапгрейдить(): не надо\n")
			return true
		}
		// Пробуем улучшить шахту
		_стр := strings.TrimPrefix(стр, "<a class=\"simple-but border mb5\" href=\"")
		_стр = strings.TrimSuffix(_стр, "\">")
		// https://wartank.ru/building-upgrade/Armory?4-1.ILinkListener-upgradeLink-link
		// <a class="simple-but border mb5" href="Armory?50-1.ILinkListener-upgradeLink-link">
		ссылка := "https://wartank.ru/building-upgrade/" + _стр
		списСтр = сам.Сеть().ВебВоркер().Получ(ссылка)
		// Проверить, что постройка состоялась
		for _, стр := range списСтр {
			if strings.Contains(стр, "ILinkListener-upgradeLink-link") {
				log.Printf("Арсенал.проапгрейдить().фнКупить(): покупка арсенала не прошла\n\tlink=%v\n\tстр=\n\t%v\n", ссылка, стр)
				return false // Покупка не оплачена
			}
		}
		сам.вЛог.Добавить("Арсенал.проапгрейдить().фнКупить(): ок\n")
		return true
	}

	фнПодтверждение := func() bool {
		for _, стр = range списСтр {
			// <a class="simple-but border w50 mXa mb10" w:id="confirmLink" href="../wicket/page?7-1.ILinkListener-confirmLink"><span><span>да, подтверждаю</span></span></a>
			if strings.Contains(стр, `ILinkListener-confirmLink`) {
				еслиНайти = true
				break
			}
		}
		if !еслиНайти {
			сам.вЛог.Добавить("Арсенал.проапгрейдить().фнПодтверждение(): не надо\n")
			return true
		}
		// Пробуем построить шахту
		_стр := strings.TrimPrefix(стр, `<a class="simple-but border w50 mXa mb10" w:id="confirmLink" href="..`)
		_стр = strings.TrimSuffix(_стр, `"><span><span>да, подтверждаю</span></span></a>`)
		// https://wartank.ru/wicket/page?6-1.ILinkListener-confirmLink
		ссылка := "https://wartank.ru" + _стр
		списСтр = сам.Сеть().ВебВоркер().Получ(ссылка)
		// Проверить, что постройка состоялась
		for _, стр := range списСтр {
			if strings.Contains(стр, "<title>Вы сделали слишком большую паузу</title>") {
				сам.вЛог.Добавить("ОШИБКА Арсенал.проапгрейдить().фнПодтверждение(): подтверждение покупка склада топлива не прошла\n\tlink=%v\n\tстр=\n\t%v\n", ссылка, стр)
				return false // Покупка не оплачена
			}
		}
		сам.вЛог.Добавить("Арсенал.проапгрейдить().фнПодтверждение(): ок\n")
		return true
	}

	фнКомплекс := func() {
		count := 5
		for count > 0 {
			if фнКупить() {
				if фнПодтверждение() {
					break
				}
			}
			count--
		}
	}
	фнКомплекс()
	return true
}

// Проверяет на забрать оружейную
func (сам *АренаАрсенал) забрать() bool {
	var (
		strOut      string
		ind         int
		еслиНайдено bool
		lstBase     = сам.СписПолучить()
	)
	for ind, strOut = range lstBase {
		if strings.Contains(strOut, `Моя амуниция`) {
			еслиНайдено = true
			ind += 17
			strOut = lstBase[ind]
			break
		}
	}
	if !еслиНайдено {
		return false
	}
	if !strings.Contains(strOut, `"><span><span>Забрать</span></span></a>`) {
		return false
	}
	lstLink := strings.Split(strOut, `<a class="simple-but border" href="`)
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Забрать</span></span></a>`)
	// https://wartank.ru/buildings?80-1.ILinkListener-buildings-0-building-rootBlock-actionPanel-takeProductionLink
	strLink = "https://wartank.ru/" + lstLink[0]
	var (
		лстАрсенал []string
	)
	time.Sleep(time.Millisecond * 100)
	res := сам.Сеть().Get(strLink)
	if res.IsErr() {
		log.Printf("Арсенал.забрать(): при выполнении Get-запроса? err=\n\t%v\n", res.Error())
		return false
	}
	лстАрсенал = res.Unwrap()
	if len(лстАрсенал) == 0 {
		log.Printf("Арсенал.забрать(): len lstBase(%v)==0", len(lstBase))
		return false
	}
	for _, strOut = range лстАрсенал {
		if strings.Contains(strOut, `<title>Производство</title>`) {
			return false
		}
	}
	сам.СтрОбновить(лстАрсенал)
	return true
}
*/
//====================================
// Фугасы -- возвращает объект числа фугасов
func (сам *АренаАрсенал) Фугасы() ИСтатПарам {
	return сам.фугас
}

// Бронебойки -- возвращает объект бронебойных снарядов
func (сам *АренаАрсенал) Бронебойки() ИСтатПарам {
	return сам.бронебойка
}

// Кумулятивы -- возвращает объект бронебойных снарядов
func (сам *АренаАрсенал) Кумулятивы() ИСтатПарам {
	return сам.кумулятив
}

// Ремки -- возвращает объект ремкомплектов
func (сам *АренаАрсенал) Ремки() ИСтатПарам {
	return сам.ремка
}
