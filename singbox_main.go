package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"fmt"
	"os"
	"sync"
	"unsafe"

	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json"
)

var (
	mu       sync.Mutex
	instance *box.Box
	cancelFn context.CancelFunc
	lastErr  *C.char
)

func setErr(s string) {
	if lastErr != nil {
		C.free(unsafe.Pointer(lastErr))
	}
	lastErr = C.CString(s)
}

func stopLocked() {
	if instance != nil {
		_ = instance.Close()
		instance = nil
	}
	if cancelFn != nil {
		cancelFn()
		cancelFn = nil
	}
}

// SingboxStart 读取 JSON 配置并启动 sing-box。返回 0 成功。
//
//export SingboxStart
func SingboxStart(home, cfg *C.char) (code C.int) {
	mu.Lock()
	defer mu.Unlock()
	defer func() {
		if r := recover(); r != nil {
			setErr(fmt.Sprint("panic: ", r))
			code = 3
		}
	}()
	stopLocked()
	_ = os.Chdir(C.GoString(home))
	data, err := os.ReadFile(C.GoString(cfg))
	if err != nil {
		setErr(err.Error())
		return 1
	}
	ctx, cancel := context.WithCancel(context.Background())
	ctx = box.Context(ctx,
		include.InboundRegistry(),
		include.OutboundRegistry(),
		include.EndpointRegistry(),
		include.DNSTransportRegistry(),
		include.ServiceRegistry())
	opts, err := json.UnmarshalExtendedContext[option.Options](ctx, data)
	if err != nil {
		cancel()
		setErr(err.Error())
		return 2
	}
	inst, err := box.New(box.Options{Context: ctx, Options: opts})
	if err != nil {
		cancel()
		setErr(err.Error())
		return 2
	}
	if err := inst.Start(); err != nil {
		_ = inst.Close()
		cancel()
		setErr(err.Error())
		return 2
	}
	instance, cancelFn = inst, cancel
	setErr("")
	return 0
}

//export SingboxStop
func SingboxStop() {
	mu.Lock()
	defer mu.Unlock()
	defer func() { _ = recover() }()
	stopLocked()
}

//export SingboxLastError
func SingboxLastError() *C.char {
	mu.Lock()
	defer mu.Unlock()
	return lastErr
}

func main() {}
