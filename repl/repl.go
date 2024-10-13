package repl

import (
	"bufio"
	"fmt"
	"io"
	"morty/evaluator"
	"morty/lexer"
	"morty/object"
	"morty/parser"
	"morty/read"
)

func Start(file *bufio.Reader, out io.Writer) {
	env := object.NewEnvironment()
	io.WriteString(out, "RESULTS:\n")

	for {

		line, isPrefix, err := read.Readln(file)
		if err != nil {
			return
		}
		if isPrefix {
			fmt.Printf("Line length exceeded")
			return
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "  parsing errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
