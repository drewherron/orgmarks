package models

// Deduplicate removes duplicate bookmarks with the same URL, keeping only the first occurrence.
// It walks the tree in depth-first order and tracks seen URLs.
func Deduplicate(root *Folder) {
	seen := make(map[string]bool)
	deduplicateFolder(root, seen)
}

// deduplicateFolder recursively removes duplicate bookmarks from a folder
func deduplicateFolder(folder *Folder, seen map[string]bool) {
	// Filter children, keeping only non-duplicate bookmarks
	filtered := make([]Node, 0, len(folder.Children))

	for _, child := range folder.Children {
		if child.IsFolder() {
			// Recursively deduplicate subfolders
			subfolder := child.(*Folder)
			deduplicateFolder(subfolder, seen)
			filtered = append(filtered, child)
		} else {
			// Check if bookmark URL has been seen
			bookmark := child.(*Bookmark)
			if !seen[bookmark.URL] {
				seen[bookmark.URL] = true
				filtered = append(filtered, child)
			}
			// If URL was already seen, skip this bookmark (don't append)
		}
	}

	folder.Children = filtered
}

// RemoveEmptyFolders recursively removes folders that have no children.
// This is useful after deduplication when folders may have become empty.
func RemoveEmptyFolders(root *Folder) {
	removeEmptyFoldersRecursive(root)
}

// removeEmptyFoldersRecursive recursively removes empty folders from a folder
func removeEmptyFoldersRecursive(folder *Folder) {
	// First, recursively process all subfolders
	for _, child := range folder.Children {
		if child.IsFolder() {
			subfolder := child.(*Folder)
			removeEmptyFoldersRecursive(subfolder)
		}
	}

	// Then filter out empty folders
	filtered := make([]Node, 0, len(folder.Children))
	for _, child := range folder.Children {
		if child.IsFolder() {
			subfolder := child.(*Folder)
			// Only keep folders that have children
			if len(subfolder.Children) > 0 {
				filtered = append(filtered, child)
			}
			// Empty folders are skipped (not appended)
		} else {
			// Always keep bookmarks
			filtered = append(filtered, child)
		}
	}

	folder.Children = filtered
}
