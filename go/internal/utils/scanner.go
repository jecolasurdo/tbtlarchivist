package utils

import (
	"errors"
	"io"
)

// Scanner is a copy of the bufio.Scanner with the following modifications:
// 1)  There is no default splitFunc. It must be supplied via NewScanner. Since
//     The splitFunc is explicitly supplied via NewScanner, the Split() method,
//     found in the bufio implementation, has been omitted here.
// 2a) There is no hard-coded limit to the number of calls that fail to advance
//     the reader.
// 2b) The Scanner relies on a Backoff structure which provides the behavior
//     for consecutive failures to advance. If the scanner fails to advance, the
//     Backoff.Wait() method is called. If Backoff.Wait returns an error, the
//     Scanner is halted and the error is returned via Scanner.Err(). If the
//     scanner succeeds in advancing the reader, the Backoff.Reset() method is
//     called.
// 3)  This Scanner elides the bufio.Scanner.Text() method.
type Scanner struct {
	r            io.Reader // The reader provided by the client.
	split        SplitFunc // The function to split the tokens.
	maxTokenSize int       // Maximum size of a token; modified by tests.
	token        []byte    // Last token returned by split.
	buf          []byte    // Buffer used as argument to split.
	start        int       // First non-processed byte in buf.
	end          int       // End of data in buf.
	err          error     // Sticky error.
	scanCalled   bool      // Scan has been called; buffer is in use.
	done         bool      // Scan has finished.
	backoff      *Backoff
}

// SplitFunc is identical to a bufio.SplitFunc
type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)

// Errors returned by Scanner.
var (
	ErrTooLong         = errors.New("bufio.Scanner: token too long")
	ErrNegativeAdvance = errors.New("bufio.Scanner: SplitFunc returns negative advance count")
	ErrAdvanceTooFar   = errors.New("bufio.Scanner: SplitFunc returns advance count beyond input")
	ErrBadReadCount    = errors.New("bufio.Scanner: Read returned impossible count")
)

const (
	// MaxScanTokenSize > See bufio.MaxScanTokenSize
	MaxScanTokenSize = 64 * 1024

	startBufSize = 4096 // Size of initial allocation for buffer.
)

// NewScanner returns a new Scanner to read from r.  splitFunc defines the
// behavior for token idenfification.  backoff defines behavior for repeated
// ineffective reads.
func NewScanner(r io.Reader, splitFunc SplitFunc, backoff *Backoff) *Scanner {
	if backoff == nil {
		panic("A Backoff object must be supplied.")
	}
	return &Scanner{
		r:            r,
		split:        splitFunc,
		maxTokenSize: MaxScanTokenSize,
		backoff:      backoff,
	}
}

// Err > See bufio.Scanner.Err()
func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

// Bytes > See bufio.Scanner.Bytes()
func (s *Scanner) Bytes() []byte {
	return s.token
}

// ErrFinalToken > See bufio.ErrFinalToken
var ErrFinalToken = errors.New("final token")

// Scan advances the Scanner to the next token, which will then be available
// through the Bytes method. It returns false when the scan stops, either by
// reaching the end of the input or an error.  After Scan returns false, the
// Err method will return any error that occurred during scanning, except that
// if it was io.EOF, Err will return nil.  The bufio.Scanner.Scan()
// implementation emits errors and/or panics if many consecutive reads occur
// without effectively advancing the scanner. This implementation defers that
// behavior to a Backoff structure that is supplied via the NewScanner method.
func (s *Scanner) Scan() bool {
	if s.done {
		return false
	}
	s.scanCalled = true
	// Loop until we have a token.
	for {
		// See if we can get a token with what we already have.
		// If we've run out of data but have an error, give the split function
		// a chance to recover any remaining, possibly empty token.
		if s.end > s.start || s.err != nil {
			advance, token, err := s.split(s.buf[s.start:s.end], s.err != nil)
			if err != nil {
				if err == ErrFinalToken {
					s.token = token
					s.done = true
					return true
				}
				s.setErr(err)
				return false
			}
			if !s.advance(advance) {
				return false
			}
			s.token = token
			if token != nil {
				if s.err == nil || advance > 0 {
					// We advanced the scanner without error, reset the backoff
					// state.
					s.backoff.Reset()
				} else {
					// We failed to advance the scanner. Apply a wait and see
					// if we should continue or not.
					err = s.backoff.Wait()
					if err != nil {
						s.setErr(err)
						return false
					}
				}
				return true
			}
		}
		// We cannot generate a token with what we are holding.
		// If we've already hit EOF or an I/O error, we are done.
		if s.err != nil {
			// Shut it down.
			s.start = 0
			s.end = 0
			return false
		}
		// Must read more data.
		// First, shift data to beginning of buffer if there's lots of empty space
		// or space is needed.
		if s.start > 0 && (s.end == len(s.buf) || s.start > len(s.buf)/2) {
			copy(s.buf, s.buf[s.start:s.end])
			s.end -= s.start
			s.start = 0
		}
		// Is the buffer full? If so, resize.
		if s.end == len(s.buf) {
			// Guarantee no overflow in the multiplication below.
			const maxInt = int(^uint(0) >> 1)
			if len(s.buf) >= s.maxTokenSize || len(s.buf) > maxInt/2 {
				s.setErr(ErrTooLong)
				return false
			}
			newSize := len(s.buf) * 2
			if newSize == 0 {
				newSize = startBufSize
			}
			if newSize > s.maxTokenSize {
				newSize = s.maxTokenSize
			}
			newBuf := make([]byte, newSize)
			copy(newBuf, s.buf[s.start:s.end])
			s.buf = newBuf
			s.end -= s.start
			s.start = 0
		}
		// Finally we can read some input.
		for {
			n, err := s.r.Read(s.buf[s.end:len(s.buf)])
			if n < 0 || len(s.buf)-s.end < n {
				s.setErr(ErrBadReadCount)
				break
			}
			s.end += n
			if err != nil {
				s.setErr(err)
				break
			}
			// We succeeded in reading, reset the backoff behavior.
			if n > 0 {
				s.backoff.Reset()
				break
			}
			// We did not accomplish a read. Call Wait and see if we should
			// continue.
			err = s.backoff.Wait()
			if err != nil {
				s.setErr(err)
				break
			}
		}
	}
}

// advance consumes n bytes of the buffer. It reports whether the advance was legal.
func (s *Scanner) advance(n int) bool {
	if n < 0 {
		s.setErr(ErrNegativeAdvance)
		return false
	}
	if n > s.end-s.start {
		s.setErr(ErrAdvanceTooFar)
		return false
	}
	s.start += n
	return true
}

// setErr records the first error encountered.
func (s *Scanner) setErr(err error) {
	if s.err == nil || s.err == io.EOF {
		s.err = err
	}
}

// Buffer > See bufio.Scanner.Buffer()
func (s *Scanner) Buffer(buf []byte, max int) {
	if s.scanCalled {
		panic("Buffer called after Scan")
	}
	s.buf = buf[0:cap(buf)]
	s.maxTokenSize = max
}
