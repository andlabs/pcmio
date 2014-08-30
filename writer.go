// 30 august 2014

package pcmio

// Writer represents a synchronous stream to an output device, such as a speaker.
type Writer interface {
	// Write behaves like io.Writer, writing the data in p (which should have the format chosen when the Writer was opened) to the output device.
	// It has one special rule: err will never be nil if n != len(p).
	Write(p []byte) (n int, err error)
}

// OpenDefaultWriter opens the default output device for writing.
// It returns an error if the specified sample format and rate cannot be used with the default output device.
// Package pcmio does not do any on-the-fly sample format conversion, even if the the underlying audio subsystem can.
func OpenDefaultWriter(format Format, rate uint) (Writer, error) {
	return openDefaultWriter(format, rate)
}
