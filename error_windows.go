// 30 august 2014

package pcmio

import (
	"fmt"
	"syscall"
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

// TODO add when

type waveOutError C.MMRESULT

func (err waveOutError) Error() string {
	errtext := make([]uint16, C.MAXERRORLENGTH + 1)
	converr := C.waveOutGetErrorTextW(C.MMRESULT(err), C.LPWSTR(unsafe.Pointer(&errtext[0])), C.MAXERRORLENGTH + 1)
	if converr != C.MMSYSERR_NOERROR {
		return fmt.Sprintf("waveOut error 0x%X (error 0x%X converting to string)", err, converr)
	}
	return syscall.UTF16ToString(errtext)
}
