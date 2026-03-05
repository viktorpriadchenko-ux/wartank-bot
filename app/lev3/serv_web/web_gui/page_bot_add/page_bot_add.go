// package page_bot_add -- страница добавления бота
package page_bot_add

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev1/web_render"
)

// СтраницаДобавлениеБота -- страница добавления бота
type СтраницаБотаДобавить struct {
	лог    ILogBuf
	прилож ИПриложение
	рендер ИВебРендер
}

//go:embed page_bot_add.tmpl.html
var стрБотДобавить string

// НовСтраницаБотовДобавить
func НовСтраницаБотовДобавить(конт IKernelCtx) *СтраницаБотаДобавить {
	лог := NewLogBuf()
	Hassert(конт != nil, "НовСтраницаБотовДобавить(): ИЯдроКонтекст==nil")
	сам := &СтраницаБотаДобавить{
		лог:    лог,
		прилож: конт.Get("мод_сервер").Val().(ИПриложение),
		рендер: web_render.НовВебРендер(стрБотДобавить),
	}
	файбер := конт.Get("fiberApp").Val().(*fiber.App)
	файбер.Post("/gui/bot/add/show", сам.гетБотНов)
	файбер.Post("/gui/bot/add/make", сам.постДобавитьБота)
	return сам
}

type AddBotRequest struct {
	Логин_ string `form:"login_bot"`
	Пароль string `form:"password_bot"`
}

// Запрос добавления нового бота на ботоферму
func (сам *СтраницаБотаДобавить) постДобавитьБота(кнт *fiber.Ctx) error {
	запрос := &AddBotRequest{}
	if err := кнт.BodyParser(запрос); err != nil {
		return сам.показатьФормуСОшибкой(кнт, "Ошибка чтения формы: "+err.Error())
	}
	сам.лог.Debug("добавитьБота(): : %#+v\n", запрос)
	if запрос.Логин_ == "" {
		return сам.показатьФормуСОшибкой(кнт, "Логин не может быть пустым")
	}
	if запрос.Пароль == "" {
		return сам.показатьФормуСОшибкой(кнт, "Пароль не может быть пустым")
	}
	if res := сам.прилож.ServBots().НовБот(запрос.Логин_, запрос.Пароль, true); res.IsErr() {
		return сам.показатьФормуСОшибкой(кнт, "Не удалось добавить бота: "+res.Error().Error())
	}
	// Успех — htmx перезагружает всю страницу
	кнт.Set("HX-Redirect", "/")
	return кнт.SendStatus(200)
}

// показатьФормуСОшибкой -- возвращает форму с сообщением об ошибке
func (сам *СтраницаБотаДобавить) показатьФормуСОшибкой(кнт *fiber.Ctx, сообщение string) error {
	кнт.Set("Content-type", "text/html; charset=utf-8")
	эррБлок := fmt.Sprintf(`<div class="row p-3"><div class="col"><div class="alert alert-danger">%s</div></div></div>`, сообщение)
	стрРез := strings.Replace(сам.рендер.Получ(), "</form>", эррБлок+"</form>", 1)
	return кнт.SendString(стрРез)
}

// Показывает страницу добавления бота
func (сам *СтраницаБотаДобавить) гетБотНов(кнт *fiber.Ctx) error {
	сам.лог.Debug("СтраницаБотаДобавить.гетБотНов()\n")
	кнт.Set("Content-type", "text/html; charset=utf-8")
	стрРез := сам.рендер.Получ()
	return кнт.SendString(стрРез)
}

/*
// Проверка на куки
func (сам *СтраницаБотаДобавить) _кукиПроверить(кнт *fiber.Ctx) error {
	сам.лог.Debug("СтраницаБотаДобавить.кукиПроверить()\n")
	имя := кнт.Cookies("login")
	if имя != "svi" {
		return кнт.Redirect("/gui/login-show", http.StatusSeeOther)
	}
	return кнт.Next()
}
*/
