package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/uncommented/priv8/parser"
)

// CLI flags
var (
	versionFlag bool
	targetFile  string
	outputFile  string
	dryRun      bool
	verbose     bool
)

func main() {
	// Parse command-line flags
	flag.BoolVar(&versionFlag, "version", false, "Print version information")
	flag.StringVar(&targetFile, "file", "", "Target file to analyze")
	flag.StringVar(&outputFile, "output", "", "Output file for sanitized content (default: adds .sanitized suffix)")
	flag.BoolVar(&dryRun, "dry-run", false, "Only detect issues without modifying files")
	flag.BoolVar(&verbose, "verbose", false, "Show detailed processing information")

	flag.Parse()

	// Print version and exit if requested
	if versionFlag {
		printVersion()
		os.Exit(0)
	}

	// Check if a file was specified
	if targetFile == "" {
		color.Red("Error: No target file specified. Use --file flag.")
		flag.Usage()
		os.Exit(1)
	}

	// Check if file exists
	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		color.Red("Error: File '%s' does not exist", targetFile)
		os.Exit(1)
	}

	// Determine output file if not specified
	if outputFile == "" && !dryRun {
		outputFile = targetFile + ".sanitized"
	}

	// Process the file
	if err := processFile(targetFile, outputFile); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}

	color.Green("Processing complete!")
}

// processFile analyzes a bash script for sensitive data and optionally sanitizes it
func processFile(filePath string, outputPath string) error {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Create bash parser
	bashParser, err := parser.NewBashParser()
	if err != nil {
		return fmt.Errorf("failed to create bash parser: %w", err)
	}

	// Parse the bash script
	_, err = bashParser.Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse bash script: %w", err)
	}

	// Write sanitized content to output file
	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write sanitized file: %w", err)
	}

	color.Green("Sanitized file written to %s", outputPath)
	return nil
}
