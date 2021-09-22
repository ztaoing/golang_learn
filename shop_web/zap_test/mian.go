/**
* @Author:zhoutao
* @Date:2021/9/9 1:38 下午
* @Desc:
 */

package main

import "go.uber.org/zap"

func main() {
	logger, _ := zap.NewProduction() // 生产环境的用法
	// logger, _ := zap.NewDevelopment() // 测试环境的用法
	defer logger.Sync() // 在结束之前把缓存刷到指定的地方。flushes buffer, if any
	url := "https://imooc.com"
	sugar := logger.Sugar()

	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}
