package main

import (
	"fmt"
	"time"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32        = windows.NewLazySystemDLL("user32.dll")
	procKeyState  = user32.NewProc("GetAsyncKeyState")
	procMapVKey   = user32.NewProc("MapVirtualKeyW")
	procToUnicode = user32.NewProc("ToUnicode")
)

const (
	mapVkVkToChar = 2
)

func getAsyncKeyState(vKey int) bool {
	state, _, _ := procKeyState.Call(uintptr(vKey))
	return state&0x8000 != 0
}

func vKeyToString(vKey int) string {
	scanCode, _, _ := procMapVKey.Call(uintptr(vKey), mapVkVkToChar)

	var buf [4]uint16
	keyState := make([]byte, 256)
	_, _, _ = procToUnicode.Call(
		uintptr(vKey),
		uintptr(scanCode),
		uintptr(unsafe.Pointer(&keyState[0])),
		uintptr(unsafe.Pointer(&buf[0])),
		4,
		0,
	)

	// Decode UTF-16 buffer to a string
	return string(utf16.Decode(buf[:]))
}

func main() {
	fmt.Println("press any key")

	for {
		for vKey := 0; vKey < 256; vKey++ {
			if getAsyncKeyState(vKey) {
				key := vKeyToString(vKey)
				if key != "" {
					fmt.Printf("Key pressed: %s(VK code: %d)\n", key, vKey)
				} else {
					fmt.Printf("Key pressed: VK code %d\n", vKey)
				}
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}
