package utils_test

import (
	"encoding/hex"
	"fmt"
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
		hexData    string
		atEOF      bool
		expAdvance int
		expToken   []byte
		expErr     error
	}{
		{
			name:       "no data",
			hexData:    "",
			atEOF:      false,
			expAdvance: 0,
			expToken:   nil,
			expErr:     nil,
		},
		{
			name:       "initial read insufficient",
			hexData:    "000E",
			atEOF:      false,
			expAdvance: 2,
			expToken:   nil,
			expErr:     nil,
		},
		{
			// Should read first 4 bytes (00 0E E7 C2) and ignore the remaining.
			name:       "initial read",
			hexData:    "000EE7C2FFFF",
			atEOF:      false,
			expAdvance: 976834,
			expToken:   nil,
			expErr:     nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			frameScanner := new(utils.FrameScanner)

			data := make([]byte, hex.DecodedLen(len(testCase.hexData)))
			_, err := hex.Decode(data, []byte(testCase.hexData))
			if err != nil {
				panic(fmt.Sprintf("malformed test case: %v\n%v", testCase.name, err))
			}

			advance, token, err := frameScanner.ScanFrames(data, testCase.atEOF)
			assert.Equal(t, testCase.expAdvance, advance, "incorrect advance value")
			assert.Equal(t, testCase.expToken, token, "incorrect token value")
			assert.Equal(t, testCase.expErr, err, "incorrect error value")
		})
	}
}
