package main

import (
	"fmt"
	"log"
	"os"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Service struct {
	Entries []*Entry `@@*`
}

type Entry struct {
	Key		string `@Ident "="`
	Value	string `@Ident`
}

var (
	serviceLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"Ident", `[a-zA-Z0-9/_\-.:,]+`},
		{"Punct", `=`},
		{"whitespace", `\s+`},
	})
	parser = participle.MustBuild[Service](
		participle.Lexer(serviceLexer),
	)
)

func CheckServiceValidity(service *Service) error {
	if service.Entries[0].Key != "Service" {
		return fmt.Errorf("line 1: no service specified")
	}
	if service.Entries[1].Key != "Command" {
		return fmt.Errorf("line 2: no command specified")
	}
	if service.Entries[2].Key != "RunLevel" {
		return fmt.Errorf("line 3: no run level specified")
	}

	return nil
}

func ServiceParser(fileName string) (Service, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return Service{}, fmt.Errorf("error opening file: %w", err)
	}
	content := string(file)

	service, err := parser.ParseString("", content)
	if err != nil {
		return Service{},fmt.Errorf("error parsing file: %w", err)
	}

	err = CheckServiceValidity(service)
	if err != nil {
		return Service{},fmt.Errorf("error parsing file: %w", err)
	}

	return *service, nil
}
