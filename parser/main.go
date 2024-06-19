package main

import (
    "fmt"
    "os"
    "github.com/alecthomas/participle/v2"
    "github.com/alecthomas/participle/v2/lexer"
)

type FSTAB struct {
	MountPoints []*MountPoint `@@*`
}

type MountPoint struct {
	FileSystem	string `@Ident`
	MountPoint	string `@Ident`
	Type		string `@Ident`
	Options		string `@Ident`
	Dump		string `@Ident?`
	Pass		string `@Ident?`
}

var (
	fstabLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"Comment", `(?:#)[^\n]*\n?`},
		{"Ident", `[a-zA-Z0-9/_\-.:,=]+`},
		{"Whitespace", `[ \t\n\r]+`},
	})
	parser = participle.MustBuild[FSTAB](
		participle.Lexer(fstabLexer),
		participle.Elide("Comment", "Whitespace"),
		participle.UseLookahead(2),
	)
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
        return
	}

	file := os.Args[1]

    content, err := os.ReadFile(file)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }
    content_str := string(content)

    fstab, err := parser.ParseString("", content_str)
    if err != nil {
        fmt.Println("Error parsing file:", err)
        return
    }

    for _, mp := range fstab.MountPoints {
        fmt.Printf("FileSystem: %s\nMountPoint: %s\nType: %s\nOptions: %s\nDump: %s\nPass: %s\n\n",
            mp.FileSystem, mp.MountPoint, mp.Type, mp.Options, mp.Dump, mp.Pass)
    }
}
