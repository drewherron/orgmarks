package models

import "strings"

// MergeFolders merges two folder trees into a single tree.
// Folders with matching titles (case-insensitive) have their children combined.
// Children from folder1 are added first, followed by children from folder2.
// This ordering ensures that when deduplication is applied afterward,
// bookmarks from folder1 take precedence over duplicates in folder2.
func MergeFolders(folder1, folder2 *Folder) *Folder {
	// Create the merged root folder
	merged := &Folder{
		Title:        folder1.Title,
		AddDate:      folder1.AddDate,
		LastModified: folder1.LastModified,
	}

	// Build a map of folder1's children by normalized title (for folders only)
	folderMap := make(map[string]*Folder)
	for _, child := range folder1.Children {
		if child.IsFolder() {
			subfolder := child.(*Folder)
			normalizedTitle := strings.ToLower(strings.TrimSpace(subfolder.Title))
			folderMap[normalizedTitle] = subfolder
		}
	}

	// First, add all children from folder1
	for _, child := range folder1.Children {
		if child.IsFolder() {
			// Folders will be merged later, add a placeholder for now
			continue
		} else {
			// Add bookmarks directly
			merged.AddChild(child)
		}
	}

	// Build a list of folders from folder2 that need to be merged or added
	folder2Subfolders := make([]*Folder, 0)
	for _, child := range folder2.Children {
		if child.IsFolder() {
			folder2Subfolders = append(folder2Subfolders, child.(*Folder))
		} else {
			// Add bookmarks from folder2
			merged.AddChild(child)
		}
	}

	// Now handle folder merging
	// First, recursively merge folders that exist in both trees
	processedFolders := make(map[string]bool)
	for _, folder1Child := range folder1.Children {
		if !folder1Child.IsFolder() {
			continue
		}
		subfolder1 := folder1Child.(*Folder)
		normalizedTitle := strings.ToLower(strings.TrimSpace(subfolder1.Title))

		// Look for a matching folder in folder2
		var matchingFolder2 *Folder
		for _, folder2Child := range folder2Subfolders {
			folder2Normalized := strings.ToLower(strings.TrimSpace(folder2Child.Title))
			if folder2Normalized == normalizedTitle {
				matchingFolder2 = folder2Child
				break
			}
		}

		if matchingFolder2 != nil {
			// Recursively merge the two folders
			mergedSubfolder := MergeFolders(subfolder1, matchingFolder2)
			merged.AddChild(mergedSubfolder)
			processedFolders[normalizedTitle] = true
		} else {
			// Only in folder1, add as-is
			merged.AddChild(subfolder1)
			processedFolders[normalizedTitle] = true
		}
	}

	// Add folders from folder2 that weren't matched
	for _, folder2Child := range folder2Subfolders {
		normalizedTitle := strings.ToLower(strings.TrimSpace(folder2Child.Title))
		if !processedFolders[normalizedTitle] {
			merged.AddChild(folder2Child)
		}
	}

	return merged
}
