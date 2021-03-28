/**
* @Author:zhoutao
* @Date:2021/3/28 下午2:22
* @Desc:
 */

package simplelru

import (
	"container/list"
	"errors"
)

//非线程安全的lru
type LRU struct {
	size      int //缓存的大小
	evictList *list.List
	items     map[interface{}]*list.Element
	onEvict   EvictCallback
}

type EvictCallback func(key, val interface{})

type entry struct {
	key interface{}
	val interface{}
}

func NewLRU(size int, onEvict EvictCallback) (*LRU, error) {
	if size < 0 {
		return nil, errors.New("size must be positive")
	}
	c := &LRU{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
		onEvict:   onEvict,
	}
	return c, nil
}

//清空所有缓存
func (l *LRU) Purge() {
	for k, v := range l.items {
		if l.onEvict != nil {
			l.onEvict(k, v.Value.(*entry))
		}
		delete(l.items, k)
	}
	//清除所有元素后，初始化链表
	l.evictList.Init()
}

//return 是否有元素被清除
func (l *LRU) Add(key, val interface{}) (evicted bool) {
	//key是否已经存在
	if ent, ok := l.items[key]; ok {
		//移动到链表的头部
		l.evictList.MoveToFront(ent)
		//更新value
		ent.Value.(*entry).val = val
		return false
	}
	//添加新元素
	ent := &entry{key: key, val: val}
	//新元素添加到链表头部
	entry := l.evictList.PushFront(ent)
	//存储到map中
	l.items[key] = entry

	//如果添加元素后链表长度超过限制则清理最久没有使用的数据
	evict := l.evictList.Len() > l.size
	if evict {
		l.removeOldest()
	}
	return evict
}

func (l *LRU) Get(key interface{}) (val interface{}, ok bool) {
	if ent, ok := l.items[key]; ok {
		//将元素移动到链表头部
		l.evictList.MoveToFront(ent)
		//todo
		if ent.Value.(*entry) == nil {
			return nil, false
		}
		return ent.Value.(*entry).val, true
	}
	//不存在则返回对应元素的初始值
	return
}

//只查找，不更新位置
func (l *LRU) Contains(key interface{}) (ok bool) {
	_, ok = l.items[key]
	return ok
}

//return val 但不更新位置
func (l *LRU) Peek(key interface{}) (val interface{}, ok bool) {
	var ent *list.Element
	if ent, ok = l.items[key]; ok {
		return ent.Value.(*entry).val, true
	}
	return nil, ok
}

//清除存在的元素
func (l *LRU) Remove(key interface{}) (present bool) {
	if ent, ok := l.items[key]; ok {
		l.removeElement(ent)
		return true
	}
	return false
}

func (l *LRU) RemoveOldest() (key, val interface{}, ok bool) {
	//获得链表尾部的元素
	ent := l.evictList.Back()
	if ent != nil {
		l.removeElement(ent)
		kv := ent.Value.(*entry)
		return kv.key, kv.val, true
	}
	return nil, nil, false
}

//获得最久没有使用的元素的key val
func (l *LRU) GetOldest() (key, val interface{}, ok bool) {
	ent := l.evictList.Back()
	if ent != nil {
		kv := ent.Value.(*entry)
		return kv.key, kv.val, true
	}
	return nil, nil, false
}

func (l *LRU) Keys() []interface{} {
	keys := make([]interface{}, len(l.items))
	i := 0
	for ent := l.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	return keys
}

//链表的长度
func (l *LRU) Len() int {
	return l.evictList.Len()
}

func (l *LRU) Resize(size int) (evicted int) {
	diff := l.Len() - size
	//扩容
	if diff < 0 {
		diff = 0
	}
	//缩容，先清除部分数据
	for i := 0; i < diff; i++ {
		l.RemoveOldest()
	}
	l.size = size
	return diff
}

func (l *LRU) removeOldest() {
	ent := l.evictList.Back()
	if ent != nil {
		l.removeElement(ent)
	}
}

func (l *LRU) removeElement(e *list.Element) {
	//在链表中清除
	l.evictList.Remove(e)
	//在map中清除
	kv := e.Value.(*entry)
	delete(l.items, kv.key)
	if l.onEvict != nil {
		l.onEvict(kv.key, kv.val)
	}
}
