package parser

import (
	"io"

	"github.com/drewherron/orgmarks/models"
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

// Parse reads the HTML bookmark file and returns the root folder
func (p *HTMLParser) Parse() (*models.Folder, error) {
	root := &models.Folder{
		Title: "Bookmarks",
	}

	// Stack to track folder nesting when we encounter DL tags
	folderStack := []*models.Folder{root}

	for {
		tt := p.Next()
		if tt == html.ErrorToken {
			if p.Err() == io.EOF {
				break
			}
			return nil, p.Err()
		}

		token := p.Token()

		switch tt {
		case html.StartTagToken:
			switch token.Data {
			case "dl":
				// DL starts a new list level - no action needed yet
				// The current folder is already on the stack
			case "dt":
				// DT is a list item - could be folder (H3) or bookmark (A)
				// We'll handle these in their respective cases
			}
		case html.EndTagToken:
			switch token.Data {
			case "dl":
				// End of a list level - pop the folder stack
				if len(folderStack) > 1 {
					folderStack = folderStack[:len(folderStack)-1]
				}
			}
		}
	}

	return root, nil
}
