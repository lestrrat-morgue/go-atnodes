package atnodes

import (
  "github.com/lestrrat/go-lex"
)

const (
  ItemText  lex.LexItemType = lex.ItemDefaultMax + 1 + iota
  ItemRange
  ItemOpenSquareBracket
  ItemCloseSquareBracket
)

type Lexer struct {
  *lex.StringLexer
}

func NewLexer(input string) *Lexer {
  l := &Lexer {
    lex.NewStringLexer(input, lexText),
  }
  return l
}

func lexText(l lex.Lexer, ctx interface {}) lex.LexFn {
  al := ctx.(*Lexer)

  for {
    if al.Peek() == '[' {
      if al.Cursor() > al.LastCursor() {
        al.Emit(ItemText)
      }
      return lexRange
    }
    if al.Next() == lex.EOF {
      break
    }
  }

  if al.Cursor() > al.LastCursor() {
    al.Emit(ItemText)
  }

  l.Emit(lex.ItemEOF)
  return nil
}

func lexRange(l lex.Lexer, ctx interface {}) lex.LexFn {
  al := ctx.(*Lexer)
  if c := al.Next(); c != '[' {
    return al.EmitErrorf("Expected '[', got %c", c)
  }
  al.Emit(ItemOpenSquareBracket)

  for {
    c := al.Next()
    switch c {
    case lex.EOF:
      return al.EmitErrorf("Expected range, got EOF")
    case ']':
      al.Backup()
      l.Emit(ItemRange)
      al.Next()
      l.Emit(ItemCloseSquareBracket)
      return lexText
    }
  }
}
