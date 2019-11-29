package main

import (
	"log"
	"os"
	parser2 "pl0compiler/parser"
	"pl0compiler/scanner"
	"pl0compiler/token"
)

type Config struct {
	filepath string
}

func main() {
	config := Config{}
	args := os.Args
	if len(args) > 1 {
		config.filepath = args[1]
	}
	log.Printf("SourceCode FilePath: %s\n", config.filepath)
	parser := parser2.Parser{}
	parser.Init(config.filepath)
	parser.Start()
	parser.Sout()
	parser.Interpret()
}

func testparser(filepath string) {
	scanner := scanner.Scanner{}
	scanner.Init(filepath)
	for {
		tok, lit, num := scanner.Scan()
		if tok == token.EOFSYM {
			break
		}
		if lit == "" {
			lit = "NULL"
		}
		if tok != token.BADTOKEN {
			log.Printf("<TYPE: %d,\tVALUE: %s,\tNUM: %d >\t类型说明：%s", tok, lit, num, tok.String())
		} else {
			log.Printf("Error: <TYPE: %d,\tVALUE: %s,\tNUM: %d >\t类型说明：%s", tok, lit, num, tok.String())
		}
	}
}
