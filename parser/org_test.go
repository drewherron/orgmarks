package parser

import (
	"os"
	"strings"
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

// TestParseEmptyOrg tests parsing an empty org file
func TestParseEmptyOrg(t *testing.T) {
	org := ``

	parser := NewOrgParser(strings.NewReader(org))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse empty org: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	// Should have just the root folder, no children
	if len(root.Children) != 0 {
		t.Errorf("Expected 0 children, got %d", len(root.Children))
	}
}

// TestParseOrgWithOnlyFolders tests parsing org file with only folders (no bookmarks)
func TestParseOrgWithOnlyFolders(t *testing.T) {
	org := `* Folder One
** Subfolder A
** Subfolder B
* Folder Two`

	parser := NewOrgParser(strings.NewReader(org))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse org with only folders: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	// Count folders and bookmarks
	folderCount := 0
	bookmarkCount := 0
	models.Walk(root, 0, func(node models.Node, depth int) {
		if node.IsFolder() {
			folderCount++
		} else {
			bookmarkCount++
		}
	})

	if bookmarkCount != 0 {
		t.Errorf("Expected 0 bookmarks, got %d", bookmarkCount)
	}

	if folderCount != 5 { // root + 2 top-level + 2 subfolders
		t.Errorf("Expected 5 folders, got %d", folderCount)
	}
}

// TestParseOrgWithSpecialCharacters tests parsing org with special characters in titles and URLs
func TestParseOrgWithSpecialCharacters(t *testing.T) {
	org := `* Folder with [brackets] & special chars
** Bookmark with "quotes" & <angle>
[[https://example.com/page?foo=1&bar=2]]
Some description with * asterisks * and more`

	parser := NewOrgParser(strings.NewReader(org))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse org with special characters: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	if len(root.Children) == 0 {
		t.Fatal("No children found")
	}

	folder := root.Children[0].(*models.Folder)
	if !strings.Contains(folder.Title, "[brackets]") {
		t.Errorf("Expected folder title to contain '[brackets]', got: %s", folder.Title)
	}

	if len(folder.Children) == 0 {
		t.Fatal("No bookmarks in folder")
	}

	bookmark := folder.Children[0].(*models.Bookmark)
	if !strings.Contains(bookmark.Title, "\"quotes\"") {
		t.Errorf("Expected bookmark title to contain '\"quotes\"', got: %s", bookmark.Title)
	}
	if !strings.Contains(bookmark.URL, "&") {
		t.Errorf("Expected bookmark URL to contain '&', got: %s", bookmark.URL)
	}
	if !strings.Contains(bookmark.Description, "*") {
		t.Errorf("Expected description to contain '*', got: %s", bookmark.Description)
	}
}

// TestParseOrgWithLongTagList tests parsing bookmarks with many tags
func TestParseOrgWithLongTagList(t *testing.T) {
	org := `* Bookmark with many tags :tag1:tag2:tag3:tag4:tag5:tag6:tag7:tag8:tag9:tag10:tag11:tag12:
[[https://example.com]]`

	parser := NewOrgParser(strings.NewReader(org))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse org with long tag list: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	if len(root.Children) == 0 {
		t.Fatal("No children found")
	}

	bookmark := root.Children[0].(*models.Bookmark)
	if len(bookmark.Tags) != 12 {
		t.Errorf("Expected 12 tags, got %d: %v", len(bookmark.Tags), bookmark.Tags)
	}
}

// TestParseOrgDeeplyNested tests parsing deeply nested folder structures
func TestParseOrgDeeplyNested(t *testing.T) {
	org := `* Level 1
** Level 2
*** Level 3
**** Level 4
***** Level 5
****** Level 6
******* Deep Bookmark
[[https://deep.example.com]]`

	parser := NewOrgParser(strings.NewReader(org))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse deeply nested org: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	// Navigate to the deepest bookmark and verify it exists
	maxDepth := 0
	foundDeepBookmark := false
	models.Walk(root, 0, func(node models.Node, depth int) {
		if depth > maxDepth {
			maxDepth = depth
		}
		if !node.IsFolder() && depth == 7 {
			foundDeepBookmark = true
		}
	})

	if maxDepth < 6 {
		t.Errorf("Expected max depth of at least 6, got %d", maxDepth)
	}

	if !foundDeepBookmark {
		t.Error("Did not find deeply nested bookmark at expected depth")
	}
}

// TestParseOrgWithEmptyFolder tests parsing folders with no bookmarks
func TestParseOrgWithEmptyFolder(t *testing.T) {
	org := `* Empty Folder
* Another Folder
** Bookmark
[[https://example.com]]`

	parser := NewOrgParser(strings.NewReader(org))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse org with empty folder: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	if len(root.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(root.Children))
	}

	// First child should be the empty folder
	if !root.Children[0].IsFolder() {
		t.Error("First child should be a folder")
	}

	emptyFolder := root.Children[0].(*models.Folder)
	if len(emptyFolder.Children) != 0 {
		t.Errorf("Expected empty folder to have 0 children, got %d", len(emptyFolder.Children))
	}
}

// TestParseOrgWithLinkTitleFormat tests parsing [[URL][Title]] format
func TestParseOrgWithLinkTitleFormat(t *testing.T) {
	org := `* Bookmark with separate title
[[https://example.com][Example Website]]
This is a description`

	parser := NewOrgParser(strings.NewReader(org))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse org with link title format: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	if len(root.Children) == 0 {
		t.Fatal("No children found")
	}

	bookmark := root.Children[0].(*models.Bookmark)
	if bookmark.URL != "https://example.com" {
		t.Errorf("Expected URL 'https://example.com', got: %s", bookmark.URL)
	}
	// Note: headline title should be used, not link title
	if bookmark.Title != "Bookmark with separate title" {
		t.Errorf("Expected title 'Bookmark with separate title', got: %s", bookmark.Title)
	}
	if bookmark.Description != "This is a description" {
		t.Errorf("Expected description 'This is a description', got: %s", bookmark.Description)
	}
}

// TestParseOrgWithShortcutURL tests parsing #+SHORTCUTURL property
func TestParseOrgWithShortcutURL(t *testing.T) {
	org := `* Bookmark with shortcut
#+SHORTCUTURL: myshortcut
[[https://example.com]]`

	parser := NewOrgParser(strings.NewReader(org))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse org with shortcut URL: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	if len(root.Children) == 0 {
		t.Fatal("No children found")
	}

	bookmark := root.Children[0].(*models.Bookmark)
	if bookmark.ShortcutURL != "myshortcut" {
		t.Errorf("Expected shortcut 'myshortcut', got: %s", bookmark.ShortcutURL)
	}
}

// TestParseOrgMixedContent tests parsing org with mixed folders and bookmarks at same level
func TestParseOrgMixedContent(t *testing.T) {
	org := `* Folder One
** Bookmark One
[[https://one.example.com]]
* Bookmark Two
[[https://two.example.com]]
* Folder Two
** Bookmark Three
[[https://three.example.com]]`

	parser := NewOrgParser(strings.NewReader(org))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse org with mixed content: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	if len(root.Children) != 3 {
		t.Errorf("Expected 3 children at root level, got %d", len(root.Children))
	}

	// Check order: Folder One, Bookmark Two, Folder Two
	if !root.Children[0].IsFolder() {
		t.Error("First child should be a folder")
	}
	if root.Children[1].IsFolder() {
		t.Error("Second child should be a bookmark")
	}
	if !root.Children[2].IsFolder() {
		t.Error("Third child should be a folder")
	}
}
