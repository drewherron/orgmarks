// Package main provides the orgmarks CLI tool for converting between
// browser bookmark HTML files and Org-mode format.
//
// orgmarks supports bidirectional conversion with full metadata preservation
// including tags, shortcuts, timestamps, and nested folder structures.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/drewherron/orgmarks/converter"
	"github.com/drewherron/orgmarks/models"
	"github.com/drewherron/orgmarks/parser"
)

const version = "1.0.0"

func main() {
	// Define flags
	inputFile := flag.String("i", "", "Input file (required)")
	outputFile := flag.String("o", "", "Output file (required)")
	deduplicate := flag.Bool("deduplicate", false, "Remove duplicate bookmarks (keep first occurrence)")
	deleteEmpty := flag.Bool("delete-empty", false, "Remove empty folders after processing")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("orgmarks version %s\n", version)
		os.Exit(0)
	}

	// Validate flags
	if *inputFile == "" || *outputFile == "" {
		fmt.Fprintln(os.Stderr, "Usage: orgmarks -i <input-file> -o <output-file>")
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  orgmarks -i bookmarks.html -o bookmarks.org")
		fmt.Fprintln(os.Stderr, "  orgmarks -i bookmarks.org -o bookmarks.html")
		os.Exit(1)
	}

	// Check input file exists
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Input file '%s' does not exist\n", *inputFile)
		os.Exit(1)
	}

	// Detect file formats from extensions
	inputExt := strings.ToLower(filepath.Ext(*inputFile))
	outputExt := strings.ToLower(filepath.Ext(*outputFile))

	// Validate format combination
	if inputExt == outputExt {
		fmt.Fprintln(os.Stderr, "Error: Input and output must have different formats")
		os.Exit(1)
	}

	// Check if output file exists and prompt for confirmation
	if _, err := os.Stat(*outputFile); err == nil {
		fmt.Fprintf(os.Stderr, "orgmarks: overwrite '%s'? ", *outputFile)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError reading input: %v\n", err)
			os.Exit(1)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Fprintln(os.Stderr, "Operation cancelled")
			os.Exit(0)
		}
	}

	// Determine conversion direction
	if (inputExt == ".html" || inputExt == ".htm") && outputExt == ".org" {
		// HTML → Org
		if err := htmlToOrg(*inputFile, *outputFile, *deduplicate, *deleteEmpty); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully converted %s → %s\n", *inputFile, *outputFile)
	} else if inputExt == ".org" && (outputExt == ".html" || outputExt == ".htm") {
		// Org → HTML
		if err := orgToHTML(*inputFile, *outputFile, *deduplicate, *deleteEmpty); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully converted %s → %s\n", *inputFile, *outputFile)
	} else {
		fmt.Fprintln(os.Stderr, "Error: Unsupported file format combination")
		fmt.Fprintln(os.Stderr, "Supported conversions:")
		fmt.Fprintln(os.Stderr, "  .html → .org")
		fmt.Fprintln(os.Stderr, "  .org → .html")
		os.Exit(1)
	}
}

// htmlToOrg converts HTML bookmark file to org-mode
func htmlToOrg(inputFile, outputFile string, deduplicate, deleteEmpty bool) error {
	// Open input file
	in, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer in.Close()

	// Parse HTML
	htmlParser := parser.NewHTMLParser(in)
	root, err := htmlParser.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Apply deduplication if requested
	if deduplicate {
		models.Deduplicate(root)
	}

	// Remove empty folders if requested
	if deleteEmpty {
		models.RemoveEmptyFolders(root)
	}

	// Create output file
	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Convert to org
	if err := converter.ToOrg(root, out); err != nil {
		return fmt.Errorf("failed to convert to org: %w", err)
	}

	return nil
}

// orgToHTML converts org-mode bookmark file to HTML
func orgToHTML(inputFile, outputFile string, deduplicate, deleteEmpty bool) error {
	// Open input file
	in, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer in.Close()

	// Parse org
	orgParser := parser.NewOrgParser(in)
	root, err := orgParser.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse org: %w", err)
	}

	// Apply deduplication if requested
	if deduplicate {
		models.Deduplicate(root)
	}

	// Remove empty folders if requested
	if deleteEmpty {
		models.RemoveEmptyFolders(root)
	}

	// Create output file
	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Convert to HTML
	if err := converter.ToHTML(root, out); err != nil {
		return fmt.Errorf("failed to convert to HTML: %w", err)
	}

	return nil
}
