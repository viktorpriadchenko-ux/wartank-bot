package btn_bron

import (
	"fmt"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/wui"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/dict_bot"
)

type BtnBron struct {
	ctx     IKernelCtx
	Btn     IWuiButton
	app     ИПриложение
	dictBot *dict_bot.DictBot
}

func NewBtnBron() *BtnBron {
	ctx := GetKernelCtx()
	sf := &BtnBron{
		ctx:     ctx,
		app:     ctx.Get("мод_сервер").Val().(ИПриложение),
		dictBot: dict_bot.GetDictBot(),
	}
	sf.Btn = wui.NewWuiButton("Броня:", sf.click)
	sf.Btn.Hx().Trigger().Set("every 5s")
	sf.Btn.Hx().Swap().Set("outerHTML")
	return sf
}

func (sf *BtnBron) Html() string {
	return sf.Btn.Html()
}

func (sf *BtnBron) click(dict map[string]string) string {
	bot := sf.dictBot.Get(dict)
	if bot == nil {
		sf.Btn.Text().Set("Броня: не бота")
		return sf.Btn.Html()
	}
	strOut := fmt.Sprintf("[Броня: %v]", bot.Стата().Броня().ЗначСтр())
	sf.Btn.Text().Set(strOut)
	return sf.Btn.Html()
}
