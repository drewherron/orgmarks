package models

import "time"

// Node is the interface that both Bookmark and Folder implement
// This allows building a tree structure with mixed node types
type Node interface {
	IsFolder() bool
	GetTitle() string
}

// Bookmark represents a single bookmark entry
type Bookmark struct {
	URL          string    // The bookmark URL
	Title        string    // The bookmark title/name
	Tags         []string  // Tags associated with the bookmark
	ShortcutURL  string    // Firefox SHORTCUTURL attribute (optional)
	AddDate      time.Time // When the bookmark was added
	LastModified time.Time // When the bookmark was last modified
	Description  string    // Optional description text (below the link in org-mode)
}

// Folder represents a bookmark folder/directory
type Folder struct {
	Title        string      // The folder name
	Children     []Node      // Child nodes (can be bookmarks or folders)
	AddDate      time.Time   // When the folder was created
	LastModified time.Time   // When the folder was last modified
}

// IsFolder returns false for Bookmark nodes
func (b *Bookmark) IsFolder() bool {
	return false
}

// GetTitle returns the bookmark title
func (b *Bookmark) GetTitle() string {
	return b.Title
}

// IsFolder returns true for Folder nodes
func (f *Folder) IsFolder() bool {
	return true
}

// GetTitle returns the folder title
func (f *Folder) GetTitle() string {
	return f.Title
}
