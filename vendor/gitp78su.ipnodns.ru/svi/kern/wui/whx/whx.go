// package whx -- HTMX-атрибуты WUI-объекта
package whx

import (
	"gitp78su.ipnodns.ru/svi/kern/wui/hx_swap"
	"gitp78su.ipnodns.ru/svi/kern/wui/hx_swap_oob"
	"gitp78su.ipnodns.ru/svi/kern/wui/hx_target"
	"gitp78su.ipnodns.ru/svi/kern/wui/hx_trigger"
	"gitp78su.ipnodns.ru/svi/kern/wui/hx_url"
	"gitp78su.ipnodns.ru/svi/kern/wui/hx_vals"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

// WuiHx -- HTMX-атрибуты WUI-объекта
type WuiHx struct {
	url     IHxUrl
	trigger IHxTrigger
	target  IHxTarget
	swap    IHxSwap
	oob     IHxSwapOob
	vals    IHxVals
}

// NewWuiHx -- возвращает новые атрибуты HTMX для WUI-объекта
func NewWuiHx(path string) *WuiHx {
	sf := &WuiHx{
		url:     hx_url.NewHxUrl(path),
		trigger: hx_trigger.NewHxTrigger(),
		target:  hx_target.NewHxTarget(),
		swap:    hx_swap.NewHxSwap(),
		oob:     hx_swap_oob.NewHxSwapOob(),
		vals:    hx_vals.NewHxVals(),
	}
	_ = IWuiHx(sf)
	return sf
}

// String -- возвращает строку тэгов
func (sf *WuiHx) String() string {
	strOut := sf.url.String() + " " // Не может быть пустым
	trig := sf.trigger.Get()
	if trig != "" {
		strOut += sf.trigger.String() + " "
	}
	targ := sf.target.Get()
	if targ != "" {
		strOut += sf.target.String() + " "
	}
	swap := sf.swap.Get()
	if swap != "" {
		strOut += sf.swap.String() + " "
	}
	oob := sf.oob.Get()
	if oob != "" {
		strOut += sf.oob.String() + " "
	}
	valsLen := sf.vals.Len()
	if valsLen != 0 {
		strOut += sf.vals.String()
	}
	return strOut
}

// Vals -- возвращает тэг переменных запроса
func (sf *WuiHx) Vals() IHxVals {
	return sf.vals
}

// Url -- возвращает тэг URL
func (sf *WuiHx) Url() IHxUrl {
	return sf.url
}

// Trigger -- возвращает тэг триггера запроса
func (sf *WuiHx) Trigger() IHxTrigger {
	return sf.trigger
}

// Target -- возвращает объект цели замены
func (sf *WuiHx) Target() IHxTarget {
	return sf.target
}

// Oob -- возвращает тэг внеполосной замены
func (sf *WuiHx) Oob() IHxSwapOob {
	return sf.oob
}

// Swap -- возвращает тэг замены
func (sf *WuiHx) Swap() IHxSwap {
	return sf.swap
}
