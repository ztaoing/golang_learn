/**
* @Author:zhoutao
* @Date:2022/11/21 14:08
* @Desc:
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// 在客户端请求时需要设置header，则需要使用NewRequest和DefaultClient.Do

func main() {
	req, err := http.NewRequest("GET", "http://www.baidu.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("key", "value")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	byts, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(byts))
}
