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

	type readCase struct {
		hexData    string
		atEOF      bool
		expAdvance int
		expToken   []byte
		expErr     error
	}

	testCases := []struct {
		name  string
		reads []readCase
	}{
		{
			name: "no data",
			reads: []readCase{
				{
					hexData:    "",
					atEOF:      false,
					expAdvance: 0,
					expToken:   nil,
					expErr:     nil,
				},
			},
		},
		{
			name: "initial read insufficient",
			reads: []readCase{
				{

					hexData:    "000E",
					atEOF:      false,
					expAdvance: 0,
					expToken:   nil,
					expErr:     nil,
				},
			},
		},
		{
			// Should read first 4 bytes (00 0E E7 C2) and ignore the remaining.
			name: "initial read",
			reads: []readCase{
				{
					hexData:    "000EE7C2FFFF",
					atEOF:      false,
					expAdvance: 976834,
					expToken:   nil,
					expErr:     nil,
				},
			},
		},
		{
			name: "initial read",
			reads: []readCase{
				{
					hexData:    "7C2FFF",
					atEOF:      false,
					expAdvance: 0,
					expToken:   nil,
					expErr:     nil,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			frameScanner := new(utils.FrameScanner)
			for _, readCase := range testCase.reads {
				data := make([]byte, hex.DecodedLen(len(readCase.hexData)))
				_, err := hex.Decode(data, []byte(readCase.hexData))
				if err != nil {
					panic(fmt.Sprintf("malformed test case: %v\n%v", testCase.name, err))
				}

				advance, token, err := frameScanner.ScanFrames(data, readCase.atEOF)
				assert.Equal(t, readCase.expAdvance, advance, "incorrect advance value")
				assert.Equal(t, readCase.expToken, token, "incorrect token value")
				assert.Equal(t, readCase.expErr, err, "incorrect error value")
			}
		})
	}
}
