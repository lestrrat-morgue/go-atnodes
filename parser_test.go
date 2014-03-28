package atnodes

import (
  "testing"
)

func TestParser(t *testing.T) {
  p := NewParser()
  t.Logf("%#v", p.ParseString("www[101-110].example.com"))
}