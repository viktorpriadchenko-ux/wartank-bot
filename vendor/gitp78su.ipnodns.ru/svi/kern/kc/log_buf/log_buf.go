// package log_buf -- потокобезопасный буфер лога
package log_buf

import (
	"fmt"
	"sync"

	"gitp78su.ipnodns.ru/svi/kern/kc/log_buf/log_msg"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// logBuf -- потокобезопасный буфер лога
type logBuf struct {
	sync.RWMutex
	lst    []ILogMsg
	lstErr []ILogMsg
}

// NewLogBuf -- возвращает новый потокобезопасный буфер лога
func NewLogBuf() ILogBuf {
	sf := &logBuf{
		lst:    []ILogMsg{},
		lstErr: []ILogMsg{},
	}
	return sf
}

// GetErr -- возвращает сообщение ошибки по номеру
func (sf *logBuf) GetErr(num int) ILogMsg {
	sf.RLock()
	defer sf.RUnlock()
	if len(sf.lstErr) == 0 {
		return log_msg.NewLogMsg(log_msg.DEBUG, "not error msg")
	}
	if num >= len(sf.lstErr) {
		return sf.lstErr[len(sf.lstErr)-1]
	}
	if num <= 0 {
		return sf.lstErr[0]
	}
	return sf.lstErr[num]
}

// Get -- возвращает сообщение по номеру
func (sf *logBuf) Get(num int) ILogMsg {
	sf.RLock()
	defer sf.RUnlock()
	if len(sf.lst) == 0 {
		return log_msg.NewLogMsg(log_msg.DEBUG, "*no msg*")
	}
	if num >= len(sf.lst) {
		return log_msg.NewLogMsg(log_msg.DEBUG, "*no msg*")
	}
	if num <= 0 {
		return log_msg.NewLogMsg(log_msg.DEBUG, "*no msg*")
	}
	return sf.lst[num]
}

type tMsg struct {
	text string
	args []any
}

// Debug -- сообщение отладки
func (sf *logBuf) Debug(fMsg string, args ...any) {
	sf.Lock()
	defer sf.Unlock()
	msg := tMsg{
		text: fMsg,
		args: args,
	}
	strMsg := fmt.Sprintf(msg.text, msg.args...)
	_msg := log_msg.NewLogMsg(log_msg.DEBUG, strMsg)
	sf.lst = append(sf.lst, _msg)
	sf.checkLen()
}

// Info -- информационные сообщения
func (sf *logBuf) Info(fMsg string, args ...any) {
	sf.Lock()
	defer sf.Unlock()
	msg := tMsg{
		text: fMsg,
		args: args,
	}
	strMsg := fmt.Sprintf(msg.text, msg.args...)
	_msg := log_msg.NewLogMsg(log_msg.INFO, strMsg)
	sf.lst = append(sf.lst, _msg)
	sf.checkLen()
}

// Warn -- предупреждающие сообщения
func (sf *logBuf) Warn(fMsg string, args ...any) {
	sf.Lock()
	defer sf.Unlock()
	msg := tMsg{
		text: fMsg,
		args: args,
	}
	strMsg := fmt.Sprintf(msg.text, msg.args...)
	_msg := log_msg.NewLogMsg(log_msg.WARN, strMsg)
	sf.lst = append(sf.lst, _msg)
	sf.checkLen()
}

// Err -- сообщения об ошибках
func (sf *logBuf) Err(fMsg string, args ...any) {
	sf.Lock()
	defer sf.Unlock()
	msg := tMsg{
		text: fMsg,
		args: args,
	}
	strMsg := fmt.Sprintf(msg.text, msg.args...)
	_msg := log_msg.NewLogMsg(log_msg.ERROR, strMsg)
	sf.lst = append(sf.lst, _msg)
	sf.lstErr = append(sf.lstErr, _msg)
	sf.checkLen()
	sf.checkLenErr()
}

// Size -- возвращает размер буфера
func (sf *logBuf) Size() int {
	sf.RLock()
	defer sf.RUnlock()
	return len(sf.lst)
}

// Проверяет длину общую лога
func (sf *logBuf) checkLen() {
	for len(sf.lst) > 100 {
		sf.lst = sf.lst[1:]
	}
}

// Проверяет длину лога ошибок
func (sf *logBuf) checkLenErr() {
	for len(sf.lstErr) > 100 {
		sf.lstErr = sf.lstErr[1:]
	}
}
