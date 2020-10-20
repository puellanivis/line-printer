package main

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"

	"github.com/stretchr/testify/mock"
)

// LinePrinter is an interface that defines how to print a line
type LinePrinter interface {
	PrintLine(string) error
}

// LinePrinterImpl implements LinePrinter
type LinePrinterImpl struct {
	writer io.Writer
}

var _ LinePrinter = LinePrinterImpl{}

// NewLinePrinter returns a LinePrinterImpl
func NewLinePrinter(w io.Writer) *LinePrinterImpl {
	lp := LinePrinterImpl{
		writer: w,
	}

	return &lp
}

// PrintLine implements LinePrinter
func (lp LinePrinterImpl) PrintLine(s string) error {
	if !strings.HasSuffix(s, "\n") {
		s = s + "\n"
	}

	n, err := lp.writer.Write([]byte(s))
	if n < 0 {
		return errors.New("negative write value")
	}

	if err != nil {
		return err
	}

	return nil
}

// LinePrinterMock implements LinePrinter
type LinePrinterMock struct {
	mock.Mock
}

var _ LinePrinter = LinePrinterMock{}

// PrintLine implements LinePrinter
func (m LinePrinterMock) PrintLine(s string) error {
	args := m.Called(s)

	return args.Error(0)
}

// LinePrinterFunc is a function that implements LinePrinter
type LinePrinterFunc func(s string) error

// PrintLine implements LinePrinter
func (f LinePrinterFunc) PrintLine(s string) error {
	return f(s)
}

var _ LinePrinter = LinePrinterFunc(nil)

// LinePrinterClosure returns a function that prints a line and logs the error if any.
func LinePrinterClosure(lp LinePrinter, log *log.Logger) LinePrinterFunc {
	return func(s string) error {
		err := lp.PrintLine(s)
		if err != nil {
			log.Printf("error occurred printing %s", s)
		}
		return err
	}
}

func main() {
	lp := NewLinePrinter(os.Stdout)
	logger := log.New(os.Stdout, "", log.LstdFlags)

	printLine := LinePrinterClosure(lp, logger)

	printLine.PrintLine("Good morning!")
}
