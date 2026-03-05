package btn_gold

import (
	"fmt"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/wui"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/dict_bot"
)

type BtnGold struct {
	ctx     IKernelCtx
	Btn     IWuiButton
	app     ИПриложение
	dictBot *dict_bot.DictBot
}

func NewBtnGold() *BtnGold {
	ctx := GetKernelCtx()
	sf := &BtnGold{
		ctx:     ctx,
		app:     ctx.Get("мод_сервер").Val().(ИПриложение),
		dictBot: dict_bot.GetDictBot(),
	}
	sf.Btn = wui.NewWuiButton("Золото:", sf.click)
	sf.Btn.Hx().Trigger().Set("every 5s")
	sf.Btn.Hx().Swap().Set("outerHTML")
	return sf
}

func (sf *BtnGold) Html() string {
	return sf.Btn.Html()
}

func (sf *BtnGold) click(dict map[string]string) string {
	bot := sf.dictBot.Get(dict)
	if bot == nil {
		sf.Btn.Text().Set("Золото: не бота")
		return sf.Btn.Html()
	}
	ангар := bot.КонтБот().Get("ангар").Val().(ИАренаАнгар)
	strOut := fmt.Sprintf("[Золото: %v]", ангар.Золото().ЗначСтр())
	sf.Btn.Text().Set(strOut)
	return sf.Btn.Html()
}
