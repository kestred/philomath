package interpreter

import (
	bc "github.com/kestred/philomath/bytecode"
)

// for now, Evaluate will return whatever the result of the last instruction is
func Evaluate(insts []bc.Instruction, consts []bc.Data, totalRegisters bc.Register) bc.Data {
	registers := make([]bc.Data, uint(totalRegisters)+1)

	for _, inst := range insts {
		switch inst.Code {
		case bc.NOOP:
			continue
		case bc.LOAD_CONST:
			registers[inst.Out] = consts[inst.Left]
		case bc.INT64_ADD:
			left := int64(registers[inst.Left])
			right := int64(registers[inst.Right])
			registers[inst.Out] = bc.Data(left + right)
		case bc.INT64_SUBTRACT:
			left := int64(registers[inst.Left])
			right := int64(registers[inst.Right])
			registers[inst.Out] = bc.Data(left - right)
		case bc.INT64_MULTIPLY:
			left := int64(registers[inst.Left])
			right := int64(registers[inst.Right])
			registers[inst.Out] = bc.Data(left * right)
		case bc.INT64_DIVIDE:
			left := int64(registers[inst.Left])
			right := int64(registers[inst.Right])
			registers[inst.Out] = bc.Data(left / right)
		}
	}

	last := insts[len(insts)-1].Out
	return registers[last]
}
