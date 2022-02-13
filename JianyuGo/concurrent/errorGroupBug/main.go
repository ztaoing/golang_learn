/**
* @Author:zhoutao
* @Date:2022/2/13 15:47
* @Desc:
 */

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	// 单独开一个协程去做其他的事情，不参与waitGroup
	go WriteChangeLog(ctx)

	for i := 0; i < 3; i++ {
		g.Go(func() error {
			return errors.New("访问redis失败\n")
		})
	}
	if err := g.Wait(); err != nil {
		fmt.Printf("appear error and err is %s", err.Error())
	}
	time.Sleep(1 * time.Second)
}

/**
因为这个ctx是WithContext方法返回的一个带取消的ctx，我们把这个ctx当作父context传入WriteChangeLog方法中了，
如果errGroup取消了，也会导致上下文的context都取消了，所以WriteChangelog方法就一直执行不到。
*/
func WriteChangeLog(ctx context.Context) error {

	select {
	case <-ctx.Done():
		return nil
	case <-time.After(time.Millisecond * 50): // 不能这样用哈！
		fmt.Println("write changelog")
	}
	return nil
}

// 运行结果
// appear error and err is 访问redis失败
