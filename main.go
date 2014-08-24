// 23 august 2014

package main

import (
	"fmt"
	"unsafe"
)

// #cgo LDFLAGS: -lasound
// #include <alsa/asoundlib.h>
// char *cdefault = "default";
import "C"

func main() {
	var pcm *C.snd_pcm_t

	chkerror := func(err C.int, when string) {
		if err != 0 {
			panic(fmt.Errorf("error %s: %s", when, C.GoString(C.snd_strerror(err))))
		}
	}

	chkerror(C.snd_pcm_open(&pcm, C.cdefault, C.SND_PCM_STREAM_PLAYBACK, 0), "opening the PCM stream")
	defer func() {
		chkerror(C.snd_pcm_close(pcm), "closing PCM stream")
	}()
	chkerror(C.snd_pcm_set_params(pcm, C.SND_PCM_FORMAT_U8, C.SND_PCM_ACCESS_RW_INTERLEAVED,
		1, 44100, 0, 250 * 1000), "setting up PCM stream")

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

	for i := 0; i < len(buffer) * 20; {
		written := C.snd_pcm_writei(pcm, unsafe.Pointer(&buffer[i % len(buffer)]),
			C.snd_pcm_uframes_t(len(buffer) - i % len(buffer)))
		if written <= 0 {
			chkerror(C.snd_pcm_recover(pcm, C.int(written), 1), "recovering from PCM write error")
		} else {
			i += int(written)
		}
	}
}
