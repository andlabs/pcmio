// 23 august 2014

package pcmio

import (
	"testing"
)

func TestOutput(t *testing.T) {
	w, err := OpenDefaultWriter(U8, 44100)
	if err != nil {
		t.Errorf("error opening default writer: %v", err)
	}

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
		n, err := w.Write(buffer)
		if err != nil {
			t.Errorf("error writing iteration %d: %v", i, err)
		} else if n != len(buffer) {
			t.Errorf("n (%d) != len(buffer) with nil error; iteration %d", n, i)
		}
	}
}
