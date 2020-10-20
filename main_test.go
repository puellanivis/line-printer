package main

import (
	"bytes"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinePrinterClosure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setup         func(lp *LinePrinterMock)
		input         string
		expectedError error
	}{
		{
			name: "test hello world",
			setup: func(lp *LinePrinterMock) {
				lp.On("PrintLine", "hello world").Return(nil)
			},
			input:         "hello world",
			expectedError: nil,
		},
		{
			name: "test error",
			setup: func(lp *LinePrinterMock) {
				lp.On("PrintLine", "error").Return(errors.New("example error"))
			},
			input:         "error",
			expectedError: errors.New("example error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var logBuffer bytes.Buffer

			logger := log.New(&logBuffer, "", 0)

			mockLinePrinter := new(LinePrinterMock)
			if tt.setup != nil {
				tt.setup(mockLinePrinter)
			}

			lpc := LinePrinterClosure(mockLinePrinter, logger)

			err := lpc.PrintLine(tt.input)
			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError)
				assert.Equal(t, logBuffer.String(), "error occurred printing error\n")
			} else {
				assert.NoError(t, err)
			}

			mockLinePrinter.AssertExpectations(t)
		})
	}
}
