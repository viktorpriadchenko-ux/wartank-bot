// package page_bot_list -- страница списка ботов
package page_bot_list

import (
	_ "embed"
	"fmt"
	// "net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev1/web_render"
)

// СтраницаЛогин -- страница списка ботов
type СтраницаСписокБотов struct {
	лог    ILogBuf
	прилож ИПриложение
	рендер ИВебРендер
}

//go:embed block_list_bot.html
var стрСписок string

// НовСтраницаСписокБотов
func НовСтраницаСписокБотов(конт IKernelCtx) *СтраницаСписокБотов {
	лог := NewLogBuf()
	Hassert(конт != nil, "НовСтраницаСписокБотов(): ИЯдроКонтекст==nil")
	сам := &СтраницаСписокБотов{
		лог:    лог,
		прилож: конт.Get("мод_сервер").Val().(ИПриложение),
		рендер: web_render.НовВебРендер(стрСписок),
	}
	файбер := конт.Get("fiberApp").Val().(*fiber.App)
	файбер.Post("/gui/bot/list_show", сам.получСписокБотов)
	return сам
}

var стрСсылкаШаблон = `
<div>
    <a class="btn btn-primary" hx-post="/gui/bot/{.number}/show" hx-target="#main">{.name}</a><br><br>
</div>
`

// Возвращает страницу логина
func (сам *СтраницаСписокБотов) получСписокБотов(кнт *fiber.Ctx) error {
	сам.лог.Debug("СтраницаСписокБотов.получСписокБотов()\n")
	кнт.Response().Header.Add("Content-type", "text/html; charset=utf-8")
	кнт.Response().Header.Add("Cache-Control", "no-cache")
	списокБотов := сам.прилож.ServBots().ListBot()
	стрСсылки := ""
	for _, бот := range списокБотов {
		стрНомер := fmt.Sprint(бот.Номер())
		стрНомер = strings.ReplaceAll(стрСсылкаШаблон, "{.number}", стрНомер)
		стрНомер = strings.ReplaceAll(стрНомер, "{.name}", бот.Имя())
		стрСсылки += стрНомер + "<br>\n"
	}

	сам.рендер.Доб("{.list_bot}", стрСсылки)
	стрРез := сам.рендер.Получ()
	return кнт.SendString(стрРез)
}

/*
// Проверка на куки
func (сам *СтраницаСписокБотов) кукиПроверить(кнт *fiber.Ctx) error {
	сам.лог.Debug("СтраницаСписокБотов.кукиПроверить()\n")
	имя := кнт.Cookies("login")
	if имя != "svi" {
		return кнт.Redirect("/gui/login-show", http.StatusSeeOther)
	}
	return кнт.Next()
}
*/
