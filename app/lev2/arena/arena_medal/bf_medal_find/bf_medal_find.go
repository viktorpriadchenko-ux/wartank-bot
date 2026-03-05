// package bf_medal_find -- поиск медалей

package bf_medal_find

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
)

// СлаваНайти -- ищет медали бота
func МедалиНайти(конт ILocalCtx) {
	медали := конт.Get("арена_медали").Val().(ИАрена)
	// Найти строку с упоминанием оставшегося времени конвоя
	медали.Обновить()
	lstStr := медали.СписПолучить()
	var (
		стрМедаль string
	)
	фнНайти := func() bool {
		// <a class="simple-but border" href="current?128-1.ILinkListener-currentMedal-takeAwardLink"><span><span>Получить медаль</span></span></a>
		for _, val := range lstStr {
			if strings.Contains(val, `-currentMedal-takeAwardLink"><span><span>Получить медаль</span></span></a>`) {
				стрМедаль = val
				return true
			}
		}
		return false
	}

	// Статистика: получить бота для счётчиков
	бот := конт.Get("бот").Val().(ИБот)

	if фнНайти() {
		// Вырезаем ссылку  на медаль
		// <a class="simple-but border" href="current?128-1.ILinkListener-currentMedal-takeAwardLink"><span><span>Получить медаль</span></span></a>
		стрМедаль = strings.TrimPrefix(стрМедаль, `<a class="simple-but border" href="`)
		стрМедаль = strings.TrimSuffix(стрМедаль, `"><span><span>Получить медаль</span></span></a>`)
		// https://wartank.ru/medals/current?137-1.ILinkListener-currentMedal-takeAwardLink
		// https://wartank.ru/medals/current?169-1.ILinkListener-currentMedal-takeAwardLink
		ссыль := "https://wartank.ru/medals/" + стрМедаль
		lstStr = медали.Сеть().ВебВоркер().Получ(ссыль)
		бот.Статистика().МедальДобавить()
		for фнНайти() {
			// Вырезаем ссылку  на медаль
			// <a class="simple-but border" href="current?128-1.ILinkListener-currentMedal-takeAwardLink"><span><span>Получить медаль</span></span></a>
			стрМедаль = strings.TrimPrefix(стрМедаль, `<a class="simple-but border" href="`)
			стрМедаль = strings.TrimSuffix(стрМедаль, `"><span><span>Получить медаль</span></span></a>`)
			// https://wartank.ru/medals/current?137-1.ILinkListener-currentMedal-takeAwardLink
			// https://wartank.ru/medals/current?169-1.ILinkListener-currentMedal-takeAwardLink
			ссыль := "https://wartank.ru/medals/" + стрМедаль
			lstStr = медали.Сеть().ВебВоркер().Получ(ссыль)
			бот.Статистика().МедальДобавить()
		}
	}
}
