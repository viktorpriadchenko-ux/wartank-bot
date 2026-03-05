// package wwidget -- базовый виджет WUI
package wwidget

import (
	"crypto/rand"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

// WuiWidget -- базовый виджет WUI
type WuiWidget struct {
	id string
}

// NewWuiWidget -- возвращает новый базовый виджет WUI
func NewWuiWidget() *WuiWidget {
	sf := &WuiWidget{
		id: "wui_" + rand.Text(),
	}
	_ = IWuiWidget(sf)
	return sf
}

// Id - возвращает ID виджета
func (sf *WuiWidget) Id() string {
	return sf.id
}

const (
	strBeg = `<div id="{.id}"> WuiWidget.Html(): id={.id}, not implemented </div>`
)

// Html -- возвращает HTML представление виджета
func (sf *WuiWidget) Html() string {
	strRes := strings.ReplaceAll(strBeg, "{.id}", sf.id)
	return strRes
}
