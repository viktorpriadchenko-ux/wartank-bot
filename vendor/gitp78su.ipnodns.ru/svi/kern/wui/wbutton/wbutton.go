// package wbutton -- WUI-кнопка
package wbutton

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"

	"gitp78su.ipnodns.ru/svi/kern/wui/wctx"
	"gitp78su.ipnodns.ru/svi/kern/wui/whx"
	"gitp78su.ipnodns.ru/svi/kern/wui/wtext"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
	"gitp78su.ipnodns.ru/svi/kern/wui/wwidget"
)

// WuiButton -- WUI-кнопка
type WuiButton struct {
	IWuiWidget
	text   IWuiText
	fnBack func(map[string]string) string
	hx     IWuiHx
}

// NewWuiButton -- возвращает новую WUI-кнопку
func NewWuiButton(text string, fnBack func(map[string]string) string) *WuiButton {
	Hassert(fnBack != nil, "NewWuiButton(): fnBack==nil")
	sf := &WuiButton{
		IWuiWidget: wwidget.NewWuiWidget(),
		text:       wtext.NewWuiText(text),
		fnBack:     fnBack,
	}
	sf.hx = whx.NewWuiHx("/wui/click/" + sf.Id())
	wCtx := wctx.GetWuiCtx()
	wCtx.Set(sf.Id(), sf, "WUI-кнопка")
	_ = IWuiButton(sf)
	return sf
}

// Hx -- возвращает атрибуты HTMX
func (sf *WuiButton) Hx() IWuiHx {
	return sf.hx
}

// Text -- возвращает текст кнопки
func (sf *WuiButton) Text() IWuiText {
	return sf.text
}

// Click -- событие нажатия
func (sf *WuiButton) Click(dict map[string]string) string {
	return sf.fnBack(dict)
}

const (
	strBeg = `<span id="{.id}" class="btn btn-primary" {.hx}>{.txt}</span>`
)

// Html -- возвращает HTML-представление текста
func (sf *WuiButton) Html() string {
	strRes := strings.ReplaceAll(strBeg, "{.id}", sf.Id())
	strRes = strings.ReplaceAll(strRes, "{.txt}", sf.text.Get())
	strRes = strings.ReplaceAll(strRes, "{.hx}", sf.hx.String())
	return strRes
}
