package utils_test

import (
	"testing"

	"github.com/jecolasurdo/tbtlarchivist/go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_FrameScanner(t *testing.T) {
	// incomplete frame header after first read
	// incomplete frame header after first read
	// incomplete record
	// malformed frame header
	// no records
	// single record
	// empty record
	// multiple records

	testCases := []struct {
		name       string
		data       []byte
		atEOF      bool
		expAdvance int
		expToken   []byte
		expErr     error
	}{
		{
			name:       "no data",
			data:       []byte{},
			atEOF:      false,
			expAdvance: 0,
			expToken:   nil,
			expErr:     nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			frameScanner := new(utils.FrameScanner)
			advance, token, err := frameScanner.ScanFrames(testCase.data, testCase.atEOF)
			assert.Equal(t, testCase.expAdvance, advance, "incorrect advance value")
			assert.Equal(t, testCase.expToken, token, "incorrect token value")
			assert.Equal(t, testCase.expErr, err, "incorrect error value")
		})
	}
}
