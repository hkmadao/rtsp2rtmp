package utils

import "github.com/deepch/vdk/av"

func OrDonePacket(done <-chan interface{}, c <-chan *av.Packet) <-chan *av.Packet {
	valStream := make(chan *av.Packet)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok { // 外界关闭数据流
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

func ToPacket(done <-chan interface{}, valueStream <-chan interface{}) <-chan av.Packet {
	stringStream := make(chan av.Packet)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(av.Packet):
			}
		}
	}()
	return stringStream
}
