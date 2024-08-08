package fstab

import (
	"fmt"
	"os"
	"os/exec"
	"golang.org/x/sys/unix"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

type FSTAB struct {
	MountPoints []*MountPoint `@@*`
}

type MountPoint struct {
	Source	string  `@Ident`
	Target	string  `@Ident`
	Type	string  `@Ident`
	Options	string  `@Ident`
	Dump	string  `@Ident?`
	Pass	string  `@Ident?`
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

func parseMountOptions(options string) (uintptr, string) {
	var flags uintptr
	var remainingOptions []string

	opts := strings.Split(options, ",")
	for _, opt := range opts {
		switch opt {
		case "ro":
			flags |= unix.MS_RDONLY
		case "nosuid":
			flags |= unix.MS_NOSUID
		case "nodev":
			flags |= unix.MS_NODEV
		case "noexec":
			flags |= unix.MS_NOEXEC
		case "exec":
			// Ignore, implicit if not noexec
		case "sync":
			flags |= unix.MS_SYNCHRONOUS
		case "dirsync":
			flags |= unix.MS_DIRSYNC
		case "remount":
			flags |= unix.MS_REMOUNT
		case "mand":
			flags |= unix.MS_MANDLOCK
		case "noatime":
			flags |= unix.MS_NOATIME
		case "nodiratime":
			flags |= unix.MS_NODIRATIME
		case "relatime":
			flags |= unix.MS_RELATIME
		default:
			remainingOptions = append(remainingOptions, opt)
		}
	}
	return flags, strings.Join(remainingOptions, ",")
}

func MountSwap(mp *MountPoint) error {
	cmd := exec.Command("swapon", mp.Source)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(stdout))
	}

	return nil
}

func Mount(mp *MountPoint) error {
	if mp.Type == "swap" {
		return MountSwap(mp)
	}

	flags, remainingOptions := parseMountOptions(mp.Options)

	err := os.MkdirAll(mp.Target, 0755)
	if (err != nil) {
		return err
	}

	err = unix.Mount(mp.Source, mp.Target, mp.Type, flags, remainingOptions)

	if err != nil  {
		return err
	}

	return nil
}

func FstabParser(file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	content_str := string(content)

	fstab, err := parser.ParseString("", content_str)
	if err != nil {
		return fmt.Errorf("error parsing file: %w", err)
	}

	for _, mp := range fstab.MountPoints {
		err = Mount(mp)
		if err != nil {
			return fmt.Errorf("error mounting mountpoint: %w", err)
		}
	}

	return nil
}
