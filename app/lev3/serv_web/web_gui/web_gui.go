// package web_gui -- веб-сервер для графики
package web_gui

import (
	_ "embed"
	// "net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev3/serv_web/web_gui/btn_list_bot"
	"wartank/app/lev3/serv_web/web_gui/btn_login"
	"wartank/app/lev3/serv_web/web_gui/page_bot_add"
	"wartank/app/lev3/serv_web/web_gui/page_bot_list"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show"
	"wartank/app/lev3/serv_web/web_gui/page_login"
)

// ВебГип -- веб-сервер для графики
type ВебГуи struct {
	прилож   ИПриложение
	конт     IKernelCtx
	лог      ILogBuf
	кнпЛогин *btn_login.КнпЛогин
	кнпБоты  *btn_list_bot.КнпСписБот
}

// НовВебГуи -- возвращает новый веб-сервер для графики
func НовВебГуи() *ВебГуи {
	лог := NewLogBuf()
	конт := GetKernelCtx()
	сам := &ВебГуи{
		прилож:   конт.Get("мод_сервер").Val().(ИПриложение),
		конт:     конт,
		лог:      лог,
		кнпЛогин: btn_login.НовКнпЛогин(),
	}
	сам.кнпБоты = btn_list_bot.НовКнпСписБот(сам.кнпЛогин)
	файбер := конт.Get("fiberApp").Val().(*fiber.App)
	// файбер.Get("/", сам.кукиПроверить) //, сам.индекс)
	файбер.Get("/", сам.индекс)
	_ = page_login.НовСтраницаЛогин(конт)
	_ = page_bot_list.НовСтраницаСписокБотов(конт)
	_ = page_bot_add.НовСтраницаБотовДобавить(конт)
	_ = page_bot_show.НовСтраницаБотПоказать(конт)
	return сам
}

//go:embed index.tmpl.html
var стрГлавная string

// Возвращает индексную страницу
func (сам *ВебГуи) индекс(кнт *fiber.Ctx) error {
	сам.лог.Debug("ВебГуи.индекс()\n")
	кнт.Set("Content-type", "text/html; charset=utf-8")
	кнт.Set("Cache-Control", "no-cache")
	_, _ = кнт.WriteString("\n\n")
	стрВых := strings.ReplaceAll(стрГлавная, "{.btn_lst_bot}", сам.кнпБоты.Html())
	стрВых = strings.ReplaceAll(стрВых, "{.url_login}", сам.кнпЛогин.Кнп.Hx().Url().String())
	return кнт.SendString(стрВых)
}

/*
// Проверка на куки
func (сам *ВебГуи) кукиПроверить(кнт *fiber.Ctx) error {
	сам.лог.Debug("ВебГуи.кукиПроверить()\n")
	имя := кнт.Cookies("login")
	if имя != "svi" {
		сам.конт.Set("еслиЛогин", false, "Признак логина")
	} else {
		сам.конт.Set("еслиЛогин", true, "Признак логина")
	}
	return кнт.Next()
}
*/
