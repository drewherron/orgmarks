package parser

import (
	"io"

	"golang.org/x/net/html"
)

// HTMLParser wraps the golang.org/x/net/html tokenizer
type HTMLParser struct {
	tokenizer *html.Tokenizer
}

// NewHTMLParser creates a new HTML parser from a reader
func NewHTMLParser(r io.Reader) *HTMLParser {
	return &HTMLParser{
		tokenizer: html.NewTokenizer(r),
	}
}

// Next advances to the next token and returns the token type
func (p *HTMLParser) Next() html.TokenType {
	return p.tokenizer.Next()
}

// Token returns the current token
func (p *HTMLParser) Token() html.Token {
	return p.tokenizer.Token()
}

// Err returns any error encountered during parsing
func (p *HTMLParser) Err() error {
	return p.tokenizer.Err()
}
