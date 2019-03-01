# lcache

Встраиваемая библиотека кэширования для Go.

## Характерные особенности

* Написана на чистом Go
* Для использования встраивается в пooриложение
* Хранит данные в первоначальном виде, тем самым снижаются расходы на сериализацию данных
* Не порождает фоновых потоков
* Оптимизирована для использования в многопоточной среде и потокобезопасное
* Позволяет задать время жизни значений
* Поддержка доступа через кэш
* Поддержка кэширования только одной сущности (Atom)

## Установка

```
go get github.com/belfinor/lcache
```

## Работа с библиотекой

### Подключение

```
import "github.com/belfinor/lcache"
```

### Создание объекта кэша

```
var cache *lcache.Cache = lcache.New(&lcache.Config{TTL: 86400, Size: 102400, Nodes: 16})
```

В этом примере создается объект кэша, в котором будут хранится максимум 102400 значений, распределенные внутри по 16 нодами (для оптимизации параллельного доступа, а для синхронизации используется RWMutex, свой для каждй ноды) с временем жизни 86400 секунд.

В программе можно создавать неограниченное число объектов-кэшей индивидуальных под конкретные нужды (новый кэш не порождает новых горутин).

В случае достижения максимального количества элементов

### Манипулирование значениями в кэше

В качестве ключей кэша используются строки, а в качестве значений могут быть использованы произвольные объекты.

#### Set/Get

При помощи методов Get/Set можно задать получить значение из кэша. Если значения в кэше нет, то возвращается *nil*.

```
cache.Set("1", "11")
cache.Set("2",int64(22))

fmt.Println( cache.Get("1").(string), cache.Get("2").(int64), cache.Get("3") == nil )
```

#### Доступ через кэш (Fetch)

```
tm := cache.Fetch( "123", func(k string ) interface{} {
  return time.Now().Unix()
  }).(int64)
```

Метод *Fetch* работает следующим образом, если значение есть в кэше, то он сразу возвращает его. Если его нет, то для получения
значения используется переданная функция. В случае отличного от *nil* значения оно помещается в кэш и возвращается в качестве значения *Fetch*. Если *nil*, то в кэш ничего не кладется и возвращается *nil*.

#### Delete

Удаляет переданный ключ из кэша

```
cache.Delete("123")
```

#### Удаление всех ключей

```
cache.Flush()
```

#### Инкрементальный счетчик

```
cache.Inc("inc")
cache.Inc("inc")
cache.Inc("inc")

fmt.Println( cache.Get("inc").(int64)) // 3
```
