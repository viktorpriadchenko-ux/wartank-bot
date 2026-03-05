// package wctx -- глобальный контекст графики
package wctx

import (
	"sync"

	"gitp78su.ipnodns.ru/svi/kern/kc/local_ctx"
	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

var (
	wCtx  IWuiCtx
	block sync.Mutex
)

// GetWuiCtx -- возвращает глобальный контекст графики
func GetWuiCtx() IWuiCtx {
	block.Lock()
	defer block.Unlock()
	if wCtx != nil {
		return wCtx
	}
	kCtx := kctx.GetKernelCtx()
	wCtx = local_ctx.NewLocalCtx(kCtx.Ctx())
	return wCtx
}
