# SVI

Шипков В.И.

Заметки, чтобы не забыть.

## Ускорение фильтра Блума

Фильтр Блума не может откатывать признаки, при удалении ключа.

**Идея**: Можно просчитывать фильтр Блума в зафиксированных слоях LSM. Тогда количество
сложений фильтров из отдельных слоёв резко сокращается. Фактически фильтр Блума будет пересчитываться только в горячем слое L0.

### Работа фильтра Блума

Все частные хэши ключей складываются по OR в общую маску.
Если частный хэш ключа наложен на общую маску

### Недостатки фильтра Блума

Для различения ключей в количестве 1 млрд (при наличии нормальной хэш-функции) --
надо иметь около 1.8 млрд бит. Т.е. около 200 МБ памяти. Сканировать такой объём памяти
за разумные сроки -- не реально. Здесь проще пойти по пути префиксного дерева.

## Иерархия ключей

<img src="./nosql_keys_ierarchy.png" alt="drawing" width="200"/>

## Интерфейсы иерархии ключей

<img src="./nosql_model_interfaces.png" alt="drawing" width="600"/>

## Структура индексного файла

## Устройство NoSQL

![Схема БД](./nosql_model_lsm.png)

## Запуск zapret

```bash
export NFQWS_OPT="
--filter-tcp=80 --dpi-desync=fake,multisplit --dpi-desync-ttl=0 --dpi-desync-fooling=md5sig,badsum <HOSTLIST> --new
--filter-tcp=443 --dpi-desync=fake,multidisorder --dpi-desync-split-pos=method+2,midsld,5 --dpi-desync-ttl=0 --dpi-desync-fooling=md5sig,badsum,badseq --dpi-desync-repeats=15 --dpi-desync-fake-tls=/opt/zapret/files/fake/tls_clienthello_www_google_com.bin <HOSTLIST> --new
--filter-udp=443 --dpi-desync=fake --dpi-desync-repeats=15 --dpi-desync-ttl=0  --dpi-desync-any-protocol --dpi-desync-cutoff=d4 --dpi-desync-fooling=md5sig,badsum --dpi-desync-fake-quic=/opt/zapret/files/fake/quic_initial_www_google_com.bin <HOSTLIST>
" & sudo /opt/zapret/init.d/sysv/zapret start
```

## Языковые модели

https://www.youtube.com/watch?v=dbSyu26x4aM 1С песня

```text
llama3.1:8b
qwen2.5-coder:1.5b-base
qwen2.5-coder:7b (на тестах заявлено, что это самая крутая модель!!!)
~ (долго думает) deepseek-r1:8b
```

Запуск сервера `ollama`:

```bash
ollama serve
```

## VPN

### Запуск

```bash
sudo wg-quick up wg0
```

### Отключение

```bash
sudo wg-quick down wg0
```
