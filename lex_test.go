package atnodes

import (
	"testing"
)

func TestLex(t *testing.T) {
	l := NewLexer("www[101-110].example.com")
	go l.Run(l)

	for item := range l.Items() {
		t.Logf("%#v", item)
	}
}
