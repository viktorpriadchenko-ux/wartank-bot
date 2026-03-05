package btn_modules

import (
	_ "embed"
	"fmt"
	"strings"

	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/wui"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

type BtnModules struct {
	btn  IWuiButton
	kCtx IKernelCtx
}

// NewBtnModules -- возвращает новую кнопку модулей
func NewBtnModules() *BtnModules {
	sf := &BtnModules{
		kCtx: kctx.GetKernelCtx(),
	}
	sf.btn = wui.NewWuiButton("Modules", sf.clickMonolit)
	sf.btn.Hx().Target().Set("#modules")
	return sf
}

// Html -- возвращает HTML-представление кнопки
func (sf *BtnModules) Html() string {
	return sf.btn.Html()
}

//go:embed block_modules.html
var strBlockModules string

//go:embed block_row.html
var strBlockRow string

// Событие клика по кнопке
func (sf *BtnModules) clickMonolit(dict map[string]string) string {
	mon := sf.kCtx.Get("monolit").Val().(IKernelMonolit)
	chLst := mon.Ctx().SortedList()
	strOut := ``
	for _, val := range chLst {
		if !strings.Contains(val.Key(), "module_") {
			continue
		}
		lstKey := strings.Split(val.Key(), "_")
		id := lstKey[1]
		strRow := strBlockRow
		strRow = strings.ReplaceAll(strRow, "{.id}", id)
		strRow = strings.ReplaceAll(strRow, "{.key}", val.Key())
		moduleName := string(val.Val().(IKernelModule).Name())
		strRow = strings.ReplaceAll(strRow, "{.name}", moduleName)
		type_ := fmt.Sprintf("%#T", val.Val())
		type_ = strings.ReplaceAll(type_, ".", ".<br>")
		strRow = strings.ReplaceAll(strRow, "{.type}", type_)
		strRow = strings.ReplaceAll(strRow, "{.createAt}", string(val.CreateAt()))
		strRow = strings.ReplaceAll(strRow, "{.updateAt}", string(val.UpdateAt()))
		strRow = strings.ReplaceAll(strRow, "{.comment}", val.Comment())
		strOut += strRow
	}
	strOut = strings.ReplaceAll(strBlockModules, "{.mod_block}", strOut)
	return strOut
}
