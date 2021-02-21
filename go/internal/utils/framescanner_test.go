package utils_test

import (
	"bufio"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	mock_io "github.com/jecolasurdo/tbtlarchivist/go/internal/mocks/io"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/utils"
)

// IOReader is a placeholder for io.Reader used for generating mock readers.
// This interface should not be used directly. See mocks/io/mock_reader.
type IOReader interface {
	io.Reader
}

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
		name         string
		inboundBytes []byte
		expResult    [][]byte
	}{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			reader := mock_io.NewMockIOReader(ctrl)
			scanner := bufio.NewScanner(reader)
			frameScanner := new(utils.FrameScanner)
			scanner.Split(frameScanner.ScanFrames)
			actualResult := [][]byte{}
			for scanner.Scan() {
				actualResult = append(actualResult, scanner.Bytes())
			}
		})
	}
}
