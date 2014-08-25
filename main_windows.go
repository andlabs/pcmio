// 24 august 2014

package main

import (
	"fmt"
	"unsafe"
)

// #cgo LDFLAGS: -lwinmm
// #include <windows.h>
import "C"

func main() {
	var hwo C.HWAVEOUT
	var format C.WAVEFORMATEX
	var header C.WAVEHDR

	checkerr := func(err C.MMRESULT, when string) {
		if err != C.MMSYSERR_NOERROR {
			panic(fmt.Sprintf("error %s: %d", when, err))	// TODO
		}
	}

	format.wFormatTag = C.WAVE_FORMAT_PCM
	format.nChannels = 1
	format.nSamplesPerSec = 44100
	format.wBitsPerSample = 8
	format.nBlockAlign = format.nChannels * format.wBitsPerSample
	format.nAvgBytesPerSec = format.nSamplesPerSec * C.DWORD(format.nBlockAlign)
	format.cbSize = 0
	checkerr(C.waveOutOpen(&hwo, C.WAVE_MAPPER, &format,
		0, 0, C.CALLBACK_NULL | C.WAVE_FORMAT_DIRECT), "opening PCM stream")
	defer func() {
		checkerr(C.waveOutClose(hwo), "closing PCM stream")
	}()

	buffer := make([]byte, 256*440)
	k := 0
	for i := 0; i < 440; i++ {
		min, max, step := 0, 256, 1
		if i % 2 == 1 {
			min, max, step = 255, -1, -1
		}
		for j := min; j != max; j += step {
			buffer[k] = byte(j)
			k++
		}
	}

	for i := 0; i < 20; i++ {
		header.lpData = C.LPSTR(unsafe.Pointer(&buffer[0]))
		header.dwBufferLength = C.DWORD(len(buffer))
		header.dwFlags = 0
		header.dwLoops = 0
		checkerr(C.waveOutPrepareHeader(hwo, &header, C.UINT(unsafe.Sizeof(header))),
			"pareparing PCM data for playack")
		checkerr(C.waveOutWrite(hwo, &header, C.UINT(unsafe.Sizeof(header))),
			"playing PCM data")
	}
}
