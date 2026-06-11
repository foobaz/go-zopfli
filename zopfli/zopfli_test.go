package zopfli

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"testing"
)

func TestIssue2(t *testing.T) {
	tests := []struct {
		in []byte
	}{
		// Issue 2.
		{[]byte{}},
		{[]byte{0}},
		// Just test a bigger one too.
		{make([]byte, 8192)},
	}

	opt := DefaultOptions()
	//opt.Verbose = true
	//opt.VerboseMore = true

	buf := &bytes.Buffer{}
	for _, tt := range tests {
		err := GzipCompress(&opt, tt.in, buf)
		if err != nil {
			t.Fatal("fail to compress:", err)
		}

		zr, err := gzip.NewReader(buf)
		if err != nil {
			t.Fatal(err)
		}
		out, err := ioutil.ReadAll(zr)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(tt.in, out) {
			t.Fatal("wanted %#v, got %#v", tt.in, out)
		}
		buf.Reset()
		t.Logf("encode len %v ok", len(tt.in))
	}
}
