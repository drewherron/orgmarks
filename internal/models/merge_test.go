package models

import (
	"testing"
)

func TestMergeFoldersSimple(t *testing.T) {
	// Create two simple folder trees
	folder1 := &Folder{
		Title: "Root",
		Children: []Node{
			&Bookmark{Title: "Bookmark 1", URL: "https://example.com/1"},
			&Bookmark{Title: "Bookmark 2", URL: "https://example.com/2"},
		},
	}

	folder2 := &Folder{
		Title: "Root",
		Children: []Node{
			&Bookmark{Title: "Bookmark 3", URL: "https://example.com/3"},
			&Bookmark{Title: "Bookmark 4", URL: "https://example.com/4"},
		},
	}

	merged := MergeFolders(folder1, folder2)

	// Should have all 4 bookmarks
	if len(merged.Children) != 4 {
		t.Errorf("Expected 4 children, got %d", len(merged.Children))
	}

	// Verify order: folder1's bookmarks first, then folder2's
	expectedTitles := []string{"Bookmark 1", "Bookmark 2", "Bookmark 3", "Bookmark 4"}
	for i, child := range merged.Children {
		if child.GetTitle() != expectedTitles[i] {
			t.Errorf("Expected child %d to be '%s', got '%s'", i, expectedTitles[i], child.GetTitle())
		}
	}
}

func TestMergeFoldersWithMatchingSubfolders(t *testing.T) {
	// Create two trees with matching "Finance" folders
	folder1 := &Folder{
		Title: "Root",
		Children: []Node{
			&Folder{
				Title: "Finance",
				Children: []Node{
					&Bookmark{Title: "Bank 1", URL: "https://bank1.com"},
				},
			},
		},
	}

	folder2 := &Folder{
		Title: "Root",
		Children: []Node{
			&Folder{
				Title: "Finance",
				Children: []Node{
					&Bookmark{Title: "Bank 2", URL: "https://bank2.com"},
				},
			},
		},
	}

	merged := MergeFolders(folder1, folder2)

	// Should have 1 Finance folder (merged)
	if len(merged.Children) != 1 {
		t.Errorf("Expected 1 child (merged Finance folder), got %d", len(merged.Children))
	}

	financeFolder := merged.Children[0].(*Folder)
	if financeFolder.Title != "Finance" {
		t.Errorf("Expected 'Finance' folder, got '%s'", financeFolder.Title)
	}

	// Finance folder should have both bookmarks
	if len(financeFolder.Children) != 2 {
		t.Errorf("Expected Finance folder to have 2 bookmarks, got %d", len(financeFolder.Children))
	}
}

func TestMergeFoldersCaseInsensitive(t *testing.T) {
	// Create two trees with "Finance" and "finance" - should match
	folder1 := &Folder{
		Title: "Root",
		Children: []Node{
			&Folder{
				Title: "Finance",
				Children: []Node{
					&Bookmark{Title: "Bank 1", URL: "https://bank1.com"},
				},
			},
		},
	}

	folder2 := &Folder{
		Title: "Root",
		Children: []Node{
			&Folder{
				Title: "finance",
				Children: []Node{
					&Bookmark{Title: "Bank 2", URL: "https://bank2.com"},
				},
			},
		},
	}

	merged := MergeFolders(folder1, folder2)

	// Should have 1 folder (case-insensitive match)
	if len(merged.Children) != 1 {
		t.Errorf("Expected 1 merged folder, got %d", len(merged.Children))
	}

	financeFolder := merged.Children[0].(*Folder)
	// Should preserve the first folder's capitalization
	if financeFolder.Title != "Finance" {
		t.Errorf("Expected 'Finance' (from folder1), got '%s'", financeFolder.Title)
	}

	// Should have both bookmarks
	if len(financeFolder.Children) != 2 {
		t.Errorf("Expected 2 bookmarks, got %d", len(financeFolder.Children))
	}
}

func TestMergeFoldersWithDifferentSubfolders(t *testing.T) {
	// Create two trees with different subfolders
	folder1 := &Folder{
		Title: "Root",
		Children: []Node{
			&Folder{
				Title: "Finance",
				Children: []Node{
					&Bookmark{Title: "Bank 1", URL: "https://bank1.com"},
				},
			},
		},
	}

	folder2 := &Folder{
		Title: "Root",
		Children: []Node{
			&Folder{
				Title: "Shopping",
				Children: []Node{
					&Bookmark{Title: "Amazon", URL: "https://amazon.com"},
				},
			},
		},
	}

	merged := MergeFolders(folder1, folder2)

	// Should have 2 folders (both preserved)
	if len(merged.Children) != 2 {
		t.Errorf("Expected 2 folders, got %d", len(merged.Children))
	}

	// Verify both folders exist
	foundFinance := false
	foundShopping := false
	for _, child := range merged.Children {
		if child.IsFolder() {
			folder := child.(*Folder)
			if folder.Title == "Finance" {
				foundFinance = true
			}
			if folder.Title == "Shopping" {
				foundShopping = true
			}
		}
	}

	if !foundFinance {
		t.Error("Finance folder not found in merged result")
	}
	if !foundShopping {
		t.Error("Shopping folder not found in merged result")
	}
}

func TestMergeFoldersDeepNesting(t *testing.T) {
	// Create two trees with deeply nested matching folders
	folder1 := &Folder{
		Title: "Root",
		Children: []Node{
			&Folder{
				Title: "Work",
				Children: []Node{
					&Folder{
						Title: "Projects",
						Children: []Node{
							&Bookmark{Title: "Project A", URL: "https://project-a.com"},
						},
					},
				},
			},
		},
	}

	folder2 := &Folder{
		Title: "Root",
		Children: []Node{
			&Folder{
				Title: "Work",
				Children: []Node{
					&Folder{
						Title: "Projects",
						Children: []Node{
							&Bookmark{Title: "Project B", URL: "https://project-b.com"},
						},
					},
				},
			},
		},
	}

	merged := MergeFolders(folder1, folder2)

	// Navigate to the deeply nested Projects folder
	if len(merged.Children) != 1 {
		t.Fatalf("Expected 1 child at root, got %d", len(merged.Children))
	}

	workFolder := merged.Children[0].(*Folder)
	if workFolder.Title != "Work" {
		t.Errorf("Expected 'Work' folder, got '%s'", workFolder.Title)
	}

	if len(workFolder.Children) != 1 {
		t.Fatalf("Expected 1 child in Work folder, got %d", len(workFolder.Children))
	}

	projectsFolder := workFolder.Children[0].(*Folder)
	if projectsFolder.Title != "Projects" {
		t.Errorf("Expected 'Projects' folder, got '%s'", projectsFolder.Title)
	}

	// Projects folder should have both bookmarks
	if len(projectsFolder.Children) != 2 {
		t.Errorf("Expected 2 bookmarks in Projects folder, got %d", len(projectsFolder.Children))
	}
}

func TestMergeFoldersMixedContent(t *testing.T) {
	// Mix of matching folders, non-matching folders, and bookmarks
	folder1 := &Folder{
		Title: "Root",
		Children: []Node{
			&Bookmark{Title: "Root Bookmark 1", URL: "https://root1.com"},
			&Folder{
				Title: "Shared",
				Children: []Node{
					&Bookmark{Title: "Shared 1", URL: "https://shared1.com"},
				},
			},
			&Folder{
				Title: "Only in F1",
				Children: []Node{
					&Bookmark{Title: "Unique 1", URL: "https://unique1.com"},
				},
			},
		},
	}

	folder2 := &Folder{
		Title: "Root",
		Children: []Node{
			&Bookmark{Title: "Root Bookmark 2", URL: "https://root2.com"},
			&Folder{
				Title: "Shared",
				Children: []Node{
					&Bookmark{Title: "Shared 2", URL: "https://shared2.com"},
				},
			},
			&Folder{
				Title: "Only in F2",
				Children: []Node{
					&Bookmark{Title: "Unique 2", URL: "https://unique2.com"},
				},
			},
		},
	}

	merged := MergeFolders(folder1, folder2)

	// Count children: 2 bookmarks + 3 folders
	if len(merged.Children) != 5 {
		t.Errorf("Expected 5 children (2 bookmarks + 3 folders), got %d", len(merged.Children))
	}

	// Verify we have exactly 2 bookmarks at root
	bookmarkCount := 0
	folderCount := 0
	for _, child := range merged.Children {
		if child.IsFolder() {
			folderCount++
		} else {
			bookmarkCount++
		}
	}

	if bookmarkCount != 2 {
		t.Errorf("Expected 2 bookmarks at root, got %d", bookmarkCount)
	}
	if folderCount != 3 {
		t.Errorf("Expected 3 folders at root, got %d", folderCount)
	}

	// Verify the "Shared" folder has 2 bookmarks
	for _, child := range merged.Children {
		if child.IsFolder() {
			folder := child.(*Folder)
			if folder.Title == "Shared" {
				if len(folder.Children) != 2 {
					t.Errorf("Expected Shared folder to have 2 bookmarks, got %d", len(folder.Children))
				}
			}
		}
	}
}

func TestMergeFoldersPreservesOrder(t *testing.T) {
	// Verify that folder1's content comes before folder2's content
	folder1 := &Folder{
		Title: "Root",
		Children: []Node{
			&Bookmark{Title: "A1", URL: "https://a1.com"},
			&Bookmark{Title: "A2", URL: "https://a2.com"},
		},
	}

	folder2 := &Folder{
		Title: "Root",
		Children: []Node{
			&Bookmark{Title: "B1", URL: "https://b1.com"},
			&Bookmark{Title: "B2", URL: "https://b2.com"},
		},
	}

	merged := MergeFolders(folder1, folder2)

	expectedOrder := []string{"A1", "A2", "B1", "B2"}
	for i, child := range merged.Children {
		if child.GetTitle() != expectedOrder[i] {
			t.Errorf("Expected position %d to be '%s', got '%s'", i, expectedOrder[i], child.GetTitle())
		}
	}
}
