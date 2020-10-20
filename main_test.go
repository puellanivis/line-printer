package main

import (
	"bytes"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type WriterMock struct {
	mock.Mock
}

func (m WriterMock) Write(b []byte) (int, error) {
	args := m.Called(b)

	return args.Int(0), args.Error(1)
}

func TestLinePrinter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setup         func(w *WriterMock)
		input         string
		expectedError error
	}{
		{
			name: "test hello world",
			setup: func(w *WriterMock) {
				w.On("Write", []byte("hello world\n")).Return(len("hello world\n"), nil)
			},
			input:         "hello world",
			expectedError: nil,
		},
		{
			name: "test error",
			setup: func(w *WriterMock) {
				w.On("Write", []byte("error\n")).Return(0, errors.New("example error"))
			},
			input:         "error",
			expectedError: errors.New("example error"),
		},
		{
			name: "test negative write",
			setup: func(w *WriterMock) {
				w.On("Write", []byte("error\n")).Return(-42, nil)
			},
			input:         "error",
			expectedError: errors.New("example error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWriter := new(WriterMock)
			if tt.setup != nil {
				tt.setup(mockWriter)
			}

			lp := NewLinePrinter(mockWriter)

			err := lp.PrintLine(tt.input)
			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockWriter.AssertExpectations(t)
		})
	}
}

func TestMain(t *testing.T) {
	mockWriter := new(WriterMock)
	mockWriter.On("Write", []byte("Good morning!\n")).Return(len("Good morning!\n"), nil)

	save := output
	output = mockWriter
	defer func() {
		output = save
	}()

	main()

	mockWriter.AssertExpectations(t)
}
