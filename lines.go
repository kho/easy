package easy

import (
	"bufio"
	"io"
	"log"
)

func ForEachLine(r io.Reader, f func(string) error) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		if err := f(s.Text()); err != nil {
			return err
		}
	}
	return s.Err()
}

func ForEachLineBytes(r io.Reader, f func([]byte) error) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		if err := f(s.Bytes()); err != nil {
			return err
		}
	}
	return s.Err()
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

func ForEachLineBytesN(r io.Reader, every int, f func([]byte) error) error {
	n := 0
	return ForEachLineBytes(r, func(x []byte) error {
		err := f(x)
		n++
		if n%every == 0 {
			log.Print(n)
		}
		return err
	})
}
