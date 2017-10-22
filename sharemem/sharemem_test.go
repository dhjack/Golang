package sharemem

import (
    "testing"
    "github.com/donnie4w/go-logger/logger"
	"errors"
	"fmt"
	"unsafe"
    . "sharemem"

)

type uidOnlineCache struct {
	maxuid          int64
	count           int64
	capacity        int64
	set             []uint8
	rom             []uint8
	httpCurReqCount int32
}

func TestSHAREMEM() {
    cache, err := newUidOnlineCache(10)
    if err != nil {
        fmt.Printf("test share mem is not enable: %v\n", err)
    }
    cache.enable()
    fmt.Printf("test set share mem enable\n")
}

func newUidOnlineCache(size int64) (cache *uidOnlineCache, err error) {

	// 获取uidOnlineCache的大小
	cacheSize := unsafe.Sizeof(uidOnlineCache{})

    fmt.Printf("cacheSize:%v\n", cacheSize)

	// 将申请的内存大小按对齐大小改正
	capacity := GetAlignSize(uintptr(size))

    fmt.Printf("capacity:%v\n", capacity)

    g_shareMemoryKey := uintptr(0x80)
	// 获取cache的地址
	controlBuf := GetShareMemory(g_shareMemoryKey, cacheSize)

	// 获取set切片对应的数组的地址
	setBuf := GetShareMemory(g_shareMemoryKey+1, capacity)

	// 获取rom切片对应的数组的地址
	romBuf := GetShareMemory(g_shareMemoryKey+2, capacity)

	// 加上偏移，得到cache地址
	cache = (*uidOnlineCache)(unsafe.Pointer(controlBuf.GetRawBuf()))

	// slice有三个变量，第一个是数组的地址，第二个是当前的长度，第三个是数组的容量。所以一一赋值
	mapToSlice := func(buf *ShmemPtr, offset uintptr) {
		*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(controlBuf)) + Align*offset)) = buf.GetRawBuf()
		*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(controlBuf)) + Align*(offset+1))) = capacity
		*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(controlBuf)) + Align*(offset+2))) = capacity
	}

	mapToSlice(setBuf, 7)
	mapToSlice(romBuf, 10)

    fmt.Printf("cache before :%v\n", cache)
	cache.httpCurReqCount = 0

	if !controlBuf.IsEnable() {
		logger.Info("share memory is unable, need sync!")
		cache.maxuid = 0
		cache.count = 0
		cache.capacity = int64(capacity)
		err = errors.New("share memory is unable")
	} else {
		logger.Info("share memory is enable without sync!")
	}

    fmt.Printf("cache after:%v\n", cache)

	return
}

func (sc *uidOnlineCache) enable() {
	GetControlPtr(uintptr(unsafe.Pointer(sc))).SetEnable()
	fmt.Println("sync done. share memory is enable!")
}
