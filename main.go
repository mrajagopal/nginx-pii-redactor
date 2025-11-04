package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	crossplane "github.com/nginxinc/nginx-go-crossplane"
)

type PIIPatterns struct {
	Name        string
	Pattern     *regexp.Regexp
	Replacement string
}

func redactor(block *crossplane.Directive) {
	PIIPatterns := []PIIPatterns{
		{
			Name:        "ip_address",
			Pattern:     regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`),
			Replacement: "[REDACTED_IP]",
		},
		{
			Name:        "domain_names",
			Pattern:     regexp.MustCompile(`\b[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*\.(com|org|net|edu|gov|mil|int|local)\b`),
			Replacement: "[REDACTED_DOMAIN]",
		},
		{
			Name:        "file_paths",
			Pattern:     regexp.MustCompile(`(/[a-zA-Z0-9._-]+)+\.(pem|key|crt|cert|p12|pfx|jks)(\s|$)`),
			Replacement: "[REDACTED_CERT_PATH]",
		},
	}

	fmt.Printf("Directive %s\n", block.Directive)
	if block.Directive != "" {
		for _, piiPattern := range PIIPatterns {
			for i, arg := range block.Args {
				fmt.Printf("Original Arg: %s\n", arg)
				redactedArg := piiPattern.Pattern.ReplaceAllString(arg, piiPattern.Replacement)
				block.Args[i] = redactedArg
				fmt.Printf("Redacted Arg: %s\n", redactedArg)
			}
		}
	}
}

func printDirectiveBlocks(i *crossplane.Directive, depth int) {
	for _, j := range i.Block {
		indent := ""
		for x := 0; x < depth; x++ {
			indent += "  "
		}
		redactor(j)
		fmt.Printf("%sSubBlock: %s\n", indent, j)
		printDirectiveBlocks(j, depth+1)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-crossplane-parser <path-to-go-file>")
		os.Exit(1)
	}
	path := os.Args[1]

	payload, err := crossplane.Parse(path, &crossplane.ParseOptions{
		SingleFile:         true,
		StopParsingOnError: false,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, i := range payload.Config[0].Parsed {
		redactor(i)
		fmt.Printf("Directive %s (%d)\n", i.Directive, i.Line)
		printDirectiveBlocks(i, 0)
	}

	s, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.WriteFile("test.json", s, 0644)
	parserWriteOptions := &crossplane.BuildOptions{
		Tabs:   true,
		Header: true,
		Indent: 4,
	}

	redactedFile, err := os.Create("nginx.redacted.conf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer redactedFile.Close()
	crossplane.Build(redactedFile, payload.Config[0], parserWriteOptions)
}
