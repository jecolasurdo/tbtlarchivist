package utils

import "io"

// FrameScanner provides methods for scanning a reader that returns records
// that are prefixed with a big endian int32 that descrives the length of the
// message to follow.
type FrameScanner struct {
	reader  io.Reader
	backoff *Backoff
}

// NewFrameScanner inststantiates a new FrameScanner.
func NewFrameScanner(reader io.ReadCloser, backoff *Backoff) *FrameScanner {
	return &FrameScanner{
		reader:  reader,
		backoff: backoff,
	}
}

// Scan attempts to advance the scanner to the next record. Scan will continue
// to return records until the reader is closed or returns an error.
func (fs *FrameScanner) Scan() bool {
	panic("not implemented")
}

// Bytes returns the current record or nil if there is no record available.
// Bytes should only be called after a call to Scan. Repeated calls to Bytes
// without a call to Scan will return the same record.
func (fs *FrameScanner) Bytes() []byte {
	panic("not implemented")
}

// Err returns any errors returned from the Scanner or its underlaying reader.
// Err will be nil if the reader is closed and scanning completes successfully.
func (fs *FrameScanner) Err() error {
	panic("not implemented")
}
