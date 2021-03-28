/**
* @Author:zhoutao
* @Date:2021/3/28 下午2:23
* @Desc:
 */

package simplelru

type LRUCache interface {
	Add(key, val interface{}) bool
	Get(key interface{}) (val interface{}, ok bool)
	//只查找，不更新位置
	Contains(key interface{}) (ok bool)
	Peek(key interface{}) (val interface{}, ok bool)
	Remove(key interface{}) bool
	RemoveOldest() (interface{}, interface{}, bool)
	GetOldest() (interface{}, interface{}, bool)
	Keys() []interface{}
	Len() int
	//清除所有缓存
	Purge()
	Resize(int) int
}
