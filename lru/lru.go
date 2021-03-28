/**
* @Author:zhoutao
* @Date:2021/3/28 下午3:31
* @Desc:
 */

package lru

import (
	"golang_learn/golang_learn/lru/simplelru"
	"sync"
)

const DefaultEvictedBufferSize = 16

//线程安全的cache
type Cache struct {
	lru                        *simplelru.LRU
	evictedKeys, evictedValues []interface{}
	onEvictedCB                func(k, v interface{})
	lock                       sync.RWMutex
}

func New(size int) (*Cache, error) {
	return NewWithEvict(size, nil)
}

func NewWithEvict(size int, onEvicted func(key, val interface{})) (c *Cache, err error) {
	c = &Cache{
		onEvictedCB: onEvicted,
	}
	if onEvicted != nil {
		c.initEvictBuffers()
		onEvicted = c.onEvicted
	}
	c.lru, err = simplelru.NewLRU(size, onEvicted)
	return c, err
}

//初始化缓冲区
func (c *Cache) initEvictBuffers() {
	c.evictedKeys = make([]interface{}, 0, DefaultEvictedBufferSize)
	c.evictedValues = make([]interface{}, DefaultEvictedBufferSize)
}

//保存被清除的key和value
func (c *Cache) onEvicted(k, v interface{}) {
	c.evictedKeys = append(c.evictedKeys, k)
	c.evictedValues = append(c.evictedValues, v)
}

//lock
func (c *Cache) Purge() {
	var ks, vs []interface{}
	//lock
	c.lock.Lock()
	c.lru.Purge()
	if c.onEvictedCB != nil && len(c.evictedKeys) > 0 {
		ks, vs = c.evictedKeys, c.evictedValues
		c.initEvictBuffers()
	}
	//unlock
	c.lock.Unlock()
	//调用回调函数
	if c.onEvictedCB != nil {
		for i := 0; i < len(ks); i++ {
			c.onEvictedCB(ks[i], vs[i])
		}
	}
}

//return 添加元素时是否有元素被清除
func (c *Cache) Add(key, val interface{}) (evicted bool) {
	var k, v interface{}
	//lock
	c.lock.Lock()
	evicted = c.lru.Add(key, val)
	if c.onEvictedCB != nil && evicted {
		k, v = c.evictedKeys[0], c.evictedValues[0]
		//清空
		c.evictedKeys, c.evictedValues = c.evictedKeys[:0], c.evictedValues[:0]
	}
	c.lock.Unlock()

	if c.onEvictedCB != nil && evicted {
		c.onEvictedCB(k, v)
	}
	return
}

func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	value, ok = c.lru.Get(key)
	c.lock.Unlock()
	return value, ok
}

func (c *Cache) Contains(key interface{}) bool {
	c.lock.RLock()
	containKey := c.lru.Contains(key)
	c.lock.RUnlock()
	return containKey
}

func (c *Cache) Peek(key interface{}) (value interface{}, ok bool) {
	c.lock.RLock()
	value, ok = c.lru.Peek(key)
	c.lock.RUnlock()
	return value, ok
}

func (c *Cache) ContainsOrAdd(key, value interface{}) (ok, evicted bool) {
	var k, v interface{}
	c.lock.Lock()
	if c.lru.Contains(key) {
		c.lock.Unlock()
		return true, false
	}
	evicted = c.lru.Add(key, value)
	if c.onEvictedCB != nil && evicted {
		k, v = c.evictedKeys[0], c.evictedValues[0]
		c.evictedKeys, c.evictedValues = c.evictedKeys[:0], c.evictedValues[:0]
	}
	c.lock.Unlock()
	if c.onEvictedCB != nil && evicted {
		c.onEvictedCB(k, v)
	}
	return false, evicted
}

func (c *Cache) PeekOrAdd(key, value interface{}) (previous interface{}, ok, evicted bool) {
	var k, v interface{}
	c.lock.Lock()
	previous, ok = c.lru.Peek(key)
	if ok {
		c.lock.Unlock()
		return previous, true, false
	}
	evicted = c.lru.Add(key, value)
	if c.onEvictedCB != nil && evicted {
		k, v = c.evictedKeys[0], c.evictedValues[0]
		c.evictedKeys, c.evictedValues = c.evictedKeys[:0], c.evictedValues[:0]
	}
	c.lock.Unlock()
	if c.onEvictedCB != nil && evicted {
		c.onEvictedCB(k, v)
	}
	return nil, false, evicted
}

func (c *Cache) Remove(key interface{}) (present bool) {
	var k, v interface{}
	c.lock.Lock()
	present = c.lru.Remove(key)
	if c.onEvictedCB != nil && present {
		k, v = c.evictedKeys[0], c.evictedValues[0]
		c.evictedKeys, c.evictedValues = c.evictedKeys[:0], c.evictedValues[:0]
	}
	c.lock.Unlock()
	if c.onEvictedCB != nil && present {
		c.onEvicted(k, v)
	}
	return
}

func (c *Cache) Resize(size int) (evicted int) {
	var ks, vs []interface{}
	c.lock.Lock()
	evicted = c.lru.Resize(size)
	if c.onEvictedCB != nil && evicted > 0 {
		ks, vs = c.evictedKeys, c.evictedValues
		c.initEvictBuffers()
	}
	c.lock.Unlock()
	if c.onEvictedCB != nil && evicted > 0 {
		for i := 0; i < len(ks); i++ {
			c.onEvictedCB(ks[i], vs[i])
		}
	}
	return evicted
}

func (c *Cache) RemoveOldest() (key, value interface{}, ok bool) {
	var k, v interface{}
	c.lock.Lock()
	key, value, ok = c.lru.RemoveOldest()
	if c.onEvictedCB != nil && ok {
		k, v = c.evictedKeys[0], c.evictedValues[0]
		c.evictedKeys, c.evictedValues = c.evictedKeys[:0], c.evictedValues[:0]
	}
	c.lock.Unlock()
	if c.onEvictedCB != nil && ok {
		c.onEvictedCB(k, v)
	}
	return
}

func (c *Cache) GetOldest() (key, value interface{}, ok bool) {
	c.lock.RLock()
	key, value, ok = c.lru.GetOldest()
	c.lock.RUnlock()
	return
}

func (c *Cache) Keys() []interface{} {
	c.lock.RLock()
	keys := c.lru.Keys()
	c.lock.RUnlock()
	return keys
}

func (c *Cache) Len() int {
	c.lock.RLock()
	length := c.lru.Len()
	c.lock.RUnlock()
	return length
}
