package utils

import (
	"time"

	"github.com/beego/beego/v2/core/logs"
)

func OrDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false { // 外界关闭数据流
					return
				}
				select { // 防止写入阻塞
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func Tee(done <-chan interface{}, in <-chan interface{}, timeout time.Duration) (<-chan interface{}, <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range in {
			var out1, out2 = out1, out2 // 私有变量覆盖
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case out1 <- val:
					out1 = nil // 置空阻塞机制完成select轮询
				case out2 <- val:
					out2 = nil
				case <-time.After(timeout): //设置超时发送时间，防止发送阻塞
					logs.Error("Tee timeout")
				}
			}
		}
	}()
	return out1, out2
}
