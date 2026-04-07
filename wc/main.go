// ccwc implements part of the Unix wc command for the Coding Challenges project.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type Config struct {
	Filename  string
	ShowBytes bool
	ShowLines bool
	ShowWords bool
	ShowChars bool
}

type Counts struct {
	Bytes int64
	Lines int64
	Words int64
	Chars int64
}

func main() {
	cfg, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cfg Config) error {
	file, err := os.Open(cfg.Filename)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	counts, err := count(file)
	if err != nil {
		return err
	}

	if err := printCounts(cfg, counts); err != nil {
		return err
	}

	return nil
}

// parseArgs parses command-line arguments into a Config.
func parseArgs(args []string) (Config, error) {
	fs := flag.NewFlagSet("ccwc", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	showBytes := fs.Bool("c", false, "count bytes")
	showLines := fs.Bool("l", false, "count lines")
	showWords := fs.Bool("w", false, "count words")
	showChars := fs.Bool("m", false, "count characters")

	usageString := "usage: ccwc [-l] [-w] [-c] [-m] <file>"

	if err := fs.Parse(args); err != nil {
		return Config{}, fmt.Errorf("ccwc: %v\n%s", err, usageString)
	}

	if fs.NArg() != 1 {
		return Config{}, fmt.Errorf("ccwc: expected exactly one file argument\n%s", usageString)
	}

	cfg := Config{
		Filename:  fs.Arg(0),
		ShowBytes: *showBytes,
		ShowLines: *showLines,
		ShowWords: *showWords,
		ShowChars: *showChars,
	}

	if !cfg.ShowBytes && !cfg.ShowLines && !cfg.ShowWords && !cfg.ShowChars {
		cfg.ShowBytes = true
		cfg.ShowLines = true
		cfg.ShowWords = true
	}

	return cfg, nil
}

// count reads all data from r and returns the computed counts.
func count(r io.Reader) (Counts, error) {
	reader := bufio.NewReader(r)
	var counts Counts
	inWord := false

	for {
		ch, size, err := reader.ReadRune()
		if err == io.EOF {
			return counts, nil
		}
		if err != nil {
			return Counts{}, err
		}

		counts.Bytes += int64(size)
		counts.Chars++

		if ch == '\n' {
			counts.Lines++
		}
		if unicode.IsSpace(ch) {
			inWord = false
			continue
		}
		if !inWord {
			inWord = true
			counts.Words++
		}
	}
}

// printCounts prints the selected counts for the given file.
func printCounts(cfg Config, counts Counts) error {
	parts := make([]string, 0, 4)

	if cfg.ShowBytes {
		parts = append(parts, fmt.Sprintf("Bytes: %8d", counts.Bytes))
	}
	if cfg.ShowLines {
		parts = append(parts, fmt.Sprintf("Lines: %8d", counts.Lines))
	}
	if cfg.ShowWords {
		parts = append(parts, fmt.Sprintf("Words: %8d", counts.Words))
	}
	if cfg.ShowChars {
		parts = append(parts, fmt.Sprintf("Chars: %8d", counts.Chars))
	}

	parts = append(parts, fmt.Sprintf("File: %s", cfg.Filename))

	_, err := fmt.Println(strings.Join(parts, " "))
	return err
}
