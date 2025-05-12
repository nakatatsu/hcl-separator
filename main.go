package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file> (<resource_type.name>|<line_number>)\n", os.Args[0])
		os.Exit(1)
	}
	filePath := os.Args[1]
	identifier := os.Args[2]

	/* -------- Parse HCL and obtain the AST -------- */
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(filePath)
	exitIfDiag(diags)

	syntaxFile, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		fmt.Fprintln(os.Stderr, "Failed to convert to HCL syntax body")
		os.Exit(1)
	}

	// Decide mode by trying to parse identifier as int
	if line, err := strconv.Atoi(identifier); err == nil {
		// -------- Line mode --------
		handleLineMode(filePath, syntaxFile.Blocks, line)
	} else {
		// -------- Resource mode --------
		handleResourceMode(filePath, syntaxFile.Blocks, identifier)
	}
}

// handleResourceMode locates a resource by "type.name"
func handleResourceMode(filePath string, blocks hclsyntax.Blocks, full string) {
	dot := strings.IndexRune(full, '.')
	if dot == -1 {
		fmt.Fprintln(os.Stderr, "Please specify the resource as type.name")
		os.Exit(1)
	}
	resType, resName := full[:dot], full[dot+1:]

	var target *hclsyntax.Block
	for _, b := range blocks {
		if b.Type == "resource" && len(b.Labels) == 2 &&
			b.Labels[0] == resType && b.Labels[1] == resName {
			target = b
			break
		}
	}
	if target == nil {
		fmt.Fprintf(os.Stderr, "Resource %s not found\n", full)
		os.Exit(1)
	}

	outputBlockJSON(filePath, target)
}

// handleLineMode locates the innermost block that contains the given line
func handleLineMode(filePath string, blocks hclsyntax.Blocks, line int) {
	target := findBlockByLine(blocks, line)
	if target == nil {
		fmt.Fprintf(os.Stderr, "No block found containing line %d\n", line)
		os.Exit(1)
	}
	outputBlockJSON(filePath, target)
}

// findBlockByLine returns the deepest block whose range includes the line.
func findBlockByLine(blocks hclsyntax.Blocks, line int) *hclsyntax.Block {
	for _, b := range blocks {
		r := b.Range()
		if line >= r.Start.Line && line <= r.End.Line {
			// try to go deeper first
			if inner := findBlockByLine(b.Body.Blocks, line); inner != nil {
				return inner
			}
			return b
		}
	}
	return nil
}

// outputBlockJSON marshals the block info and prints it.
func outputBlockJSON(filePath string, block *hclsyntax.Block) {
	startLn := block.Range().Start.Line
	endLn := block.Range().End.Line

	lines, err := readLines(filePath, startLn, endLn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	result := struct {
		StartLine int    `json:"start_line"`
		EndLine   int    `json:"end_line"`
		Content   string `json:"content"`
	}{
		StartLine: startLn,
		EndLine:   endLn,
		Content:   strings.Join(lines, "\n"),
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}

/* ---------- Helper functions ---------- */

func exitIfDiag(diags hcl.Diagnostics) {
	if diags == nil || !diags.HasErrors() {
		return
	}
	for _, d := range diags {
		fmt.Fprintln(os.Stderr, d.Error())
	}
	os.Exit(1)
}

// readLines returns the lines in [startLn, endLn] (1-based, inclusive).
func readLines(path string, startLn, endLn int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []string
	sc := bufio.NewScanner(f)
	ln := 1
	for sc.Scan() {
		if ln >= startLn && ln <= endLn {
			out = append(out, sc.Text())
		}
		if ln > endLn {
			break
		}
		ln++
	}
	return out, sc.Err()
}
