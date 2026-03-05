// package btn_monolit -- обработчик для показа блока монолита
package btn_monolit

import (
	_ "embed"

	"gitp78su.ipnodns.ru/svi/kern/wui"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

type BtnMonolit struct {
	btn IWuiButton
}

// NewBtnMonolit -- возвращает новую кнопку монолита
func NewBtnMonolit() *BtnMonolit {
	sf := &BtnMonolit{}
	sf.btn = wui.NewWuiButton("Monolit", sf.clickMonolit)
	sf.btn.Hx().Target().Set("#monolit")
	return sf
}

// Html -- возвращает HTML-представление кнопки
func (sf *BtnMonolit) Html() string {
	return sf.btn.Html()
}

//go:embed block_monolit.html
var strBlockMonolit string

// Событие клика по кнопке
func (sf *BtnMonolit) clickMonolit(dict map[string]string) string {
	return strBlockMonolit
}
