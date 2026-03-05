// package page_bot_show -- страница показа бота
package page_bot_show

import (
	_ "embed"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/alias"
	. "wartank/app/lev0/types"
	"wartank/app/lev1/web_render"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_attack"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_bron"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_fuel"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_fyne"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_glory"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_gold"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_hard"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_mine_count"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_mine_level"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_mine_mode"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_mine_time"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_mine_type"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_power"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/btn_silver"
)

// СтраницаБотПоказать -- страница показа бота
type СтраницаБотПоказать struct {
	лог             ILogBuf
	прилож          ИПриложение
	рендер          ИВебРендер
	кнпЗолото       *btn_gold.BtnGold
	кнпСеребро      *btn_silver.BtnSilver
	кнпТопливо      *btn_fuel.BtnFuel
	кнпСлава        *btn_glory.BtnGlory
	кнпАтака        *btn_attack.BtnAttack
	кнпБроня        *btn_bron.BtnBron
	кнпТочность     *btn_fyne.BtnFyne
	кнпПрочность    *btn_hard.BtnHard
	кнпМощь         *btn_power.BtnPower
	кнпШахтаУровень *btn_mine_level.BtnMineLevel
	кнпШахтаРежим   *btn_mine_mode.BtnMineMode
	кнпШахтаПродукт *btn_mine_count.BtnMineCount
	кнпШахтаТип     *btn_mine_type.BtnMineType
	кнпШахтаВремя   *btn_mine_time.BtnMineTime
}

//go:embed bot_show.tmpl.html
var стрБотПоказать string

// НовСтраницаБотПоказать
func НовСтраницаБотПоказать(конт IKernelCtx) *СтраницаБотПоказать {
	лог := NewLogBuf()
	Hassert(конт != nil, "НовСтраницаБотПоказать(): ИЯдроКонтекст==nil")
	сам := &СтраницаБотПоказать{
		лог:             лог,
		прилож:          конт.Get("мод_сервер").Val().(ИПриложение),
		рендер:          web_render.НовВебРендер(стрБотПоказать),
		кнпЗолото:       btn_gold.NewBtnGold(),
		кнпСеребро:      btn_silver.NewBtnSilver(),
		кнпТопливо:      btn_fuel.NewBtnFuel(),
		кнпСлава:        btn_glory.NewBtnGlory(),
		кнпАтака:        btn_attack.NewBtnAttack(),
		кнпБроня:        btn_bron.NewBtnBron(),
		кнпТочность:     btn_fyne.NewBtnFyne(),
		кнпПрочность:    btn_hard.NewBtnHard(),
		кнпМощь:         btn_power.NewBtnPower(),
		кнпШахтаУровень: btn_mine_level.NewBtnMineLevel(),
		кнпШахтаРежим:   btn_mine_mode.NewBtnMineMode(),
		кнпШахтаПродукт: btn_mine_count.NewBtnMineCount(),
		кнпШахтаТип:     btn_mine_type.NewBtnMineType(),
		кнпШахтаВремя:   btn_mine_time.NewBtnMineTime(),
	}
	файбер := конт.Get("fiberApp").Val().(*fiber.App)
	файбер.Post("/gui/bot/:id/show", сам.гетБотПоказ)
	// файбер.Post("/gui/bot/add", сам.кукиПроверить, сам.постДобавитьБота)
	return сам
}

// Показывает состояние бота по имени
func (сам *СтраницаБотПоказать) гетБотПоказ(кнт *fiber.Ctx) error {
	сам.лог.Debug("СтраницаБотПоказать.гетБотПоказ()\n")
	кнт.Set("Content-type", "text/html; charset=utf-8")
	кнт.Set("Cache-Control", "no-cache")
	стрНомер := кнт.Params("id")
	иНомер, ош := strconv.Atoi(стрНомер)
	if ош != nil {
		return кнт.Redirect("/gui/bot", http.StatusSeeOther)
	}
	ботНомер := АБотНомер(иНомер)

	сам.лог.Debug("СтраницаБотПоказать.гетБотПоказ(): ботНомер=%d\n", стрНомер)
	бот := сам.прилож.ServBots().Get(ботНомер)
	if бот == nil {
		return кнт.Redirect("/gui/bot", http.StatusSeeOther)
	}
	{ // Глобальные показатели
		сам.кнпЗолото.Btn.Hx().Vals().Set("id", иНомер)
		сам.кнпСеребро.Btn.Hx().Vals().Set("id", иНомер)
		сам.кнпТопливо.Btn.Hx().Vals().Set("id", иНомер)
		сам.кнпСлава.Btn.Hx().Vals().Set("id", иНомер)

		сам.рендер.Доб("{.имя}", бот.Имя())
		номер := бот.Номер()
		сам.рендер.Доб("{.id}", номер)
		сам.рендер.Доб("{.block_gold}", сам.кнпЗолото.Html())
		сам.рендер.Доб("{.block_silver}", сам.кнпСеребро.Html())
		сам.рендер.Доб("{.block_fuel}", сам.кнпТопливо.Html())
		сам.рендер.Доб("{.block_glory}", сам.кнпСлава.Html())
	}
	{ // Сила танка
		сам.кнпАтака.Btn.Hx().Vals().Set("id", иНомер)
		сам.кнпБроня.Btn.Hx().Vals().Set("id", иНомер)
		сам.кнпТочность.Btn.Hx().Vals().Set("id", иНомер)
		сам.кнпПрочность.Btn.Hx().Vals().Set("id", иНомер)
		сам.кнпМощь.Btn.Hx().Vals().Set("id", иНомер)

		сам.рендер.Доб("{.attack}", сам.кнпАтака.Html())
		сам.рендер.Доб("{.броня}", сам.кнпБроня.Html())
		сам.рендер.Доб("{.точность}", сам.кнпТочность.Html())
		сам.рендер.Доб("{.прочность}", сам.кнпПрочность.Html())
		мощь := сам.кнпМощь.Html()
		сам.рендер.Доб("{.block_power}", мощь)
	}
	{ // Шахта
		сам.кнпШахтаУровень.Btn.Hx().Vals().Set("id", иНомер)
		сам.кнпШахтаРежим.Btn.Hx().Vals().Set("id", иНомер)
		сам.кнпШахтаПродукт.Btn.Hx().Vals().Set("id", иНомер)
		сам.кнпШахтаТип.Btn.Hx().Vals().Set("id", иНомер)
		сам.кнпШахтаВремя.Btn.Hx().Vals().Set("id", иНомер)

		сам.рендер.Доб("{.шахта_уровень}", сам.кнпШахтаУровень.Html())
		сам.рендер.Доб("{.шахта_режим}", сам.кнпШахтаРежим.Html())
		сам.рендер.Доб("{.шахта_сделать_кол}", сам.кнпШахтаПродукт.Html())
		сам.рендер.Доб("{.шахта_сделать_назв}", сам.кнпШахтаТип.Html())
		сам.рендер.Доб("{.шахта_сделать_время}", сам.кнпШахтаВремя.Html())
	}
	{ // Полигон
		полигон := бот.КонтБот().Get("полигон").Val().(ИАренаПолигон)
		сам.рендер.Доб("полигон_уровень", полигон.Уровень().ЗначСтр())
		сам.рендер.Доб("полигон_режим", полигон.Состояние().Получ())
		сам.рендер.Доб("полигон_сделать_кол", полигон.ПродуктСейчас().ЗначСтр())
		сам.рендер.Доб("полигон_сделать_назв", полигон.ПродуктСейчас().Имя())
		сам.рендер.Доб("полигон_сделать_время", полигон.ПродуктВремяСейчас())
	}
	{ // Арсенал
		арс := бот.КонтБот().Get("арсенал").Val().(ИАренаАрсенал)
		сам.рендер.Доб("оружейная_уровень", арс.Уровень().ЗначСтр())
		сам.рендер.Доб("оружейная_работа", арс.ПродуктСейчас().Имя())
		сам.рендер.Доб("оружейная_режим", арс.Состояние().Получ())
		сам.рендер.Доб("оружейная_кумул", арс.Кумулятивы().ЗначСтр())
		сам.рендер.Доб("оружейная_бронебойки", арс.Бронебойки().ЗначСтр())
		сам.рендер.Доб("оружейная_фугасы", арс.Фугасы().ЗначСтр())
		сам.рендер.Доб("оружейная_ремки", арс.Ремки().ЗначСтр())
		сам.рендер.Доб("оружейная_время", арс.ПродуктВремяСейчас().Get())
		сам.рендер.Доб("мощь", бот.Стата().Мощь().ЗначСтр())
		сам.рендер.Доб("мощь", бот.Стата().Мощь().ЗначСтр())
	}
	стрРез := сам.рендер.Получ()
	return кнт.SendString(стрРез)
}

/*
// Проверка на куки
func (сам *СтраницаБотПоказать) кукиПроверить(кнт *fiber.Ctx) error {
	сам.лог.Debug("СтраницаБотПоказать.кукиПроверить()\n")
	имя := кнт.Cookies("login")
	if имя != "svi" {
		return кнт.Redirect("/gui/login-show", http.StatusSeeOther)
	}
	return кнт.Next()
}
*/
