package parser

import (
	"os"
	"strings"
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

func TestParseChromiumBookmarks(t *testing.T) {
	// Open the Chromium sample file
	file, err := os.Open("../chromium_default_bookmarks.html")
	if err != nil {
		t.Fatalf("Failed to open Chromium bookmark file: %v", err)
	}
	defer file.Close()

	// Parse the file
	parser := NewHTMLParser(file)
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse Chromium bookmarks: %v", err)
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

	t.Logf("Parsed %d total nodes from Chromium bookmarks", nodeCount)

	// Verify we can walk the tree
	folderCount := 0
	bookmarkCount := 0
	models.Walk(root, 0, func(node models.Node, depth int) {
		if node.IsFolder() {
			folderCount++
			t.Logf("Folder at depth %d: %s", depth, node.GetTitle())
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
}

// TestParseEmptyHTML tests parsing an empty/minimal HTML file
func TestParseEmptyHTML(t *testing.T) {
	html := `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
</DL><p>`

	parser := NewHTMLParser(strings.NewReader(html))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse empty HTML: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	// Should have just the root folder, no children
	if len(root.Children) != 0 {
		t.Errorf("Expected 0 children, got %d", len(root.Children))
	}
}

// TestParseHTMLWithSpecialCharacters tests parsing bookmarks with special HTML characters
func TestParseHTMLWithSpecialCharacters(t *testing.T) {
	html := `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3>Folder with &amp; ampersand &lt;&gt;</H3>
    <DL><p>
        <DT><A HREF="https://example.com/page?foo=1&amp;bar=2" ADD_DATE="1234567890">Link with &quot;quotes&quot; &amp; ampersands</A>
    </DL><p>
</DL><p>`

	parser := NewHTMLParser(strings.NewReader(html))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse HTML with special characters: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	// Check folder title
	if len(root.Children) == 0 {
		t.Fatal("No children found")
	}

	folder := root.Children[0].(*models.Folder)
	if !strings.Contains(folder.Title, "&") {
		t.Errorf("Expected folder title to contain '&', got: %s", folder.Title)
	}
	if !strings.Contains(folder.Title, "<>") {
		t.Errorf("Expected folder title to contain '<>', got: %s", folder.Title)
	}

	// Check bookmark
	if len(folder.Children) == 0 {
		t.Fatal("No bookmarks in folder")
	}

	bookmark := folder.Children[0].(*models.Bookmark)
	if !strings.Contains(bookmark.Title, "\"") {
		t.Errorf("Expected bookmark title to contain '\"', got: %s", bookmark.Title)
	}
	if !strings.Contains(bookmark.URL, "&") {
		t.Errorf("Expected bookmark URL to contain '&', got: %s", bookmark.URL)
	}
}

// TestParseMalformedHTML tests parsing HTML with missing or malformed tags
func TestParseMalformedHTML(t *testing.T) {
	// Missing closing tags - parser should still work
	html := `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3>Folder</H3>
    <DL><p>
        <DT><A HREF="https://example.com">Bookmark
</DL><p>`

	parser := NewHTMLParser(strings.NewReader(html))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse malformed HTML: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	// Should still find the folder and bookmark
	nodeCount := models.CountNodes(root)
	if nodeCount < 2 {
		t.Errorf("Expected at least 2 nodes (folder + bookmark), got %d", nodeCount)
	}
}

// TestParseHTMLWithEmptyFolder tests parsing folders with no bookmarks
func TestParseHTMLWithEmptyFolder(t *testing.T) {
	html := `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3>Empty Folder</H3>
    <DL><p>
    </DL><p>
    <DT><A HREF="https://example.com">Bookmark</A>
</DL><p>`

	parser := NewHTMLParser(strings.NewReader(html))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse HTML with empty folder: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	if len(root.Children) != 2 {
		t.Errorf("Expected 2 children (empty folder + bookmark), got %d", len(root.Children))
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

// TestParseHTMLWithLongTagList tests parsing bookmarks with many tags
func TestParseHTMLWithLongTagList(t *testing.T) {
	html := `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><A HREF="https://example.com" TAGS="tag1,tag2,tag3,tag4,tag5,tag6,tag7,tag8,tag9,tag10,tag11,tag12">Many Tags</A>
</DL><p>`

	parser := NewHTMLParser(strings.NewReader(html))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse HTML with long tag list: %v", err)
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

// TestParseHTMLDeeplyNested tests parsing deeply nested folder structures
func TestParseHTMLDeeplyNested(t *testing.T) {
	html := `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3>Level 1</H3>
    <DL><p>
        <DT><H3>Level 2</H3>
        <DL><p>
            <DT><H3>Level 3</H3>
            <DL><p>
                <DT><H3>Level 4</H3>
                <DL><p>
                    <DT><H3>Level 5</H3>
                    <DL><p>
                        <DT><A HREF="https://deep.example.com">Deep Bookmark</A>
                    </DL><p>
                </DL><p>
            </DL><p>
        </DL><p>
    </DL><p>
</DL><p>`

	parser := NewHTMLParser(strings.NewReader(html))
	root, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse deeply nested HTML: %v", err)
	}

	if root == nil {
		t.Fatal("Root folder is nil")
	}

	// Navigate to the deepest bookmark and verify it exists
	maxDepth := 0
	models.Walk(root, 0, func(node models.Node, depth int) {
		if depth > maxDepth {
			maxDepth = depth
		}
		if !node.IsFolder() && depth < 6 {
			t.Errorf("Expected bookmark at depth 6, found at depth %d", depth)
		}
	})

	if maxDepth < 5 {
		t.Errorf("Expected max depth of at least 5, got %d", maxDepth)
	}
}
