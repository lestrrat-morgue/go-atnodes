package atnodes

import (
  "fmt"
  "github.com/deckarep/golang-set"
  "github.com/lestrrat/go-lex"
  "strconv"
)

type Parser struct {
  lex.Consumer
}

func NewParser() *Parser {
  return &Parser { nil }
}

func (p *Parser) ParseString(template string) []string {
  l := NewLexer(template)
  go l.Run(l)

  p.Consumer = lex.NewItemConsume(l)

  s := mapset.NewSet()
  for item := p.Peek(); item.Type() != lex.ItemEOF; item = p.Peek() {
fmt.Printf("%#v\n", item)
    switch item.Type() {
    case ItemText:
      p.parseTerm(&s)
    case ItemOpenSquareBracket:
      p.parsePermutation(&s)
    default:
      p.unexpected("Unexpected item found: %s", item)
    }
    fmt.Printf("%v\n", s)
  }

  expanded := make([]string, s.Cardinality())
  i := 0
  for s := range s.Iter() {
    expanded[i] = s.(string)
    i++
  }
  return expanded
}

func (p *Parser) unexpected(format string, args ...interface{}) {
  panic(
    fmt.Sprintf(
      "Unexpected token found: %s",
      fmt.Sprintf(format, args...),
    ),
  )
}

func (p *Parser) parseTerm(s *mapset.Set) {
fmt.Println("parseTerm")
  term := p.Consume()
  if term.Type() != ItemText {
    p.unexpected("Expected text, got %s", term)
  }

  newSet := mapset.NewSet()
  if s.Cardinality() <= 0 {
    newSet.Add(term.Value())
  } else {
    for prefix := range s.Iter() {
      newSet.Add(prefix.(string) + term.Value())
    }
  }

  *s = newSet
}

func (p *Parser) parsePermutation(s *mapset.Set) {
fmt.Println("parsePermutation")
  openBracket := p.Consume()
  if openBracket.Type() != ItemOpenSquareBracket {
    p.unexpected("Expected '[', got %s", openBracket)
  }

  num := p.Consume()
  if num.Type() != ItemNumber {
    p.unexpected("Expected number, got %s", num)
  }

  if p.Peek().Type() == ItemMinus {
    // It's a range!
    p.Backup2(num)
    p.parseRange(s)
  }

  closeBracket := p.Consume()
  if closeBracket.Type() != ItemCloseSquareBracket {
    p.unexpected("Expected ']', got %s", openBracket)
  }
}

func (p *Parser) parseRange(s *mapset.Set) {
fmt.Println("parseRange")
  from := p.Consume()
  if from.Type() != ItemNumber {
    p.unexpected("Expected number, got %s", from)
  }
  fromInt, err := strconv.ParseInt(from.Value(), 10, 64)
  if err != nil {
    p.unexpected("Expected int, parse error: %s", err)
  }

  if hyphen := p.Consume(); hyphen.Type() != ItemMinus {
    p.unexpected("Expected '-', got %s", hyphen)
  }
  to := p.Consume()
  if to.Type() != ItemNumber {
    p.unexpected("Expected number, got %s", to)
  }
  toInt, err := strconv.ParseInt(to.Value(), 10, 64)
  if err != nil {
    p.unexpected("Expected int, parse error: %s", err)
  }
fmt.Printf("Parse range %d -> %d\n", fromInt, toInt)

  newSet := mapset.NewSet()
  if s.Cardinality() <= 0 {
    for i := fromInt; i <= toInt; i++ {
      newSet.Add(i)
    }
  } else {
    for prefix := range s.Iter() {
      for i := fromInt; i <= toInt; i++ {
        newSet.Add(fmt.Sprintf("%s%d", prefix.(string), i))
      }
    }
  }
  *s = newSet
}