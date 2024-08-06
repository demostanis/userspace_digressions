package services

import (
	"fmt"
	"os"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type ServiceFile struct {
	Entries []*Entry `@@*`
}

type Entry struct {
	Key		string `@Ident "="`
	Value	string `@Ident`
}

type Service struct {
	Service		string
	Command		string
	RunLevel	string
}

var (
	serviceLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"Ident", `[a-zA-Z0-9\/_\-.:, <>]+`},
		{"Punct", `=`},
		{"whitespace", `\s+`},
	})
	parser = participle.MustBuild[ServiceFile](
		participle.Lexer(serviceLexer),
	)
)

func CheckServiceValidity(service *ServiceFile) (Service, error) {
	if service.Entries[0].Key != "Service" {
		return Service{}, fmt.Errorf("line 1: no service specified")
	}
	if service.Entries[1].Key != "Command" {
		return Service{}, fmt.Errorf("line 2: no command specified")
	}
	if service.Entries[2].Key != "RunLevel" {
		return Service{}, fmt.Errorf("line 3: no run level specified")
	}
	
	return Service{
		Service: service.Entries[0].Value,
		Command: service.Entries[1].Value,
		RunLevel: service.Entries[2].Value,
	}, nil
}

func ServiceParser(fileName string) (Service, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return Service{}, fmt.Errorf("error opening file: %w", err)
	}
	content := string(file)

	serviceFile, err := parser.ParseString("", content)
	if err != nil {
		return Service{},fmt.Errorf("error parsing file: %w", err)
	}

	service, err := CheckServiceValidity(serviceFile)
	if err != nil {
		return Service{},fmt.Errorf("error parsing file: %w", err)
	}

	return service, nil
}
