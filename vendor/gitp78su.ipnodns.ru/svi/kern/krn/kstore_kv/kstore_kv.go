// package kstore_kv -- локальное быстрое key-value хранилище ядра
package kstore_kv

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	"gitp78su.ipnodns.ru/svi/kern/kc/local_ctx"
	"gitp78su.ipnodns.ru/svi/kern/kc/safe_bool"
	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

const (
	storeStreamName = "kstore_kv" // Имя потока для ожидателя потоков
)

// kStoreKv -- локальное хранилище ядра
type kStoreKv struct {
	sync.RWMutex
	kCtx      IKernelCtx
	ctx       ILocalCtx
	log       ILogBuf
	wg        IKernelWg
	storePath string
	db        *badger.DB
	isWork    ISafeBool
}

var (
	kernStore *kStoreKv // Глобальный объект
	block     sync.Mutex
)

// GetKernelStore -- возвращает новое локальное хранилище ядра
func GetKernelStore() IKernelStoreKv {
	block.Lock()
	defer block.Unlock()
	log.Println("GetKernelStore()")
	if kernStore != nil {
		return kernStore
	}
	ctx := kctx.GetKernelCtx()
	sf := &kStoreKv{
		kCtx:   ctx,
		ctx:    local_ctx.NewLocalCtx(ctx.Ctx()),
		wg:     ctx.Wg(),
		isWork: safe_bool.NewSafeBool(),
	}
	sf.log = sf.ctx.Log()
	sf.open()
	kernStore = sf
	ctx.Set("kernStoreKV", kernStore, "fast KV store on Badger")
	return kernStore
}

// Log -- возвращает локальный лог
func (sf *kStoreKv) Log() ILogBuf {
	return sf.log
}

// Set -- устанавливает значение по ключу
func (sf *kStoreKv) Set(key string, val []byte) Result[bool] {
	sf.Lock()
	defer sf.Unlock()
	sf.log.Debug("kStoreKv.Set(): key='%v'", key)
	fnSet := func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), val)
		return err
	}
	err := sf.db.Update(fnSet)
	if err != nil {
		err := fmt.Errorf("kStoreKv.Set(): key=%v, err=\n\t%w", key, err)
		sf.log.Err(err.Error())
		return NewErr[bool](err)
	}
	return NewOk(true)
}

// Get -- возвращает значение по ключу
func (sf *kStoreKv) Get(key string) Result[[]byte] {
	sf.RLock()
	defer sf.RUnlock()
	sf.log.Debug("kStoreKv.Get(): key='%v'", key)
	var binVal []byte
	fnGet := func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		binVal, err = item.ValueCopy(binVal)
		return err
	}
	err := sf.db.View(fnGet)
	if err != nil {
		err := fmt.Errorf("kStoreKv.Get(): key=%v, err=\n\t%v", key, err)
		sf.log.Err(err.Error())
		return NewErr[[]byte](err)
	}
	return NewOk(binVal)
}

// Delete -- удалить ключ из хранилища
func (sf *kStoreKv) Delete(key string) Result[bool] {
	sf.Lock()
	defer sf.Unlock()
	sf.log.Debug("kStoreKv.Delete(): key='%v'", key)
	fnDelete := func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	}
	err := sf.db.Update(fnDelete)
	if err != nil {
		err := fmt.Errorf("kStoreKv.Delete(): key=%v, err=\n\t%w", key, err)
		sf.log.Err(err.Error())
		return NewErr[bool](err)
	}
	return NewOk(true)
}

// Открывает базу при создании
func (sf *kStoreKv) open() {
	sf.Lock()
	defer sf.Unlock()
	sf.log.Debug("kStoreKv.open()")
	strPath := os.Getenv("LOCAL_STORE_PATH")
	Hassert(strPath != "", "kStoreKv.open(): env LOCAL_STORE_PATH not set")
	pwd, err := os.Getwd()
	Hassert(err == nil, "kStoreKv.open(): in get PWD, err=\n\t%v", err)
	sf.storePath = pwd + strPath + "/db_local"
	err = os.MkdirAll(sf.storePath, 0750)
	Hassert(err == nil, "kStoreKv.open(): in make dir %v, err=\n\t%v", sf.storePath, err)
	sf.db, err = badger.Open(badger.DefaultOptions(sf.storePath))
	Hassert(err == nil, "kStoreKv.open(): in open DB %v, err=\n\t%v", sf.storePath, err)
	res := sf.wg.Add(storeStreamName)
	Hassert(res.IsOk(), "kStoreKv.open(): in add name stream to IKernelWg, err=\n\t%v", res.Error())
	sf.isWork.Set()
	go sf.close()
	go sf.clean()
}

// Выполняет периодическую сборку мусора в файле
func (sf *kStoreKv) clean() {
	chRun := make(chan int, 2)
	defer close(chRun)
	fnClean := func() {
		sf.Lock()
		defer sf.Unlock()
		_ = sf.db.RunValueLogGC(0.7)
	}
	chRun <- 1
	for {
		select {
		case <-sf.kCtx.Ctx().Done(): // надо прекратить работу
			return
		case <-chRun: // Пора поработать
			fnClean()
		}
		time.Sleep(time.Second * 1)
	}
}

// Ожидает последнего потока под отдельной блокировкой
func (sf *kStoreKv) wait(chWait chan int) {
	for {
		time.Sleep(time.Millisecond * 5)
		if sf.wg.Len() <= 1 {
			break
		}
	}
	close(chWait)
}

// Ожидает закрытия контекста ядра, закрывает хранилище
func (sf *kStoreKv) close() {
	sf.kCtx.Done()
	sf.Lock()
	defer sf.Unlock()
	if !sf.isWork.Get() {
		return
	}
	chWait := make(chan int, 2)
	go sf.wait(chWait)
	<-chWait
	sf.isWork.Reset()
	err := sf.db.Close()
	Assert(err == nil, "kStoreKv.close(): in close DB, err=\n\t%v", err)
	sf.wg.Done(storeStreamName)
	sf.log.Debug("kStoreKv.close(): done")
}
