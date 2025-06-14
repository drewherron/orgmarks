package converter

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/drewherron/orgmarks/models"
)

// ToHTML converts a bookmark tree to Netscape Bookmark HTML format
func ToHTML(root *models.Folder, w io.Writer) error {
	// Write HTML header
	if err := writeHTMLHeader(w); err != nil {
		return err
	}

	// Write the bookmark tree
	if err := writeHTMLNode(root, 0, w); err != nil {
		return err
	}

	return nil
}

// writeHTMLHeader writes the Netscape Bookmark format header
func writeHTMLHeader(w io.Writer) error {
	header := `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
`
	_, err := w.Write([]byte(header))
	return err
}

// writeHTMLNode recursively writes a node in HTML format
func writeHTMLNode(node models.Node, depth int, w io.Writer) error {
	indent := strings.Repeat("    ", depth)

	if node.IsFolder() {
		folder := node.(*models.Folder)

		// Skip root folder (depth 0), only write its children
		if depth > 0 {
			// Write folder header
			addDate := formatTimestamp(folder.AddDate)
			lastModified := formatTimestamp(folder.LastModified)

			_, err := fmt.Fprintf(w, "%s<DT><H3 ADD_DATE=\"%s\" LAST_MODIFIED=\"%s\">%s</H3>\n",
				indent, addDate, lastModified, escapeHTML(folder.Title))
			if err != nil {
				return err
			}

			// Start nested list
			_, err = fmt.Fprintf(w, "%s<DL><p>\n", indent)
			if err != nil {
				return err
			}
		}

		// Write children
		for _, child := range folder.Children {
			if err := writeHTMLNode(child, depth+1, w); err != nil {
				return err
			}
		}

		// Close nested list
		if depth > 0 {
			_, err := fmt.Fprintf(w, "%s</DL><p>\n", indent)
			if err != nil {
				return err
			}
		}
	} else {
		// Bookmark
		bookmark := node.(*models.Bookmark)

		// Build attributes
		attrs := []string{
			fmt.Sprintf("HREF=\"%s\"", escapeHTML(bookmark.URL)),
			fmt.Sprintf("ADD_DATE=\"%s\"", formatTimestamp(bookmark.AddDate)),
			fmt.Sprintf("LAST_MODIFIED=\"%s\"", formatTimestamp(bookmark.LastModified)),
		}

		// Add tags if present
		if len(bookmark.Tags) > 0 {
			attrs = append(attrs, fmt.Sprintf("TAGS=\"%s\"", strings.Join(bookmark.Tags, ",")))
		}

		// Add shortcut URL if present
		if bookmark.ShortcutURL != "" {
			attrs = append(attrs, fmt.Sprintf("SHORTCUTURL=\"%s\"", bookmark.ShortcutURL))
		}

		// Write bookmark
		_, err := fmt.Fprintf(w, "%s<DT><A %s>%s</A>\n",
			indent, strings.Join(attrs, " "), escapeHTML(bookmark.Title))
		if err != nil {
			return err
		}
	}

	// Close root DL tag at the end
	if depth == 0 {
		_, err := w.Write([]byte("</DL><p>\n"))
		return err
	}

	return nil
}

// formatTimestamp converts time.Time to Unix timestamp string
func formatTimestamp(t time.Time) string {
	if t.IsZero() {
		// Use current time if not set
		return fmt.Sprintf("%d", time.Now().Unix())
	}
	return fmt.Sprintf("%d", t.Unix())
}

// escapeHTML escapes special HTML characters
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
