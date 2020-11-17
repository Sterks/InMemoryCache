package main

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
)

var ops uint64

type (
	Key   = string
	Value = string
)


func main()  {

	//var id uint64

	m := NewInMemoryCache()
	//m.Create("test1", "value1")
	//m.Create("test2", "value2")
	//m.Create("test3", "value3")
	//m.Create("test3", "value3")
	//m.Create("1", "value3")
	z := m.GetOrSet("test", func() Value {
		key := strconv.FormatUint(atomic.AddUint64(&ops, 1), 10)
		m.data[key] = ""
		return key
	})
	fmt.Println(z)
	//fmt.Println("-------------")
	//fmt.Println(m.Get("test3"))
	//m.GetOrSet()
}

type Cache interface {
	GetOrSet(key Key, valueFn func() Value) Value
	Get(key Key) (Value, bool)
}

// ----------------------------------------------

type InMemoryCache struct {
	dataMutex sync.RWMutex
	data      map[Key]Value
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[Key]Value),
	}
}

func (cache *InMemoryCache) Get(key Key) (Value, bool) {
	cache.dataMutex.RLock()
	defer cache.dataMutex.RUnlock()

	value, found := cache.data[key]
	return value, found
}

func (cache *InMemoryCache) Create(key string, value string){
	cache.dataMutex.RLock()
	defer cache.dataMutex.RUnlock()

	_, found := cache.Get(key)
	if found {
		key = strconv.FormatUint(atomic.AddUint64(&ops, 1), 10)
		cache.data[key] = value
	} else {
		cache.data[key] = value
	}
	fmt.Println(cache.data)
}

// GetOrSet возвращает значение ключа в случае его существования.
// Иначе, вычисляет значение ключа при помощи valueFn, сохраняет его в кэш и возвращает это значение.
func (cache *InMemoryCache) GetOrSet(key Key, valueFn func() Value) Value {
	cache.dataMutex.RLock()
	defer cache.dataMutex.RUnlock()

	_, found := cache.Get(key)
	if found {
		return cache.data[key]
	} else {
		key = valueFn()
		_ = cache.data[key]
		return key
	}
}

