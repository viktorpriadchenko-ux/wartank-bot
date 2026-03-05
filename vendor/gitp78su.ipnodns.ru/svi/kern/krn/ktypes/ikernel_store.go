package ktypes

// IKernelStoreKv -- интерфейс к локальному быстрому key-value хранилищу ядра
type IKernelStoreKv interface {
	// Get -- возвращает значение по ключу
	Get(key string) Result[[]byte]
	// Set -- устанавливает значение по ключу
	Set(key string, val []byte) Result[bool]
	// Delete -- удаляет значение по ключу
	Delete(key string) Result[bool]
	// Log -- возвращает локальный лог
	Log() ILogBuf
}
