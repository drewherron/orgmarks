package parser

import (
	"os"
	"testing"

	"github.com/drewherron/orgmarks/models"
)

func TestParseFirefoxBookmarks(t *testing.T) {
	// Open the Firefox sample file
	file, err := os.Open("../firefox_default_bookmarks.html")
	if err != nil {
		t.Fatalf("Failed to open Firefox bookmark file: %v", err)
	}
	defer file.Close()

	// Parse the file
	parser := NewHTMLParser(file)
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse Firefox bookmarks: %v", err)
	}

	// Basic validation
	if root == nil {
		t.Fatal("Root folder is nil")
	}

	if root.Title == "" {
		t.Error("Root folder has no title")
	}

	// Count nodes
	nodeCount := models.CountNodes(root)
	if nodeCount == 0 {
		t.Error("No nodes found in bookmark tree")
	}

	t.Logf("Parsed %d total nodes from Firefox bookmarks", nodeCount)

	// Verify we can walk the tree
	folderCount := 0
	bookmarkCount := 0
	models.Walk(root, 0, func(node models.Node, depth int) {
		if node.IsFolder() {
			folderCount++
		} else {
			bookmarkCount++
			// Verify bookmark has URL
			bookmark := node.(*models.Bookmark)
			if bookmark.URL == "" {
				t.Errorf("Bookmark '%s' has no URL", bookmark.Title)
			}
		}
	})

	t.Logf("Found %d folders and %d bookmarks", folderCount, bookmarkCount)

	if folderCount == 0 {
		t.Error("No folders found")
	}

	if bookmarkCount == 0 {
		t.Error("No bookmarks found")
	}

	// Verify specific features were parsed
	foundTags := false
	foundShortcut := false
	models.Walk(root, 0, func(node models.Node, depth int) {
		if !node.IsFolder() {
			bookmark := node.(*models.Bookmark)
			if len(bookmark.Tags) > 0 {
				foundTags = true
				t.Logf("Found bookmark with tags: %s -> %v", bookmark.Title, bookmark.Tags)
			}
			if bookmark.ShortcutURL != "" {
				foundShortcut = true
				t.Logf("Found bookmark with shortcut: %s -> %s", bookmark.Title, bookmark.ShortcutURL)
			}
		}
	})

	if !foundTags {
		t.Error("No bookmarks with tags found")
	}

	if !foundShortcut {
		t.Error("No bookmarks with shortcuts found")
	}
}
