package utils

import (
	"encoding/binary"
	"io"
)

const maxBufferSize = 1000 * 1024

type frameState int

const (
	frameStateReadingHeader frameState = iota
	frameStateReadingBody
)

// FrameScanner provides methods for scanning a reader that returns records
// that are prefixed with a big endian int32 that descrives the length of the
// message to follow.
type FrameScanner struct {
	reader           io.Reader
	backoff          *Backoff
	err              error
	state            frameState
	buffer           [maxBufferSize]byte
	bufferStart      int
	bufferEnd        int
	currentFrameSize int
}

// NewFrameScanner inststantiates a new FrameScanner.
func NewFrameScanner(reader io.ReadCloser, backoff *Backoff) *FrameScanner {
	return &FrameScanner{
		buffer:      [maxBufferSize]byte{},
		reader:      reader,
		backoff:     backoff,
		state:       frameStateReadingHeader,
		bufferStart: 0,
		bufferEnd:   maxBufferSize,
	}
}

// Poll begins polling the underlaying reader, returning each message on a
// channel.
func (fs *FrameScanner) Poll() <-chan []byte {
	recordSource := make(chan []byte)
	go func() {
		defer close(recordSource)
		for {
			n, err := fs.reader.Read(fs.buffer[fs.bufferStart:fs.bufferEnd])
			if err != nil {
				fs.err = err
				return
			}
			fs.bufferStart += n
			switch fs.state {
			case frameStateReadingHeader:
				if fs.bufferStart >= 4 {
					fs.currentFrameSize = int(binary.BigEndian.Uint32(fs.buffer[0:4]))
					fs.shiftBufferLeft(4)
					fs.state = frameStateReadingBody
				}
			case frameStateReadingBody:
				if fs.bufferStart >= fs.currentFrameSize {
					recordSource <- fs.buffer[0:fs.currentFrameSize]
					fs.shiftBufferLeft(fs.currentFrameSize)
					fs.state = frameStateReadingHeader
				}
			}
		}
	}()

	return recordSource
}

func (fs *FrameScanner) shiftBufferLeft(n int) {
	copy(fs.buffer[:], fs.buffer[n:])
	fs.bufferStart -= n
}

// Err returns any errors returned from the Scanner or its underlaying reader.
// Err will be nil if the reader is closed and scanning completes successfully.
func (fs *FrameScanner) Err() error {
	return fs.err
}
