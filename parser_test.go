package atnodes

import (
	"testing"
)

func TestParser(t *testing.T) {
	p := NewParser()
	t.Logf("%#v", p.ParseString("www[101-110].example.com"))
	t.Logf("%#v", p.ParseString("www[101-110,111].example.com"))
	t.Logf("%#v", p.ParseString("www[101-110,200-210].example.com"))
	t.Logf("%#v", p.ParseString("www[101,200-210].example.com"))
}
