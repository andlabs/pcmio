// 30 august 2014

package pcmio

import (
	"unsafe"
)

// #cgo CFLAGS: --std=c99
// #cgo LDFLAGS: -lasound
// #include <alsa/asoundlib.h>
// char *cdefault = "default";
import "C"

type writer struct {
	pcm		*C.snd_pcm_t
}

var formats = map[Format]C.snd_pcm_format_t{
	U8:		C.SND_PCM_FORMAT_U8,
}

func openDefaultWriter(format Format, rate uint) (*writer, error) {
	w := new(writer)
	err := C.snd_pcm_open(&w.pcm, C.cdefault, C.SND_PCM_STREAM_PLAYBACK, 0)
	if err != 0 {
		return nil, alsaerr(err)
	}
	// the 0 parameter indicates no automatic rate conversions
	err = C.snd_pcm_set_params(w.pcm, formats[format], C.SND_PCM_ACCESS_RW_INTERLEAVED, 1, C.uint(rate), 0, 250 * 1000)
	if err != 0 {
		return nil, alsaerr(err)
	}
	return w, nil
}

func (w *writer) Write(p []byte) (n int, err error) {
	for i := 0; i < len(p); {
		written := C.snd_pcm_writei(w.pcm, unsafe.Pointer(&p[i]), C.snd_pcm_uframes_t(len(p) - i))
		if written <= 0 {
			err := C.snd_pcm_recover(w.pcm, C.int(written), 1)
			if err != 0 {
				return i, alsaerr(err)
			}
		} else {
			i += int(written)
		}
	}
	return len(p), nil
}
