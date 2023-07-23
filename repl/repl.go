package main

import (
	"bufio"
	"cardboard/lexer"
	"cardboard/parser"
	"fmt"
	"os"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">>> ")
		scanner.Scan()
		if scanner.Text() == "" {
			return
		}

		input := scanner.Text()
		lex := lexer.CreateLexer(input)
		parser := parser.CreateParser(lex)
		program := parser.ParseCardBoard()

		if checkParserErrors(parser) {
			continue
		}

		fmt.Println(program.String())
	}
}

func checkParserErrors(p *parser.Parser) bool {
	errs := p.GetErrors()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		return true
	}
	return false
}
