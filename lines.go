package easy

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

type LineError struct {
	Num  int    // 1-based line number.
	Line string // the actual line.
	Err  error  // the error.
}

func (e *LineError) Error() string {
	return fmt.Sprintf("line %d: %v: %q", e.Num, e.Err, e.Line)
}

func ForEachLine(r io.Reader, f func(string) error) error {
	s := bufio.NewScanner(r)
	n := 0
	for s.Scan() {
		n++
		if err := f(s.Text()); err != nil {
			return &LineError{n, s.Text(), err}
		}
	}
	if err := s.Err(); err != nil {
		return &LineError{-1, "", err}
	} else {
		return nil
	}
}

func ForEachByteLine(r io.Reader, f func([]byte) error) error {
	s := bufio.NewScanner(r)
	n := 0
	for s.Scan() {
		n++
		if err := f(s.Bytes()); err != nil {
			return &LineError{n, s.Text(), err}
		}
	}
	if err := s.Err(); err != nil {
		return &LineError{-1, "", err}
	} else {
		return nil
	}
}

func ForEachLineN(r io.Reader, every int, f func(string) error) error {
	n := 0
	return ForEachLine(r, func(x string) error {
		err := f(x)
		n++
		if n%every == 0 {
			log.Print(n)
		}
		return err
	})
}

func ForEachByteLineN(r io.Reader, every int, f func([]byte) error) error {
	n := 0
	return ForEachByteLine(r, func(x []byte) error {
		err := f(x)
		n++
		if n%every == 0 {
			log.Print(n)
		}
		return err
	})
}
