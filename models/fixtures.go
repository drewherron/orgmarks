package models

import "time"

// SampleBookmarkTree returns a sample bookmark tree for testing
func SampleBookmarkTree() *Folder {
	root := &Folder{
		Title:        "Bookmarks Menu",
		AddDate:      time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		LastModified: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
	}

	// Bookmarks Toolbar folder
	toolbar := &Folder{
		Title:        "Bookmarks Toolbar",
		AddDate:      time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		LastModified: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
	}

	// Add some bookmarks to toolbar
	toolbar.AddChild(&Bookmark{
		Title:        "Fedora Docs",
		URL:          "https://docs.fedoraproject.org/",
		AddDate:      time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		LastModified: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
	})

	toolbar.AddChild(&Bookmark{
		Title:        "Fedora Magazine",
		URL:          "https://fedoramagazine.org/",
		Tags:         []string{"news"},
		ShortcutURL:  "magazine",
		AddDate:      time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		LastModified: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
	})

	// Email folder
	email := &Folder{
		Title:        "Email",
		AddDate:      time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		LastModified: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
	}

	email.AddChild(&Bookmark{
		Title:        "Gmail",
		URL:          "https://mail.google.com",
		Tags:         []string{"email", "google"},
		Description:  "Personal email",
		AddDate:      time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		LastModified: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
	})

	// Nested folder structure: Fedora Project > More Downloads
	fedoraProject := &Folder{
		Title:        "Fedora Project",
		AddDate:      time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		LastModified: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
	}

	moreDownloads := &Folder{
		Title:        "More Downloads",
		AddDate:      time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		LastModified: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
	}

	moreDownloads.AddChild(&Bookmark{
		Title:        "Fedora Spins",
		URL:          "https://spins.fedoraproject.org/",
		Tags:         []string{"cinnamon", "kde", "lxde", "lxqt", "mate", "soas", "xfce"},
		AddDate:      time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		LastModified: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
	})

	fedoraProject.AddChild(moreDownloads)
	toolbar.AddChild(fedoraProject)

	// Build the root structure
	root.AddChild(toolbar)
	root.AddChild(email)

	return root
}
