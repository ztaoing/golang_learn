package main

// 只支持 int 即可。

type MyMap struct {
}

func (m *MyMap) Load(key interface{}) (value interface{}, ok bool) {
	return nil, false
}

func (m *MyMap) Store(key, value interface{}) {
}

func (m *MyMap) Delete(key interface{}) {

}

func (m *MyMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	return nil, false
}

func (m *MyMap) LoadAndDelete(key interface{}) (value interface{}, loaded bool) {
	return nil, false
}

func main() {}
