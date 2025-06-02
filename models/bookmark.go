package models

import "time"

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
