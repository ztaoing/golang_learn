package main

import "fmt"

func main() {
	m := map[string]string{
		"algorithm": "1",
		"b":         "2",
		"c":         "3",
	}
	m1 := make(map[string]int) //empty map

	var m2 map[string]int // m2=nil

	fmt.Println(m, m1, m2)
	delete(m, "algorithm")
	for k, v := range m {
		fmt.Println(k, v)
	}
	//	name := m["d"]
	if courseName, ok := m["d"]; ok {
		fmt.Println(courseName, ok)
	} else {
		fmt.Println("key dese not exsits")
	}

}

//最大不重复的子串
