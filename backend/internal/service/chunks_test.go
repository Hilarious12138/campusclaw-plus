package service

import (
	"strings"
	"testing"
)

func TestChunkStrategies(t *testing.T) {
	samples := []string{"短文😀。", strings.Repeat("中文😀ab", 500), strings.Repeat("长句。\n\n", 300), "   ", "# 一级\n正文\n## 二级\n" + strings.Repeat("乙", 1000) + "\n### 三级\n尾片"}
	for _, strategy := range []string{"auto", "custom", "hierarchy"} {
		for _, body := range samples {
			o, e := (ChunkOptions{Strategy: strategy, MaxLength: 100, OverlapRatio: .5, Delimiter: "newline"}).normalized()
			if e != nil {
				t.Fatal(e)
			}
			text, basis, chunks := splitText(body, o)
			if basis != "original" || text != body {
				t.Fatal("original altered")
			}
			prev := -1
			for _, c := range chunks {
				if c.Start <= prev || c.End <= c.Start || c.End-c.Start > o.MaxLength || string([]rune(text)[c.Start:c.End]) != c.Text {
					t.Fatalf("bad chunk %#v", c)
				}
				prev = c.Start
			}
			if strings.TrimSpace(body) != "" && len(chunks) == 0 {
				t.Fatal("empty")
			}
		}
	}
	o, _ := (ChunkOptions{Strategy: "auto", MaxLength: 1, RemoveURLs: true}).normalized()
	text, _, c := splitText(strings.Repeat("a", 900), o)
	if len(c) != 2 || c[0].End != 800 || c[1].Start != 720 || len(text) != 900 {
		t.Fatal("800/80")
	}
	o, _ = (ChunkOptions{Strategy: "custom", MaxLength: 100, RemoveURLs: true, RemoveEmails: true, CollapseWhitespace: true}).normalized()
	text, basis, c := splitText("参考 https://example.com\n\n联系 a@b.com  正文", o)
	if text != "参考 联系 正文" || basis != "processed" || c[0].Text != text {
		t.Fatal(text, basis)
	}
	for _, o := range []ChunkOptions{{Strategy: "custom", MaxLength: 99}, {Strategy: "custom", MaxLength: 2001}, {Strategy: "custom", MaxLength: 800, OverlapRatio: .51}, {Strategy: "custom", MaxLength: 800, Delimiter: "bad"}, {Strategy: "bad"}} {
		if _, e := o.normalized(); e == nil {
			t.Fatal("invalid accepted")
		}
	}
	o, _ = (ChunkOptions{Strategy: "hierarchy"}).normalized()
	_, _, c = splitText("# 一\n甲\n## 二\n乙\n### 三\n丙", o)
	if len(c) != 3 || c[0].Text != "# 一\n甲\n" || !strings.HasPrefix(c[1].Text, "## 二") {
		t.Fatal("headings")
	}
}
