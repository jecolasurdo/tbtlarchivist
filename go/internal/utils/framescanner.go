package utils

import "encoding/binary"

// FrameScanner exposes a ScanFrames method which can be used as a
// bufio.SplitFunc.  See FrameScanner.ScanFrames for more details.
type FrameScanner struct {
}

// ScanFrames is a SplitFunc that terminates records based on a record length
// frame header. The first 4 bytes of each frame is expected to be a bigendian
// int32 that denotes the length of the record to follow.  See also:
// https://golang.org/pkg/bufio/#SplitFunc
func (fs *FrameScanner) ScanFrames(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 && atEOF == false {
		return
	}

	if len(data) < 4 {
		advance = 4 - len(data)
	}

	if len(data) >= 4 {
		advance = int(binary.BigEndian.Uint32(data[0:4]))
	}

	return
}
