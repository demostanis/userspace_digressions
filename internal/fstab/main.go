package fstab

import (
	"fmt"
	"os"
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

func MountMP(mp *MountPoint) {
	flags, remainingOptions := parseMountOptions(mp.Options)
	unix.Mount(mp.Source, mp.Target, mp.Type, flags, remainingOptions)
}

func MountDefaultsMPs() {
	err = os.MkdirAll("/dev", 0755)
	if err != nil {
		fmt.Println("/dev already exist")
	}
	err = os.MkdirAll("/proc", 0755)
	if err != nil {
		fmt.Println("/proc already exist")
	}
	err = os.MkdirAll("/dev/pts", 0755)
	if err != nil {
		fmt.Println("/dev/pts already exist")
	}
	err = os.MkdirAll("/dev/shm", 0755)
	if err != nil {
		fmt.Println("/dev/shm already exist")
	}

	devtmpfs	:= MountPoint{"devtmpfs", "/dev", "devtmpfs", "exec,nosuid,mode=0755,size=2M", "", ""}
	tmpfs		:= MountPoint{"tmpfs", "/dev", "tmpfs", "exec,nosuid,mode=0755,size=2M", "", ""}
	proc		:= MountPoint{"proc", "/proc", "proc", "noexec,nosuid,nodev", "", ""}
	devpts		:= MountPoint{"devpts", "/dev/pts", "devpts", "gid=5,mode=0620,noexec,nosuid", "", ""}
	shm			:= MountPoint{"shm", "/dev/shm", "tmpfs", "nodev,nosuid,noexec", "", ""}

	MountMP(&devtmpfs)
	MountMP(&tmpfs)
	MountMP(&proc)
	MountMP(&devpts)
	MountMP(&shm)
}

func FstabParser(file string) {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading  file:", err)
		return
	}
	content_str := string(content)

	fstab, err := parser.ParseString("", content_str)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	MountDefaultsMPs()
	for _, mp := range fstab.MountPoints {
		MountMP(mp)
	}
}
