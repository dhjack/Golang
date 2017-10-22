package sharemem

// #include "sharemem.h"
import "C"

import (
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"unsafe"
)

// 共享内存的控制结构
type ShmemPtr struct {
	enable bool //共享内存的数据是否可用
}

// 获取当前环境下，对齐的大小
const Align = unsafe.Sizeof(uintptr(0))

// ShememPtr 加上对齐后的大小
var controlBufSize uintptr

// 计算对齐后的大小
func GetAlignSize(size uintptr) uintptr {
	return (size + Align - 1) & (^(Align - 1))
}

func init() {
	// 计算共享内存的标志消息占用的内存大小
	controlBufSize = GetAlignSize(unsafe.Sizeof(ShmemPtr{}))
}

// 获取共享内存
// 实现用创建的方式获取。
// 1. 如果获取失败，判断是否是已经存在。
// 2. 如果是存在，直接获取
func GetShareMemory(key uintptr, size uintptr) *ShmemPtr {

	size += controlBufSize

	logger.Debug(fmt.Sprintf("sharemem key:%v, size:%v", key, size))

    shmaddr := C.getShareMemory(C.int(key), C.int(size))

	logger.Debug(fmt.Sprintf("sharemem shaddr: 0x%x", shmaddr))

	logger.Info("sharemem get share memory success!")

	return (*ShmemPtr)(unsafe.Pointer(shmaddr))
}

// 判断共享内存的数据是否可用
func (self *ShmemPtr) IsEnable() bool {
	return self.enable
}

// 设置共享内存数据可用
func (self *ShmemPtr) SetEnable() {
	self.enable = true
}

// 获取实际的buf地址
func (self *ShmemPtr) GetRawBuf() uintptr {
	return uintptr(unsafe.Pointer(self)) + controlBufSize
}

// 根据buf地址，找对应的控制结构指针
func GetControlPtr(ptr uintptr) *ShmemPtr {
	return (*ShmemPtr)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) - controlBufSize))
}
