package parser

import (
	"io"
	"strconv"
	"strings"
	"time"

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
			case "h3":
				// H3 is a folder
				folder := parseFolder(token)

				// Get the folder title from text content
				folder.Title = p.getTextContent()

				// Add to current parent folder
				if len(folderStack) == 0 {
					// Malformed HTML - no parent folder available
					folderStack = append(folderStack, root)
				}
				currentFolder := folderStack[len(folderStack)-1]
				currentFolder.AddChild(folder)

				// Push this folder onto the stack for its children
				folderStack = append(folderStack, folder)
			case "a":
				// A is a bookmark
				bookmark := parseBookmark(token)

				// Skip Firefox place: URLs (dynamic queries)
				if strings.HasPrefix(bookmark.URL, "place:") {
					// Skip text content to advance parser
					p.getTextContent()
					continue
				}

				// Get the bookmark title from text content
				bookmark.Title = p.getTextContent()

				// Skip bookmarks without URLs (malformed)
				if bookmark.URL == "" {
					continue
				}

				// Add to current parent folder
				if len(folderStack) == 0 {
					// Malformed HTML - no parent folder available
					folderStack = append(folderStack, root)
				}
				currentFolder := folderStack[len(folderStack)-1]
				currentFolder.AddChild(bookmark)
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

// parseFolder extracts folder information from an H3 token
func parseFolder(token html.Token) *models.Folder {
	folder := &models.Folder{}

	// Extract attributes
	for _, attr := range token.Attr {
		switch attr.Key {
		case "add_date":
			if ts, err := strconv.ParseInt(attr.Val, 10, 64); err == nil {
				folder.AddDate = time.Unix(ts, 0)
			}
		case "last_modified":
			if ts, err := strconv.ParseInt(attr.Val, 10, 64); err == nil {
				folder.LastModified = time.Unix(ts, 0)
			}
		}
	}

	return folder
}

// parseBookmark extracts bookmark information from an A token
func parseBookmark(token html.Token) *models.Bookmark {
	bookmark := &models.Bookmark{}

	// Extract attributes
	for _, attr := range token.Attr {
		switch attr.Key {
		case "href":
			bookmark.URL = attr.Val
		case "add_date":
			if ts, err := strconv.ParseInt(attr.Val, 10, 64); err == nil {
				bookmark.AddDate = time.Unix(ts, 0)
			}
		case "last_modified":
			if ts, err := strconv.ParseInt(attr.Val, 10, 64); err == nil {
				bookmark.LastModified = time.Unix(ts, 0)
			}
		case "tags":
			// Parse comma-separated tags
			if attr.Val != "" {
				tagList := strings.Split(attr.Val, ",")
				for _, tag := range tagList {
					trimmed := strings.TrimSpace(tag)
					if trimmed != "" {
						bookmark.Tags = append(bookmark.Tags, trimmed)
					}
				}
			}
		case "shortcuturl":
			bookmark.ShortcutURL = attr.Val
		case "icon", "icon_uri":
			// Skip icon data as per requirements
		}
	}

	return bookmark
}

// getTextContent reads text content until the closing tag
func (p *HTMLParser) getTextContent() string {
	var content string
	for {
		tt := p.Next()
		if tt == html.TextToken {
			content += p.Token().Data
		} else if tt == html.EndTagToken {
			break
		} else if tt == html.ErrorToken {
			break
		}
	}
	return content
}
