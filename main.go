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

// stringSlice is a custom flag type that allows multiple values
type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	// Define flags
	var inputFiles stringSlice
	flag.Var(&inputFiles, "i", "Input file (can be specified multiple times for merging)")
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
	if len(inputFiles) == 0 || *outputFile == "" {
		fmt.Fprintln(os.Stderr, "Usage: orgmarks -i <input-file> [-i <input-file2> ...] -o <output-file>")
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  orgmarks -i bookmarks.html -o bookmarks.org")
		fmt.Fprintln(os.Stderr, "  orgmarks -i bookmarks.org -o bookmarks.html")
		fmt.Fprintln(os.Stderr, "  orgmarks -i file1.org -i file2.org -o merged.org    # Merge multiple files")
		fmt.Fprintln(os.Stderr, "  orgmarks -i organized.org -i new.html -o final.org  # Merge different formats")
		os.Exit(1)
	}

	// Check all input files exist
	for _, inputFile := range inputFiles {
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Input file '%s' does not exist\n", inputFile)
			os.Exit(1)
		}
	}

	outputExt := strings.ToLower(filepath.Ext(*outputFile))

	// Handle multiple input files (merge mode)
	if len(inputFiles) > 1 {
		// Merge mode - output must be .org
		if outputExt != ".org" {
			fmt.Fprintln(os.Stderr, "Error: When merging multiple files, output must be .org format")
			os.Exit(1)
		}

		if err := mergeFiles(inputFiles, *outputFile, *deduplicate, *deleteEmpty); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully merged %d files → %s\n", len(inputFiles), *outputFile)
		os.Exit(0)
	}

	// Single input file - regular conversion mode
	inputFile := inputFiles[0]
	inputExt := strings.ToLower(filepath.Ext(inputFile))

	// Validate format combination for single file conversion
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
		if err := htmlToOrg(inputFile, *outputFile, *deduplicate, *deleteEmpty); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully converted %s → %s\n", inputFile, *outputFile)
	} else if inputExt == ".org" && (outputExt == ".html" || outputExt == ".htm") {
		// Org → HTML
		if err := orgToHTML(inputFile, *outputFile, *deduplicate, *deleteEmpty); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully converted %s → %s\n", inputFile, *outputFile)
	} else {
		fmt.Fprintln(os.Stderr, "Error: Unsupported file format combination")
		fmt.Fprintln(os.Stderr, "Supported conversions:")
		fmt.Fprintln(os.Stderr, "  .html → .org")
		fmt.Fprintln(os.Stderr, "  .org → .html")
		os.Exit(1)
	}
}

// mergeFiles merges multiple bookmark files into a single org file
func mergeFiles(inputFiles []string, outputFile string, deduplicate, deleteEmpty bool) error {
	// Parse the first file
	root, err := parseFile(inputFiles[0])
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", inputFiles[0], err)
	}

	// Merge remaining files
	for i := 1; i < len(inputFiles); i++ {
		nextTree, err := parseFile(inputFiles[i])
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", inputFiles[i], err)
		}
		root = models.MergeFolders(root, nextTree)
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

// parseFile parses a bookmark file (either HTML or org) and returns the root folder
func parseFile(filename string) (*models.Folder, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if ext == ".html" || ext == ".htm" {
		htmlParser := parser.NewHTMLParser(file)
		return htmlParser.Parse()
	} else if ext == ".org" {
		orgParser := parser.NewOrgParser(file)
		return orgParser.Parse()
	} else {
		return nil, fmt.Errorf("unsupported file format: %s", ext)
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
