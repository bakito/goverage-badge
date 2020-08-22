package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bakito/goverage-shield/pkg/shield"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"golang.org/x/tools/cover"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	outFile = flag.String("o", "", "go coverage file input")
	inFile  = flag.String("i", "", "shield json output file")
	label   = flag.String("l", "", "label for the shield")
)

func main() {
	flag.Parse()

	// Usage information when no arguments.
	if flag.NFlag() == 0 && flag.NArg() == 0 {
		flag.Usage()
	}

	if *inFile == "" || *outFile == "" {
		os.Exit(2)
	}

	profiles, err := cover.ParseProfiles(*inFile)

	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile(*outFile)

	if err != nil {
		panic(err)
	}
	shield := &shield.Shield{}
	_ = json.Unmarshal(b, shield)

	var out *bufio.Writer
	fd, err := os.Create(*outFile)
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	out = bufio.NewWriter(fd)
	defer out.Flush()

	var total, covered int64
	for _, profile := range profiles {
		fn := profile.FileName
		file, err := findFile(fn)
		if err != nil {
			panic(err)
		}
		funcs, err := findFuncs(file)
		if err != nil {
			panic(err)
		}
		// Now match up functions and profile blocks.
		for _, f := range funcs {
			c, t := f.coverage(profile)
			total += t
			covered += c
		}
	}
	coverage := 100.0 * float64(covered) / float64(total)
	shield.Update(*label,coverage)

	b, err = json.MarshalIndent(shield, "", "  ")
	if err != nil {
		panic(err)
	}
	out.Write(b)
}

// findFuncs parses the file and returns a slice of FuncExtent descriptors.
func findFuncs(name string) ([]*FuncExtent, error) {
	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, name, nil, 0)
	if err != nil {
		return nil, err
	}
	visitor := &FuncVisitor{
		fset:    fset,
		name:    name,
		astFile: parsedFile,
	}
	ast.Walk(visitor, visitor.astFile)
	return visitor.funcs, nil
}

// Visit implements the ast.Visitor interface.
func (v *FuncVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		start := v.fset.Position(n.Pos())
		end := v.fset.Position(n.End())
		fe := &FuncExtent{
			name:      n.Name.Name,
			startLine: start.Line,
			startCol:  start.Column,
			endLine:   end.Line,
			endCol:    end.Column,
		}
		v.funcs = append(v.funcs, fe)
	}
	return v
}

// FuncVisitor implements the visitor that builds the function position list for a file.
type FuncVisitor struct {
	fset    *token.FileSet
	name    string // Name of file.
	astFile *ast.File
	funcs   []*FuncExtent
}

// FuncExtent describes a function's extent in the source by file and position.
type FuncExtent struct {
	name      string
	startLine int
	startCol  int
	endLine   int
	endCol    int
}

// coverage returns the fraction of the statements in the function that were covered, as a numerator and denominator.
func (f *FuncExtent) coverage(profile *cover.Profile) (num, den int64) {
	// We could avoid making this n^2 overall by doing a single scan and annotating the functions,
	// but the sizes of the data structures is never very large and the scan is almost instantaneous.
	var covered, total int64
	// The blocks are sorted, so we can stop counting as soon as we reach the end of the relevant block.
	for _, b := range profile.Blocks {
		if b.StartLine > f.endLine || (b.StartLine == f.endLine && b.StartCol >= f.endCol) {
			// Past the end of the function.
			break
		}
		if b.EndLine < f.startLine || (b.EndLine == f.startLine && b.EndCol <= f.startCol) {
			// Before the beginning of the function
			continue
		}
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}
	if total == 0 {
		total = 1 // Avoid zero denominator.
	}
	return covered, total
}

// findFile finds the location of the named file in GOROOT, GOPATH etc.
func findFile(file string) (string, error) {
	dir, file := filepath.Split(file)
	pkg, err := build.Import(dir, ".", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("can't find %q: %v", file, err)
	}
	return filepath.Join(pkg.Dir, file), nil
}
