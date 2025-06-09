package parser

import (
	"bufio"
	"io"
	"strings"
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
