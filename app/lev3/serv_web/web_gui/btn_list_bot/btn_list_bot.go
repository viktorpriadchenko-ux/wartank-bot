// package btn_list_bot -- кнопка списка ботов
package btn_list_bot

import (
	_ "embed"
	"fmt"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/wui"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev1/web_render"
	"wartank/app/lev3/serv_web/web_gui/btn_login"
)

// КнпЛогин -- кнопка списка ботов на сайте
type КнпСписБот struct {
	конт     IKernelCtx
	прилож   ИПриложение
	кнп      IWuiButton
	кнпЛогин *btn_login.КнпЛогин
	рендер   ИВебРендер
}

//go:embed block_list_bot.html
var стрСписок string

// НовКнпСписБот -- возвращает новую кнопку списка ботов
func НовКнпСписБот(кнпЛогин *btn_login.КнпЛогин) *КнпСписБот {
	Hassert(кнпЛогин != nil, "НовКнпСписБот(): кнпЛогин==nil")
	конт := GetKernelCtx()
	сам := &КнпСписБот{
		конт:     конт,
		прилож:   конт.Get("мод_сервер").Val().(ИПриложение),
		кнпЛогин: кнпЛогин,
		рендер:   web_render.НовВебРендер(стрСписок),
	}
	сам.кнп = wui.NewWuiButton("_Список ботов_", сам.КликСписБот)
	сам.кнп.Hx().Target().Set("#main")
	конт.Set("кнп_спис_бот", сам, "Кнопка показа списка ботов")
	return сам
}

// Html -- возвращает HTML-представление кнопки
func (сам *КнпСписБот) Html() string {
	return сам.кнп.Html()
}

var стрСсылкаШаблон = `
<div>
    <a class="btn btn-primary" hx-post="/gui/bot/{.id}/show" hx-target="#main">{.name}</a><br><br>
</div>
`

// Событие клика по кнопке
func (сам *КнпСписБот) КликСписБот(слв map[string]string) string {
	списокБотов := сам.прилож.ServBots().ListBot()
	стрСписБот := ""
	for _, бот := range списокБотов {
		стрНомер := fmt.Sprint(бот.Номер())
		стрНомер = strings.ReplaceAll(стрСсылкаШаблон, "{.id}", стрНомер)
		стрНомер = strings.ReplaceAll(стрНомер, "{.name}", бот.Имя())
		стрСписБот += стрНомер + "<br>\n"
	}

	сам.рендер.Доб("{.list_bot}", стрСписБот)
	стрРез := сам.рендер.Получ()
	return стрРез
}
