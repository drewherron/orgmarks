// Package converter provides functions for converting between the internal
// bookmark tree representation and external formats (HTML and Org-mode).
//
// ToOrg converts the internal model to Org-mode format.
// ToHTML converts the internal model to Netscape Bookmark HTML format.
package converter

import (
	"fmt"
	"io"
	"strings"

	"github.com/drewherron/orgmarks/models"
)

// ToOrg converts a bookmark tree to org-mode format
func ToOrg(root *models.Folder, w io.Writer) error {
	// Walk the tree and write org-mode format
	err := writeOrgNode(root, 0, w)
	return err
}

// writeOrgNode recursively writes a node in org-mode format
func writeOrgNode(node models.Node, depth int, w io.Writer) error {
	if node.IsFolder() {
		folder := node.(*models.Folder)

		// Skip root folder (depth 0), only write its children
		if depth > 0 {
			// Write folder headline
			stars := strings.Repeat("*", depth)
			if _, err := fmt.Fprintf(w, "%s %s\n", stars, folder.Title); err != nil {
				return err
			}
		}

		// Write children (handles empty folders gracefully - just writes headline)
		for _, child := range folder.Children {
			if err := writeOrgNode(child, depth+1, w); err != nil {
				return err
			}
		}

		// Add blank line after empty folders for readability
		if depth > 0 && len(folder.Children) == 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
	} else {
		// Bookmark
		bookmark := node.(*models.Bookmark)

		// Write bookmark headline with title and tags
		stars := strings.Repeat("*", depth)
		headline := fmt.Sprintf("%s %s", stars, bookmark.Title)

		// Add tags if present
		if len(bookmark.Tags) > 0 {
			// Pad to align tags (approximate 80 columns)
			padding := 80 - len(headline) - len(strings.Join(bookmark.Tags, ":")) - 2
			if padding < 1 {
				padding = 1
			}
			headline += strings.Repeat(" ", padding)
			headline += ":" + strings.Join(bookmark.Tags, ":") + ":"
		}

		if _, err := fmt.Fprintln(w, headline); err != nil {
			return err
		}

		// Write SHORTCUTURL property if present
		if bookmark.ShortcutURL != "" {
			if _, err := fmt.Fprintf(w, "#+SHORTCUTURL: %s\n", bookmark.ShortcutURL); err != nil {
				return err
			}
		}

		// Write link
		if _, err := fmt.Fprintf(w, "[[%s]]\n", bookmark.URL); err != nil {
			return err
		}

		// Write description if present
		if bookmark.Description != "" {
			if _, err := fmt.Fprintln(w, bookmark.Description); err != nil {
				return err
			}
		}

		// Empty line after bookmark for readability
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}

	return nil
}
