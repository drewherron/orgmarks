# Orgmarks

A program allowing bidirectional conversion between the Netscape Bookmark HTML format (used by Firefox and Chrome) and Org-mode format, for easy bookmark organization in Emacs.

## Features

- **Bidirectional conversion**: HTML ↔ Org-mode
- **Deduplication**: Optional removal of duplicate URLs
- **Merging**: Multiple inputs produce one Org file
- **Nested folder support**: Handles nested bookmark hierarchies
- **Standards-compliant**: Compatible with Firefox, Chrome, and Chromium bookmark exports

## Why?

- **Emacs users**: Organize bookmarks in Org-mode with all of Emacs' text editing power
- **Version control**: Track bookmark changes with Git using plain-text Org format
- **Bulk editing**: Use Emacs' powerful text manipulation
- **Backup**: Maintain human-readable backups of your bookmarks
- **Deduplication**: Your browser can't do this

## Installation

### From Source

Requires Go 1.16 or later:

```bash
git clone https://github.com/drewherron/orgmarks.git
cd orgmarks
go build
```

This will create the `orgmarks` binary in the current directory.

### Install to PATH

```bash
go install github.com/drewherron/orgmarks@latest
```

### Emacs Tip

One of the main reasons I wrote this program was to be able to use org-refile on my bookmarks. I didn't want to use all my current targets when organizing bookmarks, so I added this to my `init.el`:

```elisp
;; Custom refile function for bookmarks files
(defun my/verify-refile-target-is-folder ()
  "Return t if the current heading has no direct link (children don't matter)."
  (save-excursion
    (org-back-to-heading t)
    (let ((heading-end (save-excursion
                         (outline-next-heading)
                         (point))))
      ;; Skip the headline itself and check only until the next heading
      (forward-line)
      ;; Only check this entry's content, not children
      (let ((next-heading (save-excursion
                            (if (re-search-forward "^\\*+ " heading-end t)
                                (match-beginning 0)
                              heading-end))))
        (not (re-search-forward org-link-bracket-re next-heading t))))))

(defun my/org-refile-bookmarks ()
  "Refile to bookmarks.org when in bookmarks.org or bookmarks_source.org."
  (interactive)
  (if (and buffer-file-name
           (string-match-p "bookmarks\\(_source\\)?\\.org$" buffer-file-name))
      ;; We're in bookmarks.org or bookmarks_source.org
      ;; Always refile to bookmarks.org (the clean target)
      (let ((org-refile-targets '(("~/org/bookmarks.org" :maxlevel . 10)))
            (org-refile-target-verify-function 'my/verify-refile-target-is-folder))
        (org-refile))
    ;; Not in a bookmarks file, use normal refile
    (org-refile)))

;; Bind this to C-c C-w (the normal org-refile binding)
(with-eval-after-load 'org
  (define-key org-mode-map (kbd "C-c C-w") 'my/org-refile-bookmarks))
```

This automatically detects when you're working in a bookmarks file (named either `bookmarks.org` or `bookmark_source.org`) and only shows folders (no bookmarks) as refile targets. This way, you can work in your bookmarks file (I use `~/org/bookmarks.org`) and only refile to within that same file. Or, you can convert all your messy, disorganized bookmarks to `bookmarks_source.org`, then use `C-c C-w` to refile individual bookmarks one at a time into a clean `~/org/bookmarks.org` file (which will always be the target). Feel free to change that path, of course.

## Usage

### Basic Conversion

Convert HTML bookmarks to Org-mode:

```bash
orgmarks -i bookmarks.html -o bookmarks.org
```

Convert Org-mode bookmarks to HTML:

```bash
orgmarks -i bookmarks.org -o bookmarks.html
```

### Deduplication

Remove duplicate URLs (keeps first occurrence):

```bash
orgmarks -i bookmarks.html -o bookmarks.org --deduplicate
```

### Delete Empty Folders

Remove empty folders after processing (might be especially useful after deduplication):

```bash
orgmarks -i bookmarks.html -o bookmarks.org --deduplicate --delete-empty
```

### Merging Files

You can merge multiple bookmark files into a single Org file just by supplying multiple input files. Folders with matching names (case-insensitive) are combined, and their bookmarks are merged together. You can merge any combination of `.org` and `.html` files, but with multiple inputs the output must be an Org file:

```bash
# Merge two org files
orgmarks -i organized.org -i new_bookmarks.org -o final.org

# Merge HTML bookmarks into an existing org file
orgmarks -i bookmarks.org -i browser_export.html -o bookmarks.org

# Merge multiple files with deduplication
orgmarks -i file1.org -i file2.org -i file3.html -o merged.org --deduplicate

# Merge and clean up
orgmarks -i primary.org -i additions.html -o final.org --deduplicate --delete-empty
```

**Note**: When merging with `--deduplicate`, bookmarks from the first input file take precedence over duplicates in subsequent files. This means you can list your organized bookmarks first, then add new bookmarks from your browser, and any duplicates will keep the version from your organized file.

**Another Note**: Remember that the bookmarks are added first and then deduplicated. If there are folders containing only duplicate files, it will look like we're just adding empty folders to our output file. This is actually intentional, and if you want to remove all empty folders in the output file you can use `--delete-empty`.

### Version Information

```bash
orgmarks --version
```

### Help

```bash
orgmarks --help
```

## org-mode Format

orgmarks uses the following org-mode conventions for bookmarks:

### Folders

Folders are represented as Org headlines:

```org
* Bookmarks Toolbar
** Technology
*** Programming
```

### Bookmarks

Bookmarks are headlines with a link:

```org
** My Bookmark Title
[[https://example.com]]
```

### Tags

Tags are appended to the headline in the standard Org format:

```org
** GNU.org                                                                :tech:
[[https://www.gnu.org]]
```

### Shortcut URLs (Keywords)

Shortcut URLs are stored as properties:

```org
** Wiktionary
#+SHORTCUTURL: wk
[[https://www.wiktionary.org/]]
```

These are basically aliases that can be entered in the address bar to visit a URL. They may be called something else (keyword, nickname, etc.) depending on browser. Or, the functionality may be missing entirely.

PS: Zen Browser looks pretty nice.

### Descriptions

Text after the link is treated as a description:

```org
** Gmail                                                           :email:google:
[[https://mail.google.com]]
Personal email account
```

Chrome doesn't support this, and it was removed in Firefox, but who knows, it could come back someday...

### Complete Example

```org
* Bookmarks Toolbar
** Wikipedia                                                          :reference:
[[https://wikipedia.org/]]

** Hacker News                                                             :news:
#+SHORTCUTURL: hn
[[https://news.ycombinator.com/]]

** Development
*** GitHub
[[https://github.com]]
Code hosting platform

* Email
** Gmail                                                           :email:google:
[[https://mail.google.com]]
Personal email
```

## Netscape Bookmark Format

The HTML format follows the Netscape Bookmark file specification used by Firefox and Chrome:

```html
<!DOCTYPE NETSCAPE-Bookmark-file-1>
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3 ADD_DATE="1234567890" LAST_MODIFIED="1234567890">Folder Name</H3>
    <DL><p>
        <DT><A HREF="https://example.com" ADD_DATE="1234567890" LAST_MODIFIED="1234567890" TAGS="tag1,tag2" SHORTCUTURL="keyword">Bookmark Title</A>
    </DL><p>
</DL><p>
```

## Technical Details

### Metadata Preservation

orgmarks preserves the following metadata during conversion:

- **Titles**: Bookmark and folder names
- **URLs**: Full bookmark URLs with query parameters
- **Tags**: Multiple tags per bookmark (Firefox format)
- **Shortcuts**: Keyword shortcuts for quick access (Firefox/Chrome)
- **Timestamps**: ADD_DATE and LAST_MODIFIED from HTML (but not written to Org files - see below)
- **Descriptions**: Additional text associated with bookmarks
- **Hierarchy**: Nested folder structure of any depth

**Note on timestamps**: Timestamps are not written to Org files, I thought it was too messy/cluttered and I don't think anyone cares about this anyway when it comes to web bookmarks. When converting Org back to HTML, the current time is used for ADD_DATE and LAST_MODIFIED. If you need to preserve exact timestamps, don't use this program. I was considering adding another command line option for timestamp preservation, or... you could add it!

### Special Handling

- **Firefox `place:` URLs**: These dynamic query URLs are skipped (Firefox regenerates them)
- **Icons**: ICON and ICON_URI data is ignored (browsers will regenerate favicons anyway)
- **HTML entities**: Special characters are properly escaped/unescaped
- **Empty folders**: Preserved in both formats

### Deduplication

When using `--deduplicate`, orgmarks:

1. Walks the bookmark tree depth-first
2. Tracks encountered URLs
3. Keeps only the first occurrence of each URL
4. Removes all subsequent duplicates

## Building

### Requirements

- Go 1.16 or later
- Dependencies (automatically downloaded):
  - `golang.org/x/net/html` for HTML parsing

### Build Commands

```bash
# Standard build
go build

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o orgmarks-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o orgmarks-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o orgmarks-darwin-arm64
GOOS=windows GOARCH=amd64 go build -o orgmarks-windows-amd64.exe

# Run tests
go test ./...

# Run tests with verbose output
go test ./... -v
```

## Development

### Project Structure

```
orgmarks/
├── main.go                     # CLI entry point
├── models/
│   ├── bookmark.go             # Data structures
│   ├── fixtures.go             # Test fixtures
│   └── deduplicate.go          # Deduplication logic
├── parser/
│   ├── html.go                 # HTML parser
│   ├── html_test.go            # HTML parser tests
│   ├── org.go                  # Org-mode parser
│   └── org_test.go             # Org-mode parser tests
├── converter/
│   ├── html_to_org.go          # HTML → Org converter
│   ├── org_to_html.go          # Org → HTML converter
│   └── converter_test.go       # Integration tests
├── test_bookmarks.org          # Test fixture
└── README.md
```

### Running Tests

The test suite includes:

- **Unit tests**: Parser edge cases (empty files, special characters, deep nesting, long tag lists)
- **Integration tests**: Full conversion pipelines (HTML → Org → HTML round-trips)
- **Fixtures**: Real browser bookmark exports for testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./parser/...
go test ./converter/...

# Run with coverage
go test ./... -cover
```

## Contributing

Contributions welcome! Possible areas for improvement:

- Additional bookmark format support (JSON?)
- Performance optimizations for large bookmark collections
- An option for full timestamp preservation

## License

MIT
