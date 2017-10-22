package sharemem

import (
    "testing"
    "github.com/donnie4w/go-logger/logger"
	"errors"
	"fmt"
	"unsafe"
    . "sharemem"

)

type someData struct {
	data1          int64
	data2           int64
	capacity        int64
	data3             []uint8
	data4             []uint8
	data5 int32
}

func TestSHAREMEM() {
    cache, err := newUidOnlineCache(10)
    if err != nil {
        fmt.Printf("test share mem is not enable: %v\n", err)
    }
    cache.enable()
    fmt.Printf("test set share mem enable\n")
}

func newUidOnlineCache(size int64) (cache *someData, err error) {

	// 获取someData的大小
	cacheSize := unsafe.Sizeof(someData{})

    fmt.Printf("cacheSize:%v\n", cacheSize)

	// 将申请的内存大小按对齐大小改正
	capacity := GetAlignSize(uintptr(size))

    fmt.Printf("capacity:%v\n", capacity)

    g_shareMemoryKey := uintptr(0x80)
	// 获取cache的地址
	controlBuf := GetShareMemory(g_shareMemoryKey, cacheSize)

	// 获取set切片对应的数组的地址
	data3Buf := GetShareMemory(g_shareMemoryKey+1, capacity)

	// 获取data4切片对应的数组的地址
	data4Buf := GetShareMemory(g_shareMemoryKey+2, capacity)

	// 加上偏移，得到cache地址
	cache = (*someData)(unsafe.Pointer(controlBuf.GetRawBuf()))

	// slice有三个变量，第一个是数组的地址，第二个是当前的长度，第三个是数组的容量。所以一一赋值
	mapToSlice := func(buf *ShmemPtr, offset uintptr) {
		*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(controlBuf)) + Align*offset)) = buf.GetRawBuf()
		*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(controlBuf)) + Align*(offset+1))) = capacity
		*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(controlBuf)) + Align*(offset+2))) = capacity
	}

	mapToSlice(data3Buf, 7)
	mapToSlice(data4Buf, 10)

    fmt.Printf("cache before :%v\n", cache)
	cache.data5 = 0

	if !controlBuf.IsEnable() {
		logger.Info("share memory is unable, need sync!")
		cache.data1 = 0
		cache.data2 = 0
		cache.capacity = int64(capacity)
		err = errors.New("share memory is unable")
	} else {
		logger.Info("share memory is enable without sync!")
	}

    fmt.Printf("cache after:%v\n", cache)

	return
}

func (sc *someData) enable() {
	GetControlPtr(uintptr(unsafe.Pointer(sc))).SetEnable()
	fmt.Println("sync done. share memory is enable!")
}
