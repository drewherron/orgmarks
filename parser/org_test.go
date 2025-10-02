package parser

import (
	"os"
	"testing"

	"github.com/drewherron/orgmarks/models"
)

func TestParseOrgBookmarks(t *testing.T) {
	// Open the test org file
	file, err := os.Open("../test_bookmarks.org")
	if err != nil {
		t.Fatalf("Failed to open test org file: %v", err)
	}
	defer file.Close()

	// Parse the file
	parser := NewOrgParser(file)
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse org bookmarks: %v", err)
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

	t.Logf("Parsed %d total nodes from org bookmarks", nodeCount)

	// Verify we can walk the tree
	folderCount := 0
	bookmarkCount := 0
	models.Walk(root, 0, func(node models.Node, depth int) {
		if node.IsFolder() {
			folderCount++
			t.Logf("Folder at depth %d: %s", depth, node.GetTitle())
		} else {
			bookmarkCount++
			bookmark := node.(*models.Bookmark)
			t.Logf("Bookmark at depth %d: %s -> %s", depth, bookmark.Title, bookmark.URL)

			// Verify bookmark has URL
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
	foundDescription := false
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
			if bookmark.Description != "" {
				foundDescription = true
				t.Logf("Found bookmark with description: %s -> %s", bookmark.Title, bookmark.Description)
			}
		}
	})

	if !foundTags {
		t.Error("No bookmarks with tags found")
	}

	if !foundShortcut {
		t.Error("No bookmarks with shortcuts found")
	}

	if !foundDescription {
		t.Error("No bookmarks with descriptions found")
	}
}
