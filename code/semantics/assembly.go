package semantics

import (
	"github.com/kestred/philomath/code/ast"
	"github.com/kestred/philomath/code/scanner"
	"github.com/kestred/philomath/code/token"
)

func PreprocessAssembly(cs *ast.Section) {
	// TODO: I really shouldn't need to check ALL nodes for this
	for _, node := range cs.Nodes {
		if asm, ok := node.(*ast.AsmBlock); ok {
			detectInputsOutputs(asm)
		}
	}
}

func detectInputsOutputs(asm *ast.AsmBlock) {
	var labels []string

	// TODO: Stop using this scanner
	var scan scanner.Scanner
	scan.Init("asm", []byte(asm.Source), nil)

	// HACK: Good nuff for now but it can't not be buggy
	prevTok := token.INVALID
	prevLit := ""
	prevLine := -1
	prevInst := false
	isMov := false
	offset, tok, lit := scan.Scan()
	for tok != token.END {
		pos := scan.Pos()
		isInst := false
		switch tok {
		case token.IDENT:
			if pos.Line > prevLine {
				if scan.Peek() == token.COLON {
					labels = append(labels, lit)
					asm.Inputs = removeBindings(asm.Inputs, lit)
					asm.Outputs = removeBindings(asm.Outputs, lit)
				} else {
					isInst = true
					if len(lit) >= 3 {
						isMov = (lit[0:3] == "mov")
					}
				}
			} else if !(prevTok == token.OPERATOR && prevLit == "%") {
				switch lit {
				case "near", "far", "byte", "word", "dword", "qword", "ptr":
					break
				case "NEAR", "FAR", "BYTE", "WORD", "DWORD", "QWORD", "PTR":
					break
				default:
					if stringsInclude(labels, lit) {
						break
					}

					if prevInst && isMov && scan.Peek() == token.COMMA {
						asm.Outputs = append(asm.Outputs, ast.AsmBinding{ast.Ident(lit), offset})
					} else {
						asm.Inputs = append(asm.Inputs, ast.AsmBinding{ast.Ident(lit), offset})
					}
				}
			}
		}

		prevInst = isInst
		prevLine = pos.Line
		prevLit = lit
		prevTok = tok
		offset, tok, lit = scan.Scan()
	}
}

func removeBindings(bindings []ast.AsmBinding, name string) []ast.AsmBinding {
	n := 0
	for n < len(bindings) {
		if bindings[n].Name.Literal == name {
			bindings[n] = bindings[len(bindings)-1]
			bindings = bindings[:len(bindings)-1]
		} else {
			n++
		}
	}

	if len(bindings) > 0 {
		return bindings
	} else {
		return nil
	}
}

func stringsInclude(strings []string, el string) bool {
	for _, s := range strings {
		if s == el {
			return true
		}
	}
	return false
}
