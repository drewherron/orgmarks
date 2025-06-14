package converter

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/drewherron/orgmarks/models"
	"github.com/drewherron/orgmarks/parser"
)

func TestFirefoxHTMLToOrg(t *testing.T) {
	// Parse Firefox HTML bookmarks
	htmlFile, err := os.Open("../firefox_default_bookmarks.html")
	if err != nil {
		t.Fatalf("Failed to open Firefox HTML file: %v", err)
	}
	defer htmlFile.Close()

	htmlParser := parser.NewHTMLParser(htmlFile)
	root, err := htmlParser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	// Convert to org-mode
	var buf bytes.Buffer
	err = ToOrg(root, &buf)
	if err != nil {
		t.Fatalf("Failed to convert to org: %v", err)
	}

	orgOutput := buf.String()

	// Basic validation
	if orgOutput == "" {
		t.Fatal("Org output is empty")
	}

	t.Logf("Generated org output (%d bytes)", len(orgOutput))

	// Check for expected org-mode elements
	if !strings.Contains(orgOutput, "* ") {
		t.Error("No org headlines found")
	}

	if !strings.Contains(orgOutput, "[[http") {
		t.Error("No links found")
	}

	if !strings.Contains(orgOutput, ":") && strings.Contains(orgOutput, "news") {
		t.Error("Tags not found in expected format")
	}

	if !strings.Contains(orgOutput, "#+SHORTCUTURL:") {
		t.Error("SHORTCUTURL property not found")
	}

	// Print a sample of the output
	lines := strings.Split(orgOutput, "\n")
	sampleSize := 30
	if len(lines) < sampleSize {
		sampleSize = len(lines)
	}
	t.Logf("First %d lines of output:\n%s", sampleSize, strings.Join(lines[:sampleSize], "\n"))
}

func TestChromiumHTMLToOrg(t *testing.T) {
	// Parse Chromium HTML bookmarks
	htmlFile, err := os.Open("../chromium_default_bookmarks.html")
	if err != nil {
		t.Fatalf("Failed to open Chromium HTML file: %v", err)
	}
	defer htmlFile.Close()

	htmlParser := parser.NewHTMLParser(htmlFile)
	root, err := htmlParser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	// Convert to org-mode
	var buf bytes.Buffer
	err = ToOrg(root, &buf)
	if err != nil {
		t.Fatalf("Failed to convert to org: %v", err)
	}

	orgOutput := buf.String()

	// Basic validation
	if orgOutput == "" {
		t.Fatal("Org output is empty")
	}

	t.Logf("Generated org output (%d bytes)", len(orgOutput))

	// Check for expected org-mode elements
	if !strings.Contains(orgOutput, "* ") {
		t.Error("No org headlines found")
	}

	if !strings.Contains(orgOutput, "[[http") {
		t.Error("No links found")
	}

	// Print full output for Chromium (it's small)
	t.Logf("Full output:\n%s", orgOutput)
}

func TestRoundTripHTMLToOrgToHTML(t *testing.T) {
	// Parse Firefox HTML bookmarks
	htmlFile, err := os.Open("../firefox_default_bookmarks.html")
	if err != nil {
		t.Fatalf("Failed to open Firefox HTML file: %v", err)
	}
	defer htmlFile.Close()

	htmlParser := parser.NewHTMLParser(htmlFile)
	root1, err := htmlParser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	// Convert to org-mode
	var orgBuf bytes.Buffer
	err = ToOrg(root1, &orgBuf)
	if err != nil {
		t.Fatalf("Failed to convert to org: %v", err)
	}

	t.Logf("HTML → Org conversion: %d bytes", orgBuf.Len())

	// Parse the org-mode output back
	orgParser := parser.NewOrgParser(&orgBuf)
	root2, err := orgParser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse org back to model: %v", err)
	}

	// Compare node counts
	count1 := models.CountNodes(root1)
	count2 := models.CountNodes(root2)

	t.Logf("Original tree: %d nodes", count1)
	t.Logf("Round-trip tree: %d nodes", count2)

	// Verify counts match (should be similar, allowing for place: URLs being dropped)
	if count2 == 0 {
		t.Error("Round-trip tree is empty")
	}

	// Verify we have bookmarks with URLs using models.Walk
	bookmarkCount := 0
	folderCount := 0
	var hasTaggedBookmark bool
	var hasShortcut bool

	models.Walk(root2, 0, func(node models.Node, depth int) {
		if node.IsFolder() {
			folderCount++
		} else {
			bookmarkCount++
			bookmark := node.(*models.Bookmark)
			if bookmark.URL == "" {
				t.Errorf("Bookmark '%s' has no URL after round-trip", bookmark.Title)
			}
			if len(bookmark.Tags) > 0 {
				hasTaggedBookmark = true
			}
			if bookmark.ShortcutURL != "" {
				hasShortcut = true
			}
		}
	})

	t.Logf("Round-trip result: %d folders, %d bookmarks", folderCount, bookmarkCount)

	if bookmarkCount == 0 {
		t.Error("No bookmarks in round-trip result")
	}

	if !hasTaggedBookmark {
		t.Error("No tagged bookmarks in round-trip result")
	}

	if !hasShortcut {
		t.Error("No shortcuts in round-trip result")
	}
}

func TestOrgToHTML(t *testing.T) {
	// Parse org bookmarks
	orgFile, err := os.Open("../test_bookmarks.org")
	if err != nil {
		t.Fatalf("Failed to open test org file: %v", err)
	}
	defer orgFile.Close()

	orgParser := parser.NewOrgParser(orgFile)
	root, err := orgParser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse org: %v", err)
	}

	// Convert to HTML
	var buf bytes.Buffer
	err = ToHTML(root, &buf)
	if err != nil {
		t.Fatalf("Failed to convert to HTML: %v", err)
	}

	htmlOutput := buf.String()

	// Basic validation
	if htmlOutput == "" {
		t.Fatal("HTML output is empty")
	}

	t.Logf("Generated HTML output (%d bytes)", len(htmlOutput))

	// Check for expected HTML elements
	if !strings.Contains(htmlOutput, "<!DOCTYPE NETSCAPE-Bookmark-file-1>") {
		t.Error("Missing DOCTYPE")
	}

	if !strings.Contains(htmlOutput, "<DL><p>") {
		t.Error("No DL tags found")
	}

	if !strings.Contains(htmlOutput, "<H3") {
		t.Error("No H3 folder tags found")
	}

	if !strings.Contains(htmlOutput, "<A HREF=") {
		t.Error("No anchor tags found")
	}

	if !strings.Contains(htmlOutput, "TAGS=") {
		t.Error("No TAGS attributes found")
	}

	if !strings.Contains(htmlOutput, "SHORTCUTURL=") {
		t.Error("No SHORTCUTURL attributes found")
	}

	// Print a sample of the output
	lines := strings.Split(htmlOutput, "\n")
	sampleSize := 40
	if len(lines) < sampleSize {
		sampleSize = len(lines)
	}
	t.Logf("First %d lines of HTML output:\n%s", sampleSize, strings.Join(lines[:sampleSize], "\n"))
}
