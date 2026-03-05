// package hx_url -- атрибут HTMX (URL запроса)
package hx_url

import (
	"gitp78su.ipnodns.ru/svi/kern/wui/hx_url_method"
	"gitp78su.ipnodns.ru/svi/kern/wui/hx_url_patch"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

// HxUrl -- атрибут HTMX (URL запроса)
type HxUrl struct {
	method IHxUrlMethod
	patch  IHxUrlPatch
}

// NewHxUrl -- возвращает новый URL запроса
func NewHxUrl(patch string) *HxUrl {
	sf := &HxUrl{
		method: hx_url_method.NewHxUrlMethod(),
		patch:  hx_url_patch.NewHxUrlPatch(patch),
	}
	_ = IHxUrl(sf)
	return sf
}

// String -- возвращает строковое представление тэга
func (sf *HxUrl) String() string {
	return sf.method.Get() + `="` + sf.patch.Get() + `"`
}

// Method -- возвращает метод запроса
func (sf *HxUrl) Method() IHxUrlMethod {
	return sf.method
}

// Patch -- возвращает путь запроса
func (sf *HxUrl) Patch() IHxUrlPatch {
	return sf.patch
}
