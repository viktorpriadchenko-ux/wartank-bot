// package dict_topic_serve -- словарь топиков обработчиков запросов
package dict_topic_serve

import (
	"context"
	"fmt"
	"sync"
	"time"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/kalias"
	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// dictServe -- потокобезопасный словарь обработчиков запросов
//
// Допускается только один обработчик запросов на один топик
type dictServe struct {
	sync.RWMutex
	ctx       IKernelCtx
	dictServe map[ATopic]IBusHandlerServe
}

// NewDictServe -- возвращает потокобезопасный словарь обработчиков запросов
func NewDictServe() IDictTopicServe {
	sf := &dictServe{
		ctx:       kctx.GetKernelCtx(),
		dictServe: map[ATopic]IBusHandlerServe{},
	}
	return sf
}

// Register -- регистрирует обработчик запросов
func (sf *dictServe) Register(handler IBusHandlerServe) {
	sf.Lock()
	defer sf.Unlock()
	Hassert(handler != nil, "dictServe.Register(): IBusHandlerSubscribe==nil")
	topic := handler.Topic()
	Hassert(topic != "", "dictServe.Register(): empty topic of handler")
	isTwinRegister := sf.register(handler)
	Hassert(!isTwinRegister, "dictServe.Register(): handler of topic (%v) already register", handler.Topic())
}

// Unregister -- удаляет обработчик запросов из словаря
func (sf *dictServe) Unregister(handler IBusHandlerServe) {
	sf.Lock()
	defer sf.Unlock()
	Hassert(handler != nil, "dictServe.Unregister(): IBusHandlerSubscribe==nil")
	delete(sf.dictServe, handler.Topic())
}

// SendRequest -- вызывает обработчик при поступлении запроса
func (sf *dictServe) SendRequest(topic ATopic, binReq []byte) Result[[]byte] {
	sf.RLock()
	defer sf.RUnlock()
	handler, isOk := sf.dictServe[topic]
	if !isOk {
		err := fmt.Errorf("dictServe.SendRequest(): handler for topic (%v) not exists", topic)
		return NewErr[[]byte](err)
	}
	var (
		chRes = make(chan Result[[]byte], 2)
	)
	ctx, fnCancel := context.WithTimeout(sf.ctx.Ctx(), time.Millisecond*time.Duration(TimeoutDefault))
	defer fnCancel()
	fnCall := func() {
		defer close(chRes)
		res := handler.FnBack(binReq)
		chRes <- res
	}
	go fnCall()
	select {
	case <-ctx.Done():
		err := fmt.Errorf("dictServe.SendRequest(): in call for topic (%v), err=\n\t%w", topic, ctx.Err())
		return NewErr[[]byte](err)
	case res := <-chRes:
		return res
	}
}

var TimeoutDefault = 15000

// регистрирует обработчик запросов
func (sf *dictServe) register(handler IBusHandlerServe) bool {
	topic := handler.Topic()
	_, isOk := sf.dictServe[topic]
	if isOk {
		return true
	}
	sf.dictServe[topic] = handler
	return false
}
