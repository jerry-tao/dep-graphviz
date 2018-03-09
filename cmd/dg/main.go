package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/sdboyer/gps/pkgtree"
	"hash/fnv"
	"os"
	"strings"
)

var (
	std = flag.Bool("s", false, "include stdlib, default false")
)

func main() {
	flag.Parse()
	root := flag.Arg(0)
	ptree, err := pkgtree.ListPackages(root, root)
	if err != nil {
		fmt.Println("Error: %v", err)
		os.Exit(1)
	}
	g := new(graphviz).New()
	for _, v := range ptree.Packages {
		if v.Err == nil {
			g.createNode(v.P.Name, "", v.P.Imports)
			for _, t := range v.P.Imports {
				if *std {
					g.createNode(t, "", nil)
				} else {
					if !doIsStdLib(t) {
						g.createNode(t, "", nil)
					}
				}
			}
		}
	}
	str := g.output()
	fmt.Println(str.String())

}

func doIsStdLib(path string) bool {
	i := strings.Index(path, "/")
	if i < 0 {
		i = len(path)
	}

	return !strings.Contains(path[:i], ".")
}

type graphviz struct {
	ps []*gvnode
	b  bytes.Buffer
	h  map[string]uint32
}

type gvnode struct {
	project  string
	version  string
	children []string
}

func (g graphviz) New() *graphviz {
	ga := &graphviz{
		ps: []*gvnode{},
		h:  make(map[string]uint32),
	}
	return ga
}

func (g graphviz) output() bytes.Buffer {
	g.b.WriteString("digraph {\n\tnode [shape=box];")

	for _, gvp := range g.ps {
		// Create node string
		g.b.WriteString(fmt.Sprintf("\n\t%d [label=\"%s\"];", gvp.hash(), gvp.label()))
	}

	// Store relations to avoid duplication
	rels := make(map[string]bool)

	// Create relations
	for _, dp := range g.ps {
		for _, bsc := range dp.children {
			for pr, hsh := range g.h {
				if isPathPrefix(bsc, pr) {
					r := fmt.Sprintf("\n\t%d -> %d", g.h[dp.project], hsh)

					if _, ex := rels[r]; !ex {
						g.b.WriteString(r + ";")
						rels[r] = true
					}

				}
			}
		}
	}

	g.b.WriteString("\n}")
	return g.b
}

func (g *graphviz) createNode(project, version string, children []string) {
	pr := &gvnode{
		project:  project,
		version:  version,
		children: children,
	}

	g.h[pr.project] = pr.hash()
	g.ps = append(g.ps, pr)
}

func (dp gvnode) hash() uint32 {
	h := fnv.New32a()
	h.Write([]byte(dp.project))
	return h.Sum32()
}

func (dp gvnode) label() string {
	label := []string{dp.project}

	if dp.version != "" {
		label = append(label, dp.version)
	}

	return strings.Join(label, "\\n")
}

// isPathPrefix ensures that the literal string prefix is a path tree match and
// guards against possibilities like this:
//
// github.com/sdboyer/foo
// github.com/sdboyer/foobar/baz
//
// Verify that prefix is path match and either the input is the same length as
// the match (in which case we know they're equal), or that the next character
// is a "/". (Import paths are defined to always use "/", not the OS-specific
// path separator.)
func isPathPrefix(path, pre string) bool {
	pathlen, prflen := len(path), len(pre)
	if pathlen < prflen || path[0:prflen] != pre {
		return false
	}

	return prflen == pathlen || strings.Index(path[prflen:], "/") == 0
}
