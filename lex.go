package atnodes

import (
  "github.com/lestrrat/go-lex"
)

const (
  ItemText  lex.LexItemType = lex.ItemDefaultMax + 1 + iota
  ItemNumber
  ItemMinus
  ItemComma
  ItemOpenSquareBracket
  ItemCloseSquareBracket
)

func init() {
  lex.TypeNames[ItemText] = "Text"
  lex.TypeNames[ItemMinus] = "Minus"
  lex.TypeNames[ItemComma] = "Comma"
  lex.TypeNames[ItemOpenSquareBracket] = "OpenSquareBracket"
  lex.TypeNames[ItemCloseSquareBracket] = "CloseSquareBracket"
}

type Lexer struct {
  *lex.StringLexer
}

func isNumeric(r rune) bool {
  return '0' <= r && r <= '9'
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
      return lexPermutation
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

func lexPermutation(l lex.Lexer, ctx interface {}) lex.LexFn {
  al := ctx.(*Lexer)
  if c := al.Next(); c != '[' {
    return al.EmitErrorf("Expected '[', got %c", c)
  }
  al.Emit(ItemOpenSquareBracket)

  for {
    c := al.Next()
    switch {
    case c == lex.EOF:
      return al.EmitErrorf("Expected range, got EOF")
    case isNumeric(c):
      if al.AcceptRun("0123456789") {
        l.Emit(ItemNumber)
      }
    case c == '-':
      l.Emit(ItemMinus)
    case c == ',':
      l.Emit(ItemComma)
    case c == ']':
      l.Emit(ItemCloseSquareBracket)
      return lexText
    }
  }
}
