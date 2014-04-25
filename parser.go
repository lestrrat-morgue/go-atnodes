package atnodes

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/lestrrat/go-lex"
	"sort"
	"strconv"
)

type Parser struct {
	lex.Consumer
}

func NewParser() *Parser {
	return &Parser{nil}
}

func (p *Parser) ParseString(template string) []string {
	l := NewLexer(template)
	go l.Run(l)

	p.Consumer = lex.NewItemConsume(l)

	s := mapset.NewSet()
	for item := p.Peek(); item.Type() != lex.ItemEOF; item = p.Peek() {
		switch item.Type() {
		case ItemText:
			p.parseTerm(&s)
		case ItemOpenSquareBracket:
			p.parsePermutation(&s)
		default:
			p.unexpected("Unexpected item found: %s", item)
		}
	}

	expanded := make([]string, s.Cardinality())
	i := 0
	for s := range s.Iter() {
		expanded[i] = s.(string)
		i++
	}

	sort.StringSlice(expanded).Sort()
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
	openBracket := p.Consume()
	if openBracket.Type() != ItemOpenSquareBracket {
		p.unexpected("Expected '[', got %s", openBracket)
	}

	newRange := mapset.NewSet()
	for {
		num := p.Consume()
		if num.Type() != ItemNumber {
			p.unexpected("Expected number, got %s", num)
		}

		switch p.Peek().Type() {
		case ItemMinus:
			// It's a range!
			p.Backup2(num)
			p.parseRange(&newRange)
		default:
			// It's a single element
			newRange.Add(num.Value())
		}

		// This may be followed by a comma, which would mean that there are
		// more rules to be applied
		if p.Peek().Type() != ItemComma {
			break
		}
		p.Consume() // consume the comma
	}

	closeBracket := p.Consume()
	if closeBracket.Type() != ItemCloseSquareBracket {
		p.unexpected("Expected ']', got %s", openBracket)
	}

	*s = p.permute(s, newRange)
}

func (p *Parser) parseRange(s *mapset.Set) {

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

	newSet := mapset.NewSet()
	for i := fromInt; i <= toInt; i++ {
		newSet.Add(strconv.FormatInt(int64(i), 10))
	}
	*s = s.Union(newSet)
}

func (p *Parser) permute(s *mapset.Set, newSet mapset.Set) mapset.Set {
	if s.Cardinality() <= 0 {
		// set s is empty, just return newSet
		return newSet
	}
	// otherwise, take all elements in set s, and create new set using
	// contents of newSet
	ret := mapset.NewSet()
	for prefix := range s.Iter() {
		for suffix := range newSet.Iter() {
			ret.Add(fmt.Sprintf("%s%s", prefix.(string), suffix.(string)))
		}
	}
	return ret
}
