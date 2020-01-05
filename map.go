package main

/**
考点:map初始化
map需要初始化后才能使用。
*/
type Param map[string]interface{}

type Show struct {
	*Param
}

func main() {
	s := new(Show)
	//invalid operation: s.Param["RMB"] (type *Param does not support indexing)
	s.Param["RMB"] = 10000
}
