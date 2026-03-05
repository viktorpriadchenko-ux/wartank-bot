// package bf_glory_take -- забирает призы
package bf_glory_take

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
)

// СлаваНайти -- берёт славу бота
func СлаваВзять(конт ILocalCtx) {
	проверитьМиссияРазведкаКонвой(конт)
	проверитьМиссияМастерРазведки(конт)
	проверитьМиссия6фрагов(конт)
}

// Забирает награду в конвое "Активируй боевую силу"
func проверитьМиссияРазведкаКонвой(конт ILocalCtx) {
	конвой := конт.Get("конвой").Val().(ИАренаКонвой)
	var (
		strOut      string
		еслиНайдено bool
	)
	lstConvoy := конвой.СписПолучить()
	// <a class="simple-but border" href="convoy?21-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	for _, strOut = range lstConvoy {
		if strings.Contains(strOut, `.ILinkListener-missions-cc-0-c-awardLink`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="convoy?21-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	_ссылка := strings.TrimPrefix(strOut, `<a class="simple-but border" href="`)
	_ссылка = strings.TrimSuffix(_ссылка, `"><span><span>Получить награду</span></span></a>`)
	// https://wartank.ru/convoy?23-1.ILinkListener-missions-cc-0-c-awardLink
	ссылка := "https://wartank.ru/" + _ссылка
	lstConvoy = конвой.Сеть().ВебВоркер().Получ(ссылка)
	конвой.СтрОбновить(lstConvoy)
}

// Забирает награду в конвое "Мастер дозора"
func проверитьМиссияМастерРазведки(конт ILocalCtx) {
	конвой := конт.Get("конвой").Val().(ИАренаКонвой)
	var (
		strOut      string
		еслиНайдено bool
		lstConvoy   = конвой.СписПолучить()
		ind         int
	)
	if len(lstConvoy) == 0 {
		конвой.Обновить()
		lstConvoy = конвой.СписПолучить()
	}
	for ind, strOut = range lstConvoy {
		if strings.Contains(strOut, `Проведи разведку в конвое<br/>`) {
			еслиНайдено = true
			ind += 23
			strOut = lstConvoy[ind]
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="convoy?61-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	if !strings.Contains(strOut, `ILinkListener-missions-cc-0-c-awardLink`) {
		return
	}
	lstLink := strings.Split(strOut, `<a class="simple-but border" href="`)
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Получить награду</span></span></a>`)
	// https://wartank.ru/convoy?61-1.ILinkListener-missions-cc-0-c-awardLink
	strLink = "https://wartank.ru/" + lstLink[0]
	res := конвой.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Конвой.checkMaster(): при выполнении команды GET, err=\n\t%v\n", err)
		return
	}
	конвой.СтрОбновить(res.Unwrap())
}

// Забирает награду в конвое "Уничтожь 6 врагов в конвое"
func проверитьМиссия6фрагов(конт ILocalCtx) {
	конвой := конт.Get("конвой").Val().(ИАренаКонвой)
	var (
		strOut      string
		еслиНайдено bool
	)
	lstConvoy := конвой.СписПолучить()
	// <a class="simple-but border" href="convoy?8-1.ILinkListener-missions-cc-1-c-awardLink"><span><span>Получить награду</span></span></a>
	for _, strOut = range lstConvoy {
		if strings.Contains(strOut, `.ILinkListener-missions-cc-1-c-awardLink`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="convoy?8-1.ILinkListener-missions-cc-1-c-awardLink"><span><span>Получить награду</span></span></a>
	_ссылка := strings.TrimPrefix(strOut, `<a class="simple-but border" href="`)
	_ссылка = strings.TrimSuffix(_ссылка, `"><span><span>Получить награду</span></span></a>`)
	// https://wartank.ru/convoy?15-1.ILinkListener-missions-cc-1-c-awardLink
	ссылка := "https://wartank.ru/" + _ссылка
	lstConvoy = конвой.Сеть().ВебВоркер().Получ(ссылка)
	конвой.СтрОбновить(lstConvoy)
}
