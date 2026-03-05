// package page_login -- страница логина
package page_login

import (
	_ "embed"
	// "net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev1/web_render"
)

// СтраницаЛогин -- страница логина
type СтраницаЛогин struct {
	лог    ILogBuf
	рендер ИВебРендер
}

//go:embed login.tmpl.html
var стрЛогин string

// НовСтраницаЛогин
func НовСтраницаЛогин(конт IKernelCtx) *СтраницаЛогин {
	лог := NewLogBuf()
	Hassert(конт != nil, "НовСтраницаЛогин(): ИЯдроКонтекст==nil")
	сам := &СтраницаЛогин{
		лог:    лог,
		рендер: web_render.НовВебРендер(стрЛогин),
	}
	сам.рендер.Доб("{.err}", "")
	файбер := конт.Get("fiberApp").Val().(*fiber.App)
	файбер.Post("/gui/login_show", сам.получЛогин)
	файбер.Post("/gui/login_make", сам.постЛогин)
	return сам
}

var (
	strBotList = `
	<div id="main" hx-post="/gui/bot/list_show" hx-trigger="load"  hx-swap-oob""></div>
	`
)

// Возвращает страницу логина
func (сам *СтраницаЛогин) получЛогин(кнт *fiber.Ctx) error {
	сам.лог.Debug("СтраницаЛогин.логин()\n")
	кнт.Set("Content-type", "text/html; charset=utf-8")
	имя := кнт.Cookies("login")
	if имя == "svi" {
		return кнт.SendString(strBotList)
	}
	стрРез := сам.рендер.Получ()
	return кнт.SendString(стрРез)
}

type LoginRequest struct {
	Логин_            string `form:"login"`
	Пароль_           string `form:"password"`
	КонтрольноеСлово_ string `form:"control_word"`
}

//go:embed block_list_bot.html
var стрБлокСписБотов string

// Вызывается при попытке войти
func (сам *СтраницаЛогин) постЛогин(кнт *fiber.Ctx) error {
	запрос := new(LoginRequest)
	if err := кнт.BodyParser(запрос); err != nil {
		return кнт.SendString(стрЛогин)
	}
	сам.лог.Info("СтраницаЛогин.логин(): : %#+v\n", *запрос)
	if запрос.Логин_ == "" || запрос.Пароль_ == "" {
		return кнт.SendString(стрЛогин)
	}

	if запрос.Логин_ != "svi" || запрос.Пароль_ != "Lera_07091978" {
		return кнт.SendString(стрЛогин)
	}
	кнт.Cookie(&fiber.Cookie{
		Name:     "login",
		Value:    "svi",
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})
	return кнт.SendString(стрБлокСписБотов)
}

/*
// Проверка на куки
func (сам *СтраницаЛогин) кукиПроверить(кнт *fiber.Ctx) error {
	сам.лог.Debug("СтраницаЛогин.кукиПроверить()\n")
	имя := кнт.Cookies("login")
	if имя != "svi" {
		return кнт.Redirect("/gui/login-show", http.StatusSeeOther)
	}
	return кнт.Next()
}
*/
