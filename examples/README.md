# Example Bookmark Files

This directory contains example bookmark files demonstrating the orgmarks conversion format.

## Files

### simple.org

A simple Org-mode bookmark file demonstrating:
- Nested folders (Development, News, Tools)
- Bookmarks with tags (`:programming:help:`, `:tech:startup:`)
- Shortcut URLs (`#+SHORTCUTURL: g`)
- Descriptions (text after links)

### simple.html

The HTML equivalent of `simple.org` in Netscape Bookmark format, showing:
- Nested `<DL>` and `<H3>` folder structure
- `<A>` tags with `TAGS` and `SHORTCUTURL` attributes
- Unix timestamps in `ADD_DATE` and `LAST_MODIFIED` attributes

## Usage

Convert the Org file to HTML:

```bash
orgmarks -i simple.org -o output.html
```

Convert the HTML file to Org:

```bash
orgmarks -i simple.html -o output.org
```

Round-trip test (should produce identical structure):

```bash
orgmarks -i simple.org -o temp.html
orgmarks -i temp.html -o round-trip.org
diff simple.org round-trip.org
```
