package main

import "sync"

//实现一个单例
type singleton struct {
}

var instance *singleton
var mu sync.Mutex

//加锁
func GetInstance() *singleton {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		instance = &singleton{}
	}
	return instance
}

//双重锁，避免了每次加锁，提高了代码效率
func GetInstances1() *singleton {
	if instance == nil {
		mu.Lock()
		defer mu.Unlock()
		if instance == nil {
			instance = &singleton{}
		}
	}
	return instance
}

//pool.Once实现
var once sync.Once

func GetInstnce2() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}
func main() {

}
