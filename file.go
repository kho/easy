package easy

import (
	"compress/bzip2"
	"compress/gzip"
	"errors"
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
// underlying file, when name has suffix ".gz";
//
// - A bzip2 Reader wrapped around a Closer that closes the underlying
// file, when name has suffix ".bz2";
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
			return &alsoCloseReadCloser{r, f}, nil
		} else if strings.HasSuffix(name, ".bz2") {
			r := bzip2.NewReader(f)
			return &closeOtherReadCloser{r, f}, nil
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
// - Because there is no bzip2 compressor at this moment, creating a
// ".bz2" file results in an error;
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
			return &alsoCloseWriteCloser{w, f}, nil
		} else if strings.HasSuffix(name, ".bz2") {
			return nil, noBz2Error
		} else {
			return f, nil
		}
	}
}

var noBz2Error = errors.New("bz2 compression is not supported yet (check https://code.google.com/p/go/issues/detail?id=4828)")

type alsoCloseReadCloser struct {
	top    io.ReadCloser
	bottom io.Closer
}

func (a *alsoCloseReadCloser) Read(p []byte) (int, error) {
	return a.top.Read(p)
}

func (a *alsoCloseReadCloser) Close() error {
	if err := a.top.Close(); err != nil {
		return err
	}
	if err := a.bottom.Close(); err != nil {
		return err
	}
	return nil
}

type closeOtherReadCloser struct {
	top    io.Reader
	bottom io.Closer
}

func (c *closeOtherReadCloser) Read(p []byte) (int, error) {
	return c.top.Read(p)
}

func (c *closeOtherReadCloser) Close() error {
	return c.bottom.Close()
}

type alsoCloseWriteCloser struct {
	top    io.WriteCloser
	bottom io.Closer
}

func (a *alsoCloseWriteCloser) Write(p []byte) (int, error) {
	return a.top.Write(p)
}

func (a *alsoCloseWriteCloser) Close() error {
	if err := a.top.Close(); err != nil {
		return err
	}
	if err := a.bottom.Close(); err != nil {
		return err
	}
	return nil
}

// MustOpen tries to Open(name) and panics if it fails.
func MustOpen(name string) io.ReadCloser {
	r, err := Open(name)
	if err != nil {
		panic(err)
	}
	return r
}

// MustCreate tries to Create(name) and panics if it fails.
func MustCreate(name string) io.WriteCloser {
	w, err := Create(name)
	if err != nil {
		panic(err)
	}
	return w
}
