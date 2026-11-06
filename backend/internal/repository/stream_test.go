package repository

import (
	"io"
	"strings"
	"testing"
)

type byteReader struct{ io.Reader }

func (b byteReader) Read(p []byte) (int, error) { return b.Reader.Read(p[:1]) }
func TestModelSSE(t *testing.T) {
	s := ": heart\r\nevent: message\r\ndata: {\"choices\":[{\"delta\":{\"content\":\"中文😀\"}}]}\r\n\r\ndata: {\"choices\":\n" + "data: [{\"delta\":{\"content\":\"[1]\"},\"finish_reason\":\"stop\"}]}\n\ndata: [DONE]\n\n"
	var out strings.Builder
	n := 0
	e := parseModelStream(byteReader{strings.NewReader(s)}, func(text string) error { n++; out.WriteString(text); return nil })
	if e != nil || n != 2 || out.String() != "中文😀[1]" {
		t.Fatal(out.String(), e)
	}
	for _, bad := range []string{"data: broken\n\n", "data: {\"choices\":[{\"delta\":{\"content\":\"partial\"}}]}\n\n", "data: {\"choices\":[{\"delta\":{},\"finish_reason\":\"length\"}]}\n\n"} {
		if parseModelStream(strings.NewReader(bad), func(string) error { return nil }) == nil {
			t.Fatal("bad stream accepted")
		}
	}
}
