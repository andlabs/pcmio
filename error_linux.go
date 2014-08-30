// 30 august 2014

package pcmio

// #include <alsa/asoundlib.h>
import "C"

// TODO add when

type alsaerr C.int

func (e alsaerr) Error() string {
	return C.GoString(C.snd_strerror(C.int(e)))
}
