package utils

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math"
)

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)

	return math.Float32frombits(bits)
}

func Float64ToByteLittleEndian(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

func Float64ToByteBigEndian(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)

	return bytes
}

func Int32ToByteBigEndian(number int32) []byte {
	bytes := make([]byte, 4)
	bytes[0] = byte(number >> (3 * 8))
	bytes[0] = byte(number >> (2 * 8))
	bytes[0] = byte(number >> (1 * 8))
	bytes[0] = byte(number)

	return bytes
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}

func ReverseString(s string) string {
	// 将字符串转换为rune类型的切片
	runes := []rune(s)
	// 获取字符串长度
	n := len(runes)
	// 遍历rune类型的切片，交换前后两个元素的位置
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}
	// 将rune类型的切片转换为字符串类型并返回
	return string(runes)
}

func Md5(unMd5Str string) (md5Str string) {
	md5Str = fmt.Sprintf("%x", md5.Sum([]byte(unMd5Str)))
	return
}
