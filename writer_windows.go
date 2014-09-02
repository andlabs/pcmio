// 30 august 2014

package pcmio

import (
	"syscall"
	"unsafe"
)

// #cgo CFLAGS: --std=c99
// #cgo LDFLAGS: -lkernel32 -lwinmm
// #include "winapi_windows.h"
import "C"

type writer struct {
	wo		C.HWAVEOUT
	event	C.HANDLE
}

func openDefaultWriter(format Format, rate uint) (*writer, error) {
	var fmt C.WAVEFORMATEX

	w := new(writer)
	lasterr := C.xCreateEvent(&w.event)
	if lasterr != 0 {
		return nil, syscall.Errno(lasterr)
	}

	fmt.wFormatTag = C.WAVE_FORMAT_PCM
	fmt.nChannels = 1
	fmt.nSamplesPerSec = C.DWORD(rate)
	switch format {
	case U8:
		fmt.wBitsPerSample = 8
	}
	fmt.nBlockAlign = (fmt.nChannels * fmt.wBitsPerSample) / 8
	fmt.nAvgBytesPerSec = fmt.nSamplesPerSec * C.DWORD(fmt.nBlockAlign)
	fmt.cbSize = 0
	err := C.waveOutOpen(&w.wo, C.WAVE_MAPPER, &fmt,
		C.DWORD_PTR(uintptr(unsafe.Pointer(w.event))), 0,
		C.CALLBACK_EVENT | C.WAVE_FORMAT_DIRECT)
	if err != C.MMSYSERR_NOERROR {
		return nil, waveOutError(err)
	}

	return w, nil
}

func (w *writer) Write(p []byte) (n int, err error) {
	var hdr C.WAVEHDR

	hdr.lpData = C.LPSTR(unsafe.Pointer(&p[0]))
	hdr.dwBufferLength = C.DWORD(len(p))
	hdr.dwFlags = 0
	hdr.dwLoops = 0
	lasterr := C.xResetEvent(w.event)
	if lasterr != 0 {
		return 0, syscall.Errno(lasterr)
	}
	xerr := C.waveOutPrepareHeader(w.wo, &hdr, C.UINT(unsafe.Sizeof(hdr)))
	if xerr != C.MMSYSERR_NOERROR {
		return 0, waveOutError(xerr)
	}
	xerr = C.waveOutWrite(w.wo, &hdr, C.UINT(unsafe.Sizeof(hdr)))
	if xerr != C.MMSYSERR_NOERROR {
		return 0, waveOutError(xerr)
	}
	lasterr = C.xWaitForSingleObject(w.event)
	if lasterr != 0 {
		return 0, syscall.Errno(lasterr)
	}
	xerr = C.waveOutUnprepareHeader(w.wo, &hdr, C.UINT(unsafe.Sizeof(hdr)))
	if xerr != C.MMSYSERR_NOERROR {
		return 0, waveOutError(xerr)
	}
	return len(p), nil
}
