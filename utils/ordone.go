package utils

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
