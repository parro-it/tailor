package tailor

import (
	"io"
	"io/fs"
	"os"
	"time"
)

type tailor struct {
	inner *os.File
	wait  time.Duration
}

func OpenFile(name string, wait time.Duration) (io.ReadCloser, error) {
	for {
		f, err := os.Open(name)
		if os.IsNotExist(err) {
			time.Sleep(wait)
			continue
		}
		if err != nil {
			return nil, err
		}
		return &tailor{f, wait}, nil
	}
}

// Close implements io.ReadCloser.
func (tl tailor) Close() error {
	return tl.inner.Close()
}

// Read implements io.Reader.
func (tl *tailor) Read(p []byte) (n int, err error) {
	n, err = tl.inner.Read(p)
	if err == io.EOF {
		if err = tl.waitSizeIncrease(); err != nil {
			return
		}

		return tl.Read(p[n:])
	}
	return
}

func (tl *tailor) waitSizeIncrease() (err error) {
	var st fs.FileInfo
	fname := tl.inner.Name()
	/*
		// 1) close the file
		tl.inner.Close()
		err = nil
	*/
	// 2) read the actual size of the file
	st, err = os.Stat(fname)
	if err != nil {
		return
	}

	origSz := st.Size()
	sz := origSz

	// 3) wait until the file size increase, or Close() is called
	for sz <= origSz {
		time.Sleep(tl.wait)

		st, err = os.Stat(fname)
		if err != nil {
			return
		}
		sz = st.Size()
	}
	/*
		// 4) reopen the file, seek to previous position
		tl.inner, err = os.Open(fname)
		if err != nil {
			return
		}
	*/
	_, err = tl.inner.Seek(origSz, 0)
	return
}
