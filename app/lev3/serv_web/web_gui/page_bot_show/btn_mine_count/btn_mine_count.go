package btn_mine_count

import (
	"fmt"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/wui"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/dict_bot"
)

type BtnMineCount struct {
	ctx     IKernelCtx
	Btn     IWuiButton
	app     ИПриложение
	dictBot *dict_bot.DictBot
}

func NewBtnMineCount() *BtnMineCount {
	ctx := GetKernelCtx()
	sf := &BtnMineCount{
		ctx:     ctx,
		app:     ctx.Get("мод_сервер").Val().(ИПриложение),
		dictBot: dict_bot.GetDictBot(),
	}
	sf.Btn = wui.NewWuiButton("Кол-во:", sf.click)
	sf.Btn.Hx().Trigger().Set("every 5s")
	sf.Btn.Hx().Swap().Set("outerHTML")
	return sf
}

func (sf *BtnMineCount) Html() string {
	return sf.Btn.Html()
}

func (sf *BtnMineCount) click(dict map[string]string) string {
	bot := sf.dictBot.Get(dict)
	if bot == nil {
		sf.Btn.Text().Set("Кол-во: не бота")
		return sf.Btn.Html()
	}
	шахта := bot.КонтБот().Get("шахта").Val().(ИАренаШахта)
	strOut := fmt.Sprintf("[Кол-во: %v]", шахта.ПродуктСейчас().ЗначСтр())
	sf.Btn.Text().Set(strOut)
	return sf.Btn.Html()
}
