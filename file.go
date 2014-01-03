package easy

import (
	"compress/gzip"
	"io"
	"os"
	"strings"
)

// Open opens a file for transparent sequential reading. The returned
// object can be read and closed much like os.Open. Based on name, it
// can be,
//
// - os.Stdin, when name is "-";
//
// - A *gzip.Reader wrapped around a Closer that closes both it and its
// underlying file, when name has prefix ".gz";
//
// - A normal *os.File otherwise.
func Open(name string) (io.ReadCloser, error) {
	if name == "-" {
		return os.Stdin, nil
	} else {
		f, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(name, ".gz") {
			r, err := gzip.NewReader(f)
			if err != nil {
				return nil, err
			}
			return alsoCloseReadCloser{r, f}, nil
		} else {
			return f, nil
		}
	}
}

// Create opens a file for transparent sequential writing. The
// returned object can be written and closed much like
// os.Create. Based on name, it can be,
//
// - os.Stdout, when name is "-";
//
// - A *gzip.Writer wrapped around a Closer that closes both it and
// its underlying file, when name has prefix ".gz";
//
// - A normal *os.File otherwise.
func Create(name string) (io.WriteCloser, error) {
	if name == "-" {
		return os.Stdout, nil
	} else {
		f, err := os.Create(name)
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(name, ".gz") {
			w := gzip.NewWriter(f)
			return alsoCloseWriteCloser{w, f}, nil
		} else {
			return f, nil
		}
	}
}

type alsoCloseReadCloser struct {
	top    io.ReadCloser
	bottom io.Closer
}

func (a alsoCloseReadCloser) Read(p []byte) (int, error) {
	return a.top.Read(p)
}

func (a alsoCloseReadCloser) Close() error {
	if err := a.top.Close(); err != nil {
		return err
	}
	if err := a.bottom.Close(); err != nil {
		return err
	}
	return nil
}

type alsoCloseWriteCloser struct {
	top    io.WriteCloser
	bottom io.Closer
}

func (a alsoCloseWriteCloser) Write(p []byte) (int, error) {
	return a.top.Write(p)
}

func (a alsoCloseWriteCloser) Close() error {
	if err := a.top.Close(); err != nil {
		return err
	}
	if err := a.bottom.Close(); err != nil {
		return err
	}
	return nil
}
