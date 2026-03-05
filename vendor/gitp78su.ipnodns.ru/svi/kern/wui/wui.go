// package wui -- пакет веб-интерфейса
package wui

import (
	"gitp78su.ipnodns.ru/svi/kern/wui/wbutton"
	"gitp78su.ipnodns.ru/svi/kern/wui/wctx"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

// NewWuiButton -- возвращает новую WUI-кнопку
func NewWuiButton(text string, fnClick func(map[string]string) string) IWuiButton {
	btn := wbutton.NewWuiButton(text, fnClick)
	return btn
}

// GetWuiCtx -- возвращает контекст WUI
func GetWuiCtx() IWuiCtx {
	wCtx := wctx.GetWuiCtx()
	return wCtx
}
