package parser

import (
	"bufio"
	"io"
	"strings"

	"github.com/drewherron/orgmarks/models"
)

// OrgParser parses org-mode bookmark files
type OrgParser struct {
	scanner *bufio.Scanner
}

// NewOrgParser creates a new org-mode parser from a reader
func NewOrgParser(r io.Reader) *OrgParser {
	return &OrgParser{
		scanner: bufio.NewScanner(r),
	}
}

// headline represents a parsed org-mode headline
type headline struct {
	level int      // Number of * characters (1, 2, 3, etc.)
	title string   // The headline text
	tags  []string // Tags extracted from :tag1:tag2: format
}

// parseHeadline parses a line starting with * and returns headline info
func parseHeadline(line string) *headline {
	if !strings.HasPrefix(line, "*") {
		return nil
	}

	h := &headline{}

	// Count leading asterisks
	i := 0
	for i < len(line) && line[i] == '*' {
		i++
	}
	h.level = i

	// Skip whitespace after asterisks
	for i < len(line) && line[i] == ' ' {
		i++
	}

	// Rest of the line is title (possibly with tags)
	rest := strings.TrimSpace(line[i:])

	// Extract tags from end if present (format: :tag1:tag2:)
	if strings.Contains(rest, ":") {
		// Tags are at the end, separated by colons
		lastColon := strings.LastIndex(rest, ":")
		if lastColon > 0 {
			// Check if there's another colon before it
			firstColon := strings.Index(rest, ":")
			if firstColon < lastColon {
				// Extract tags
				tagString := rest[firstColon+1 : lastColon]
				if tagString != "" {
					h.tags = strings.Split(tagString, ":")
				}
				// Title is everything before the first tag colon (trimmed)
				h.title = strings.TrimSpace(rest[:firstColon])
			} else {
				// No tags, just title
				h.title = rest
			}
		} else {
			h.title = rest
		}
	} else {
		h.title = rest
	}

	return h
}

// parseProperty parses org-mode property lines like #+SHORTCUTURL: value
func parseProperty(line string) (key, value string, ok bool) {
	line = strings.TrimSpace(line)

	// Check if it starts with #+
	if !strings.HasPrefix(line, "#+") {
		return "", "", false
	}

	// Remove #+ prefix
	rest := strings.TrimSpace(line[2:])

	// Find the colon separator
	colonIdx := strings.Index(rest, ":")
	if colonIdx == -1 {
		return "", "", false
	}

	key = strings.ToUpper(strings.TrimSpace(rest[:colonIdx]))
	value = strings.TrimSpace(rest[colonIdx+1:])

	return key, value, true
}

// parseLink parses org-mode links like [[URL]] or [[URL][title]]
// Returns url, title (if present), and ok bool
func parseLink(line string) (url, title string, ok bool) {
	line = strings.TrimSpace(line)

	// Check if line contains [[
	if !strings.Contains(line, "[[") {
		return "", "", false
	}

	// Find the opening [[
	start := strings.Index(line, "[[")
	if start == -1 {
		return "", "", false
	}

	// Find the closing ]]
	end := strings.Index(line[start:], "]]")
	if end == -1 {
		return "", "", false
	}
	end += start // Make it absolute position

	// Extract content between [[ and ]]
	content := line[start+2 : end]

	// Check if it has the format [[URL][title]]
	if strings.Contains(content, "][") {
		// Split into URL and title
		parts := strings.Split(content, "][")
		if len(parts) == 2 {
			url = strings.TrimSpace(parts[0])
			title = strings.TrimSpace(parts[1])
			return url, title, true
		}
	}

	// Simple format [[URL]]
	url = strings.TrimSpace(content)
	return url, "", true
}

// isDescriptionLine returns true if the line is plain text (not a headline, property, or link)
func isDescriptionLine(line string) bool {
	trimmed := strings.TrimSpace(line)

	// Empty lines are not description
	if trimmed == "" {
		return false
	}

	// Headlines start with *
	if strings.HasPrefix(trimmed, "*") {
		return false
	}

	// Properties start with #+
	if strings.HasPrefix(trimmed, "#+") {
		return false
	}

	// Links are [[...]]
	if strings.HasPrefix(trimmed, "[[") {
		return false
	}

	// Everything else is description text
	return true
}

// hasLink checks if any line contains an org link
func hasLink(line string) bool {
	return strings.Contains(line, "[[") && strings.Contains(line, "]]")
}

// Note: Folders vs Bookmarks distinction:
// - A headline WITH a link (anywhere in its content section) = Bookmark
// - A headline WITHOUT a link = Folder
// This will be determined during the tree building phase (step 4.9)
// by looking ahead at the content lines after each headline

// Parse reads an org-mode bookmark file and returns the root folder
func (p *OrgParser) Parse() (*models.Folder, error) {
	root := &models.Folder{
		Title: "Bookmarks",
	}

	// Stack to track current folder hierarchy by level
	folderStack := []*models.Folder{root}
	levelStack := []int{0} // Root is level 0

	var currentHeadline *headline
	var contentLines []string

	for p.scanner.Scan() {
		line := p.scanner.Text()

		// Check if this is a new headline
		if h := parseHeadline(line); h != nil {
			// Process previous headline if exists
			if currentHeadline != nil {
				p.processHeadline(currentHeadline, contentLines, &folderStack, &levelStack)
			}

			// Start new headline
			currentHeadline = h
			contentLines = []string{}
		} else {
			// Accumulate content lines for current headline
			contentLines = append(contentLines, line)
		}
	}

	// Process final headline
	if currentHeadline != nil {
		p.processHeadline(currentHeadline, contentLines, &folderStack, &levelStack)
	}

	return root, p.scanner.Err()
}

// processHeadline processes a headline and its content, creating either a folder or bookmark
func (p *OrgParser) processHeadline(h *headline, contentLines []string, folderStack *[]*models.Folder, levelStack *[]int) {
	// Check if content has a link (determines if it's a bookmark or folder)
	var linkURL string
	var shortcutURL string
	var description strings.Builder

	for _, line := range contentLines {
		// Check for link
		if url, _, ok := parseLink(line); ok && linkURL == "" {
			linkURL = url
		}

		// Check for properties
		if key, value, ok := parseProperty(line); ok {
			if key == "SHORTCUTURL" {
				shortcutURL = value
			}
		}

		// Collect description text
		if isDescriptionLine(line) {
			if description.Len() > 0 {
				description.WriteString("\n")
			}
			description.WriteString(strings.TrimSpace(line))
		}
	}

	// Determine parent folder based on level
	// Pop stack until we find the correct parent level
	for len(*levelStack) > 1 && (*levelStack)[len(*levelStack)-1] >= h.level {
		*folderStack = (*folderStack)[:len(*folderStack)-1]
		*levelStack = (*levelStack)[:len(*levelStack)-1]
	}

	// Safety check: ensure stacks are not empty (should never happen, but guard against it)
	if len(*folderStack) == 0 {
		// Malformed structure - reset to root
		root := &models.Folder{Title: "Bookmarks"}
		*folderStack = []*models.Folder{root}
		*levelStack = []int{0}
	}

	parent := (*folderStack)[len(*folderStack)-1]

	if linkURL != "" {
		// This is a bookmark
		bookmark := &models.Bookmark{
			Title:       h.title,
			URL:         linkURL,
			Tags:        h.tags,
			ShortcutURL: shortcutURL,
			Description: description.String(),
		}

		// Skip bookmarks with empty titles (malformed)
		if bookmark.Title == "" {
			return
		}

		parent.AddChild(bookmark)
	} else {
		// This is a folder
		folder := &models.Folder{
			Title: h.title,
		}

		// Skip folders with empty titles (malformed)
		if folder.Title == "" {
			return
		}

		parent.AddChild(folder)

		// Push onto stack for children
		*folderStack = append(*folderStack, folder)
		*levelStack = append(*levelStack, h.level)
	}
}
