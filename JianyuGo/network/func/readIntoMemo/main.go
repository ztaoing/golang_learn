package main

import (
	"io/ioutil"
	"log"
)

func main() {
	// 读取文件到byte slice中
	data, err := ioutil.ReadFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Data read: %s\n", data)
}
