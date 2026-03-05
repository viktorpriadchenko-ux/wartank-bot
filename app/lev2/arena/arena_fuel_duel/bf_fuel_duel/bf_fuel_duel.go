// package bf_fuel_duel -- бизнес-функция боя на топливе
package bf_fuel_duel

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// ТопливоАтаковать -- бой на топливе
func ТопливоАтаковать(конт ILocalCtx) {
	ангар := конт.Get("ангар").Val().(ИАренаАнгар)
	топливо := ангар.Топливо().Получ()
	if топливо < cons.ТопливоМин { // Минимальная ёмкость бака -- 315
		return
	}
	списСтрБой := начатьБой(конт)
	списВыстрел1 := выбратьБойСлабый(конт, списСтрБой)
	сделатьВыстрелы(конт, списВыстрел1)

	// Статистика: записать дуэль за топливо
	бот := конт.Get("бот").Val().(ИБот)
	бот.Статистика().ДуэльТопливоДобавить(true) // TODO: определять исход
	потрачено := топливо - ангар.Топливо().Получ()
	if потрачено > 0 {
		бот.Статистика().ТопливоДобавитьРасход(потрачено)
	}
}

// Идёт в атаку, если топлива больше cons.ТопливоМин
func начатьБой(конт ILocalCtx) []string {
	// Получить ссылку на атаку
	ангар := конт.Get("ангар").Val().(ИАренаАнгар)
	списАнгар := ангар.СписПолучить()
	var (
		стрАнализ string
		еслиЕсть  bool
	)
	for _, стрАтак := range списАнгар {
		if strings.Contains(стрАтак, `<span>В бой!</span>`) {
			стрАнализ = стрАтак
			еслиЕсть = true
			break
		}
	}
	Hassert(еслиЕсть, "начатьБой(): не найдена строка 'Начать бой'")
	// Вырезать ссылку на атаку
	списАнгар = strings.Split(стрАнализ, `<a class="simple-but border mb1" href="`)
	Hassert(len(списАнгар) >= 2, "начатьБой(): список строк для атаки пустой")
	бойСсылка := списАнгар[1]
	списАнгар = strings.Split(бойСсылка, `"><span><span>В бой!</span></span></a>`)
	бойСсылка = "https://wartank.ru/" + списАнгар[0]
	арена := конт.Get("арена_топливо_бой").Val().(ИАренаСтроение)
	res := арена.Сеть().Get(бойСсылка)
	res.Hassert("начатьБой(): in make GET-request to battle")
	return res.Unwrap()
}

// Выбирает первого более слабого противника и делает первый выстрел
func выбратьБойСлабый(конт ILocalCtx, списСтрБой []string) []string {
	// _mt.Println("\tAngarNet.makeSelectBattle()")
	var стрАнализ string
	// Выдернуть строку с первой ссылкой на противника
	for _, стрБой_ := range списСтрБой {
		if strings.Contains(стрБой_, `opponents-opponents-0`) {
			стрАнализ = стрБой_
			break
		}
	}
	var ссылкаБой string
	switch стрАнализ == "" {
	case true: // Такая ситуация возможна, если уже были какие-то выстрелы
		return списСтрБой
	default: // Успешный выстрел
		// Вырезать ссылку из строки
		списСтрБой = strings.Split(стрАнализ, `<td class="cntr"><a href="`)
		ссылкаБой = списСтрБой[1]
		списСтрБой = strings.Split(ссылкаБой, `"><img class="tank-img" alt="tank" src="/tankimg?`)
		ссылкаБой = "https://wartank.ru/" + списСтрБой[0]
	}
	аренаБой := конт.Get("арена_топливо_бой").Val().(ИАрена)
	res := аренаБой.Сеть().Get(ссылкаБой)
	res.Hassert("makeSelectBattle(): in GET-response select battle tank")
	return res.Unwrap() // Первый выстрел
}

// Ведёт бой в 2 выстрела (здесь только 2 и 3 выстрел -- первый сделан при слабом противнике)
func сделатьВыстрелы(конт ILocalCtx, lstShoot2 []string) {
	// _mt.Println("\tAngarNet.makeShooting()")
	var списВыстрел3 []string // Тело страницы для третьего выстрела
	аренаБой := конт.Get("арена_топливо_бой").Val().(ИАрена)
	ангар := конт.Get("ангар").Val().(ИАренаАнгар)
	фнВыстрел2 := func() { // Второй выстрел
		// _mt.Println("\tAngarNet.makeShooting().fnShoot2()")
		defer func() {

			_ = recover()
		}()
		// Получить ссылку на второй выстрел
		var strOut string
		for _, strShoot := range lstShoot2 {
			if strings.Contains(strShoot, `<span>Добить</span>`) {
				strOut = strShoot
				break
			}
		}
		var linkShoot2 string
		switch strOut == "" {
		case true: // Первый выстрел был неудачным
			for _, strShoot := range lstShoot2 {
				if strings.Contains(strShoot, `<span>Взять реванш</span>`) {
					strOut = strShoot
					break
				}
			}
			if strOut == "" { // Это ситуация для третьего выстрела
				списВыстрел3 = lstShoot2
				return
			}
			// Вырезать ссылку из строки
			lstShoot2 = strings.Split(strOut, `<a class="simple-but border" href="`)
			linkShoot2 = lstShoot2[1]
			lstShoot2 = strings.Split(linkShoot2, `"><span><span>Взять реванш</span></span></a>`)
			linkShoot2 = "https://wartank.ru/" + lstShoot2[0]
		default: // Первый выстрел был удачным
			// Вырезать ссылку из строки
			lstShoot2 = strings.Split(strOut, `<a class="simple-but border" href="`)
			linkShoot2 = lstShoot2[1]
			lstShoot2 = strings.Split(linkShoot2, `"><span><span>Добить</span></span></a>`)
			linkShoot2 = "https://wartank.ru/" + lstShoot2[0]
		}
		res := аренаБой.Сеть().Get(linkShoot2)
		res.Hassert("сделатьВыстрелы(): in Get-response shoot2")

		fuel := ангар.Топливо().Получ()
		fuel -= 30
		ангар.Топливо().Уст(fuel)
	}
	фнВыстрел2()
	фнВыстрел3 := func() { // Третий выстрел
		// _mt.Println("\tAngarNet.makeShooting().fnShoot3()")
		defer func() {
			_ = recover()
		}()
		// Получить ссылку на третий выстрел
		var strOut string
		for _, strShoot3 := range списВыстрел3 {
			if strings.Contains(strShoot3, `<span>Уничтожить</span>`) {
				strOut = strShoot3
				break
			}
		}
		linkShoot3 := ""
		switch strOut == "" {
		case true: // Если не найдена ссылка -- значит было поражение в выстреле
			if strOut == "" {
				for _, strShoot3 := range списВыстрел3 {
					if strings.Contains(strShoot3, `<span>Взять реванш</span>`) {
						strOut = strShoot3
						break
					}
				}
			}
			// Вырезать ссылку из строки
			списВыстрел3 = strings.Split(strOut, `<a class="simple-but border" href="`)
			linkShoot3 = списВыстрел3[1]
			списВыстрел3 = strings.Split(linkShoot3, `"><span><span>Взять реванш</span></span></a>`)
			linkShoot3 = "https://wartank.ru/" + списВыстрел3[0]
		default: // Успешный выстрел
			// Вырезать ссылку из строки
			списВыстрел3 = strings.Split(strOut, `<a class="simple-but border" href="`)
			linkShoot3 = списВыстрел3[1]
			списВыстрел3 = strings.Split(linkShoot3, `"><span><span>Уничтожить</span></span></a>`)
			linkShoot3 = "https://wartank.ru/" + списВыстрел3[0]
		}

		res := аренаБой.Сеть().Get(linkShoot3)
		res.Hassert("ТопливоБой.makeShooting(): in Get-response shoot3")
		fuel := ангар.Топливо().Получ()
		fuel -= 30
		ангар.Топливо().Уст(fuel)
	}
	фнВыстрел3()
	ангар.Обновить()
}
