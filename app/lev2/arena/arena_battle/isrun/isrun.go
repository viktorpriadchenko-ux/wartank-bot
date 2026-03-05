package isrun

import "sync"

/*
	Потокобезопасный признак запуска сражения
*/

type IsRun struct {
	val   bool
	block sync.RWMutex
}

func NewIsRun() *IsRun {
	return &IsRun{}
}

func (сам *IsRun) Get() bool {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.val
}
func (сам *IsRun) Set() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.val = true
}

func (сам *IsRun) Reset() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.val = false
}
