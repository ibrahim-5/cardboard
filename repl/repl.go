package repl

import (
	"bufio"
	"cardboard/eval"
	"cardboard/lexer"
	"cardboard/object"
	"cardboard/parser"
	"fmt"
	"os"
	"strings"
)

func StartREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	env := object.CreateEnvironment()

	fmt.Println("Cardboard v1.0! type :q to quit REPL.")

	for {
		fmt.Print(">>> ")
		scanner.Scan()

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			return
		} else if input == ":q" {
			fmt.Println("Ending REPL.")
			os.Exit(0)
		}

		lex := lexer.CreateLexer(input)
		parser := parser.CreateParser(lex)
		program := parser.ParseCardBoard()

		if checkParserErrors(parser) {
			continue
		}

		evaluatedProgram := eval.Eval(program, env)
		if evaluatedProgram.Type() == object.SHOW_OBJ || evaluatedProgram.Type() == object.ERROR_OBJ {
			fmt.Println("> " + evaluatedProgram.Inspect())
		}
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
