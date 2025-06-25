# Org-mode Bookmark Format Specification

This document describes the Org-mode format used by orgmarks for representing browser bookmarks.

## Overview

orgmarks represents bookmarks in org-mode format using standard Org conventions:
- Headlines represent both folders and bookmarks
- Tags use the standard `:tag1:tag2:` syntax
- Links use the standard `[[URL]]` or `[[URL][title]]` syntax
- Properties use the `#+KEY: value` syntax

## File Structure

### Root Level

The root of the bookmark tree is implicit. All top-level headlines become root-level folders or bookmarks.

```org
* Folder One
* Folder Two
* Root Level Bookmark
[[https://example.com]]
```

### Folders

A folder is any headline that does **not** contain a link in its content.

```org
* Top Level Folder
** Subfolder
*** Deeply Nested Folder
```

Empty folders are valid:

```org
* Empty Folder
* Another Folder
** This one has content
```

### Bookmarks

A bookmark is any headline that **contains** a link in its content.

The bookmark's title comes from the headline text, **not** from the link title.

```org
* My Bookmark Title
[[https://example.com]]
```

Even if the link has a title, the headline is used:

```org
* My Bookmark Title
[[https://example.com][This title is ignored]]
```

### Nesting

Folders can be nested to arbitrary depth using Org headline levels:

```org
* Level 1
** Level 2
*** Level 3
**** Level 4
***** Level 5
****** Level 6
******* Bookmark at Level 7
[[https://deep.example.com]]
```

## Metadata

### Tags

Tags are appended to headlines using standard Org syntax:

```org
* Bookmark Title                                                     :tag1:tag2:
[[https://example.com]]
```

Tags are typically aligned to column 80 for readability:

```org
* Short Title                                                    :tag1:tag2:tag3:
[[https://example.com]]
```

Multiple tags are separated by colons with no spaces:

```org
* Bookmark                                          :programming:javascript:web:
[[https://example.com]]
```

Folders can also have tags, but these are currently not used in HTML conversion:

```org
* Development Folder                                               :work:code:
** Bookmark
[[https://example.com]]
```

### Shortcut URLs (Keywords)

Firefox and Chrome support "keyword" or "shortcut" URLs that allow quick access via the address bar.

In Org format, these are stored as properties using the `#+SHORTCUTURL:` directive:

```org
* Google                                                              :search:
#+SHORTCUTURL: g
[[https://google.com]]
```

The property line must appear:
1. After the headline
2. Before the link
3. On its own line

Multiple properties can be specified (though currently only SHORTCUTURL is used):

```org
* Bookmark
#+SHORTCUTURL: mybookmark
#+CUSTOM: value
[[https://example.com]]
```

### Descriptions

Any text after the link (and after any property lines) is treated as the bookmark's description:

```org
* Gmail                                                           :email:google:
[[https://mail.google.com]]
Personal email account
```

Multi-line descriptions are supported:

```org
* Complex Bookmark
#+SHORTCUTURL: cb
[[https://example.com]]
This is the first line of the description.

This is a second paragraph.

- You can even use
- Lists and other
- Org formatting here
```

Currently, orgmarks treats all content after the link as a single description field.

### Timestamps

Org-mode supports timestamps, but orgmarks currently does not parse or generate them in Org format. Timestamps are only preserved in the internal model when converting from HTML, and are used when converting back to HTML.

When creating new Org bookmarks without timestamps, orgmarks will use the current time when converting to HTML.

## Links

### Simple Link Format

The most common format:

```org
[[URL]]
```

Example:

```org
[[https://example.com]]
```

### Link with Title Format

Org supports link titles, but orgmarks uses the **headline text** as the bookmark title, not the link title:

```org
* Bookmark Title
[[https://example.com][Ignored Title]]
```

In this case, "Bookmark Title" is used as the bookmark title, and "Ignored Title" is discarded.

### URL Schemes

Any valid URL scheme is supported:

```org
[[https://example.com]]
[[http://example.com]]
[[ftp://example.com]]
[[file:///path/to/file]]
```

You can also omit the scheme, and browsers will add `http://` automatically when importing:

```org
[[www.google.com]]
[[example.com]]
```

This makes it easier to type URLs in Emacs without worrying about the protocol.

### Query Parameters

URLs with query parameters are fully supported:

```org
[[https://example.com/page?foo=1&bar=2&baz=3]]
```

Special characters in URLs do not need to be escaped in Org format (they're handled by the HTML parser).

## Complete Example

Here's a complete example demonstrating all features:

```org
* Bookmarks Toolbar
** Quick Access
*** Google                                                            :search:
#+SHORTCUTURL: g
[[https://google.com]]
Search engine

*** Gmail                                                       :email:google:
#+SHORTCUTURL: mail
[[https://mail.google.com]]
Personal email

** Development
*** GitHub                                                        :code:tools:
[[https://github.com]]
Code hosting platform

*** Stack Overflow                                          :programming:help:
[[https://stackoverflow.com]]

** News & Reading
*** Hacker News                                                  :tech:startup:
[[https://news.ycombinator.com]]

*** Reddit                                                       :forum:social:
[[https://reddit.com]]

* Work
** Project Docs
[[https://docs.company.com/project]]
Current project documentation

** Empty Folder for Future Use

* Personal
** Shopping                                                        :ecommerce:
*** Amazon
[[https://amazon.com]]

*** eBay
[[https://ebay.com]]
```

## Conversion Notes

### From HTML to Org

When converting from HTML to Org, orgmarks:

1. Converts folder hierarchy to headline levels
2. Extracts `TAGS` attribute and converts to `:tag:` format
3. Extracts `SHORTCUTURL` attribute and creates `#+SHORTCUTURL:` property
4. Uses the `<A>` tag text as the headline title
5. Preserves timestamps internally (but doesn't write them to Org)
6. Skips Firefox `place:` URLs
7. Ignores ICON data

### From Org to HTML

When converting from Org to HTML, orgmarks:

1. Converts headline levels to nested `<DL>` structure
2. Converts `:tag:` format to `TAGS` attribute (comma-separated)
3. Converts `#+SHORTCUTURL:` property to `SHORTCUTURL` attribute
4. Uses headline text as the `<A>` tag text
5. Generates Unix timestamps for `ADD_DATE` and `LAST_MODIFIED`
6. Uses current time if no timestamps are available
7. Escapes HTML special characters (`&`, `<`, `>`, `"`)

## Special Cases

### Headlines Without Links

If a headline has no link, it's treated as a folder, even if it has other content:

```org
* This Is A Folder
Some text here, but no link means this is a folder
```

This could be useful for organization in Emacs, but tags and descriptions will be ignored by all browsers (I think).

### Bookmarks at Root Level

Bookmarks can exist at the root level (level 1):

```org
* Root Bookmark
[[https://example.com]]

* Root Folder
** Nested Bookmark
[[https://example.com/nested]]
```

### Empty Folders

Folders with no children are preserved:

```org
* Empty Folder One

* Folder With Content
** Child

* Empty Folder Two
```

### Mixed Content

Folders and bookmarks can be mixed at the same level:

```org
* Folder
** Subfolder
*** Bookmark A
[[https://a.example.com]]
** Bookmark B
[[https://b.example.com]]
** Another Subfolder
* Another Folder
```

### Special Characters

Special characters in titles are preserved as-is in Org format:

```org
* Bookmark with [brackets] & "quotes"
[[https://example.com]]

* Bookmark with <angle> brackets
[[https://example.com]]
```

When converting to HTML, these are properly escaped:
- `&` → `&amp;`
- `<` → `&lt;`
- `>` → `&gt;`
- `"` → `&quot;`

## Limitations

### What's Not Supported

1. **Org PROPERTIES drawers**: Standard Org property drawers are not parsed. Use `#+SHORTCUTURL:` instead.

2. **Timestamps**: Org `<2024-01-01>` style timestamps are not parsed or generated.

3. **Multiple links per headline**: Only the first link is recognized. Additional links are treated as description text.

4. **Link descriptions in HTML output**: The `[[URL][description]]` description is not used; headline text is always the bookmark title.

5. **Folder timestamps/metadata**: Folder-level metadata (except title and children) is not preserved when converting to HTML.

## See Also

- [Org Mode Manual](https://orgmode.org/manual/)
- [Netscape Bookmark File Format](https://docs.microsoft.com/en-us/previous-versions/windows/internet-explorer/ie-developer/platform-apis/aa753582(v=vs.85))
- orgmarks README.md for usage examples
