// package wtext -- WUI текст
package wtext

import (
	"strings"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
	"gitp78su.ipnodns.ru/svi/kern/wui/wwidget"
)

// WuiText -- текст для WUI
type WuiText struct {
	IWuiWidget
	sync.RWMutex
	val string
}

// NewWuiText -- возвращает новый текст WUI
func NewWuiText(val string) *WuiText {
	sf := &WuiText{
		IWuiWidget: wwidget.NewWuiWidget(),
		val:        val,
	}
	_ = IWuiText(sf)
	return sf
}

const (
	strBeg = `<span id="{.id}">{.txt}</span>`
)

// Html -- возвращает HTML-представление текста
func (sf *WuiText) Html() string {
	sf.RLock()
	defer sf.RUnlock()
	strRes := strings.ReplaceAll(strBeg, "{.id}", sf.Id())
	strRes = strings.ReplaceAll(strRes, "{.txt}", sf.val)
	return strRes
}

// Set -- устанавливает хранимое значение
func (sf *WuiText) Set(val string) {
	sf.Lock()
	defer sf.Unlock()
	sf.val = val
}

// Get -- возвращает хранимое значение
func (sf *WuiText) Get() string {
	sf.RLock()
	defer sf.RUnlock()
	return sf.val
}
