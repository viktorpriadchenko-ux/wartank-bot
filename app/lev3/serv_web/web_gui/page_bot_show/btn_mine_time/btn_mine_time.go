package btn_mine_time

import (
	"fmt"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/wui"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/dict_bot"
)

type BtnMineTime struct {
	ctx     IKernelCtx
	Btn     IWuiButton
	app     ИПриложение
	dictBot *dict_bot.DictBot
}

func NewBtnMineTime() *BtnMineTime {
	ctx := GetKernelCtx()
	sf := &BtnMineTime{
		ctx:     ctx,
		app:     ctx.Get("мод_сервер").Val().(ИПриложение),
		dictBot: dict_bot.GetDictBot(),
	}
	sf.Btn = wui.NewWuiButton("Время:", sf.click)
	sf.Btn.Hx().Trigger().Set("every 5s")
	sf.Btn.Hx().Swap().Set("outerHTML")
	return sf
}

func (sf *BtnMineTime) Html() string {
	return sf.Btn.Html()
}

func (sf *BtnMineTime) click(dict map[string]string) string {
	bot := sf.dictBot.Get(dict)
	if bot == nil {
		sf.Btn.Text().Set("Время: не бота")
		return sf.Btn.Html()
	}
	шахта := bot.КонтБот().Get("шахта").Val().(ИАренаШахта)
	strOut := fmt.Sprintf("[Время: %v]", шахта.ПродуктВремяСейчас().Get())
	sf.Btn.Text().Set(strOut)
	return sf.Btn.Html()
}
