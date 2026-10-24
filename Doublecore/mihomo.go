package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/hub"
	"github.com/metacubex/mihomo/hub/executor"
)

var (
	mu      sync.Mutex
	running bool
	lastErr *C.char
)

func setErr(s string) {
	if lastErr != nil {
		C.free(unsafe.Pointer(lastErr))
	}
	lastErr = C.CString(s)
}

// MihomoStart 读取配置并启动内核。返回 0 表示成功,否则可用 MihomoLastError 取错误信息。
//
//export MihomoStart
func MihomoStart(home, cfg *C.char) (code C.int) {
	mu.Lock()
	defer mu.Unlock()
	defer func() {
		if r := recover(); r != nil {
			setErr(fmt.Sprint("panic: ", r))
			code = 3
		}
	}()
	if running {
		executor.Shutdown()
		running = false
	}
	cfgPath := C.GoString(cfg)
	constant.SetHomeDir(C.GoString(home))
	constant.SetConfig(cfgPath)
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		setErr(err.Error())
		return 1
	}
	if err := hub.Parse(data); err != nil {
		setErr(err.Error())
		return 2
	}
	running = true
	setErr("")
	return 0
}

//export MihomoStop
func MihomoStop() {
	mu.Lock()
	defer mu.Unlock()
	defer func() { _ = recover() }()
	if running {
		executor.Shutdown()
		running = false
	}
}

//export MihomoLastError
func MihomoLastError() *C.char {
	mu.Lock()
	defer mu.Unlock()
	return lastErr
}

func main() {}
