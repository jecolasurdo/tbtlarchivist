package utils

import "encoding/binary"

type frameState int

const (
	frameStateReadingHeader frameState = iota
	frameStateReadingBody
)

const headerSize = 4

// FrameScanner exposes a ScanFrames method which can be used as a
// bufio.SplitFunc.  See FrameScanner.ScanFrames for more details.
type FrameScanner struct {
	state     frameState
	frameSize int
}

// ScanFrames is a SplitFunc that terminates records based on a record length
// frame header. The first 4 bytes of each frame is expected to be a bigendian
// int32 that denotes the length of the record to follow.  See also:
// https://golang.org/pkg/bufio/#SplitFunc
func (fs *FrameScanner) ScanFrames(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 && atEOF == false {
		return
	}

	switch fs.state {
	case frameStateReadingHeader:
		if len(data) < headerSize {
			return
		}
		fs.frameSize = int(binary.BigEndian.Uint32(data[0:headerSize]))
		advance = headerSize
		fs.state = frameStateReadingBody
	case frameStateReadingBody:
		if len(data) < fs.frameSize {
			return
		}
		token = data[0:fs.frameSize]
		advance = fs.frameSize
		fs.frameSize = 0
		fs.state = frameStateReadingHeader
	}

	return
}
