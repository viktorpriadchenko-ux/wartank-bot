package btn_hard

import (
	"fmt"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/wui"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev3/serv_web/web_gui/page_bot_show/dict_bot"
)

type BtnHard struct {
	ctx     IKernelCtx
	Btn     IWuiButton
	app     ИПриложение
	dictBot *dict_bot.DictBot
}

func NewBtnHard() *BtnHard {
	ctx := GetKernelCtx()
	sf := &BtnHard{
		ctx:     ctx,
		app:     ctx.Get("мод_сервер").Val().(ИПриложение),
		dictBot: dict_bot.GetDictBot(),
	}
	sf.Btn = wui.NewWuiButton("Прочность:", sf.click)
	sf.Btn.Hx().Trigger().Set("every 5s")
	sf.Btn.Hx().Swap().Set("outerHTML")
	return sf
}

func (sf *BtnHard) Html() string {
	return sf.Btn.Html()
}

func (sf *BtnHard) click(dict map[string]string) string {
	bot := sf.dictBot.Get(dict)
	if bot == nil {
		sf.Btn.Text().Set("Прочность: не бота")
		return sf.Btn.Html()
	}
	strOut := fmt.Sprintf("[Прочность: %v]", bot.Стата().Прочность().ЗначСтр())
	sf.Btn.Text().Set(strOut)
	return sf.Btn.Html()
}
