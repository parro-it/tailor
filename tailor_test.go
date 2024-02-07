package tailor

import (
	"os"
	"testing"
	"time"
)

func TestOpenFile(t *testing.T) {
	fpath := os.TempDir() + "/test.txt"
	os.Remove(fpath)
	defer os.Remove(fpath)

	err := os.WriteFile(fpath, []byte("123-456-789-"), 0644)
	if err != nil {
		t.Fatalf("WriteFile: %s", err)
	}

	tl, err := OpenFile(fpath, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("OpenFile: %s", err)
	}
	defer tl.Close()
	var buf = make([]byte, 10)
	n, err := tl.Read(buf)
	if err != nil {
		t.Fatalf("OpenFile: %s", err)
	}
	if n != 10 {
		t.Fatalf("OpenFile: expected 10, got %d", len(buf))
	}

	buf = make([]byte, 4)
	n, err = tl.Read(buf)
	if err != nil {
		t.Fatalf("OpenFile: %s", err)
	}
	if n != 2 {
		t.Fatalf("OpenFile: expected 2, got %d", n)
	}

	errs := make(chan error)
	go func() {
		time.Sleep(100 * time.Millisecond)
		err := os.WriteFile(fpath, []byte("123-456-789-ABCDEF"), 0644)
		if err != nil {
			errs <- err
			return
		}
		close(errs)
	}()

	buf = make([]byte, 6)
	n, err = tl.Read(buf)
	if err != nil {
		t.Fatalf("Read: %s", err)
	}

	if n != 6 {
		t.Fatalf("OpenFile: expected 6, got %d", n)
	}

	if string(buf) != "ABCDEF" {
		t.Fatalf("OpenFile: expected 'ABCDEF', got %s", string(buf))
	}

	if <-errs != nil {
		t.Fatalf("WriteFile in background: %s", err)
	}

}
