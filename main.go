package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file> <resource_type.name>\n", os.Args[0])
		os.Exit(1)
	}
	filePath := os.Args[1]
	full := os.Args[2]

	// Split "type.name" into type and name
	dot := strings.IndexRune(full, '.')
	if dot == -1 {
		fmt.Fprintln(os.Stderr, "Please specify the resource as type.name")
		os.Exit(1)
	}
	resType, resName := full[:dot], full[dot+1:]

	/* -------- Parse HCL and obtain the AST -------- */
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(filePath)
	exitIfDiag(diags)

	syntaxFile, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		fmt.Fprintln(os.Stderr, "Failed to convert to HCL syntax body")
		os.Exit(1)
	}

	/* -------- Locate the target resource block -------- */
	var target *hclsyntax.Block
	for _, b := range syntaxFile.Blocks {
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

	/* -------- Determine the line range -------- */
	startLn := target.Range().Start.Line
	endLn := target.Range().End.Line

	/* -------- Read only the required lines -------- */
	lines, err := readLines(filePath, startLn, endLn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	/* -------- Marshal to JSON and print -------- */
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

// Exit immediately if the diagnostics contain errors.
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
		if ln > endLn { // stop once we've read past endLn
			break
		}
		ln++
	}
	return out, sc.Err()
}
