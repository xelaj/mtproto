package tlparser

import (
	"io"
	"unicode"
)

type Cursor struct {
	source []rune
	pos    int
}

// TODO: add `Line() int` and `Col() int` methods
func NewCursor(source string) *Cursor {
	return &Cursor{
		source: []rune(source),
	}
}

func (p *Cursor) current() rune {
	return p.source[p.pos]
}

func (p *Cursor) next() error {
	if p.pos >= len(p.source)-1 {
		return io.EOF
	}

	p.pos++
	return nil
}

func (p *Cursor) Unread(count int) {
	p.pos -= count
	if p.pos <= 0 {
		p.pos = 0
	}
}

func (p *Cursor) Skip(count int) {
	p.pos += count
	if p.pos >= len(p.source)-1 {
		p.pos = len(p.source) - 1
	}
}

func (p *Cursor) SkipSpaces() {
	for unicode.IsSpace(p.current()) {
		if err := p.next(); err != nil { // EOF
			break
		}
	}
}

func (p *Cursor) ReadAt(at rune) (string, error) {
	var str []rune
	for {
		r := p.current()
		if r == at {
			return string(str), nil
		}

		str = append(str, r)
		if err := p.next(); err != nil {
			return "", err
		}
	}
}

func (p *Cursor) ReadSymbol() (rune, error) {
	r := p.current()
	err := p.next()
	return r, err
}

func (p *Cursor) ReadDigits() (string, error) {
	var digits []rune
	for unicode.IsDigit(p.current()) {
		digits = append(digits, p.current())
		if err := p.next(); err != nil {
			return "", err
		}
	}

	return string(digits), nil
}

func (p *Cursor) IsNext(s string) bool {
	for i, exp := range s {
		if p.current() != exp {
			p.Unread(i)
			return false
		}
		_ = p.next()
	}

	return true
}
