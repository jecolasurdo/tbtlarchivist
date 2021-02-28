package utils

import (
	"encoding/binary"
	"io"
)

const (
	maxBufferSize   = 1000 * 1024
	frameHeaderSize = 4
)

type frameState int

const (
	frameStateReadingHeader frameState = iota
	frameStateReadingBody
)

// FrameScanner provides methods for scanning a reader that returns records
// that are prefixed with a big endian int32 that describes the length of the
// message to follow.
type FrameScanner struct {
	reader            io.Reader
	backoff           BackoffAPI
	err               error
	state             frameState
	buffer            [maxBufferSize]byte
	bufferStart       int
	bufferEnd         int
	currentRecordSize int
	atEOF             bool
}

// NewFrameScanner instantiates a new FrameScanner that can poll the supplied
// reader. Backoff is used to apply pacing, context cancellation, and bailout
// behavior in case the reader is misbehaving. For more details, see Poll().
func NewFrameScanner(reader io.Reader, backoff BackoffAPI) *FrameScanner {
	return &FrameScanner{
		buffer:      [maxBufferSize]byte{},
		reader:      reader,
		backoff:     backoff,
		state:       frameStateReadingHeader,
		bufferStart: 0,
		bufferEnd:   maxBufferSize,
		atEOF:       false,
	}
}

// Poll begins polling the underlaying reader, and sends complete records to a
// returned channel. This continues until the underlaying reader returns an
// error, closes (returns io.EOF), or until the underlaying Backoff object's
// Wait method returns an error. In general, any time the reader is able to
// successfully read a frame header or record, it will call Backoff.Reset().
// However each read attempt that fails to read a full frame header or full
// record will call Backoff.Wait(), and evaluate whether or not Wait has
// returned an error. If wait returns an error, Poll will halt immediately, and
// the error can be evaluated via the Err method.
func (fs *FrameScanner) Poll() <-chan []byte {
	recordSource := make(chan []byte)
	go func() {
		defer close(recordSource)
		for {
			if !fs.atEOF {
				n, err := fs.reader.Read(fs.buffer[fs.bufferStart:fs.bufferEnd])
				if err != nil {
					if err == io.EOF {
						fs.atEOF = true
					} else {
						fs.err = err
						return
					}
				}
				fs.bufferStart += n
			} else if fs.bufferStart == 0 {
				// In normal operation, bufferStart will return to zero once
				// all records have been moved out of the buffer and to the
				// recordSource channel. If the buffer is empty (bufferStart is
				// zero) and the reader is closed (the reader returned EOF),
				// then we can exit cleanly. However, if the reader is
				// misbehaving and fails to send a complete final record, then
				// bufferStart will not return to zero. In that case
				// backoff.Wait() will be called either indefinitely or until
				// the backoff object decides enough is enough and returns an
				// error.
				return
			}
			switch fs.state {
			case frameStateReadingHeader:
				if fs.bufferStart >= frameHeaderSize {
					fs.currentRecordSize = int(binary.BigEndian.Uint32(fs.buffer[0:frameHeaderSize]))
					fs.shiftBufferLeft(frameHeaderSize)
					fs.state = frameStateReadingBody
					fs.backoff.Reset()
				} else {
					err := fs.backoff.Wait()
					if err != nil {
						fs.err = err
						return
					}
				}
			case frameStateReadingBody:
				if fs.bufferStart >= fs.currentRecordSize {
					record := make([]byte, fs.currentRecordSize)
					copy(record, fs.buffer[0:fs.currentRecordSize])
					recordSource <- record
					fs.shiftBufferLeft(fs.currentRecordSize)
					fs.state = frameStateReadingHeader
					fs.backoff.Reset()
				} else {
					err := fs.backoff.Wait()
					if err != nil {
						fs.err = err
						return
					}
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
