package repl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/dedisuryadi/bilang/evaluator"
	"github.com/dedisuryadi/bilang/lexer"
	"github.com/dedisuryadi/bilang/object"
	"github.com/dedisuryadi/bilang/parser"
)

const PROMPT = "bilang >>"

func Start(in io.Reader, out io.Writer) {
	var (
		env      = object.NewEnvironment()
		script   = evaluator.NewScript()
		evaluate = func(input string) {
			l := lexer.New(input)
			p := parser.New(l)
			prog := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParseErrors(out, p.Errors())
				return
			}
			evaluated := script.Eval(prog, env)
			if evaluated != nil {
				_, _ = io.WriteString(out, evaluated.Inspect())
				_, _ = io.WriteString(out, "\n")
			}
		}
	)

	defer script.Free()
	if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		inp, _ := ioutil.ReadAll(in)
		evaluate(string(inp))
		return
	}

	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		evaluate(scanner.Text())
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
