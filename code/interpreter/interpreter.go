package interpreter

import bc "github.com/kestred/philomath/code/bytecode"

func Run(prog *bc.Program) bc.Data {
	start := prog.Procedures[prog.Text["start_"]]
	main := prog.Text["main"]
	call := bc.Instruction{Op: bc.CALL, Out: start.AssignRegister(), Left: bc.Register(main)}
	start.Instructions = append(start.Instructions, call)
	return Evaluate(start)
}

// HACK: for now, Evaluate will return whatever the result of the last instruction is
func Evaluate(proc bc.Procedure) bc.Data {
	registers := make([]bc.Data, uint(proc.NextRegister))
	returnRegister := bc.Register(-1)
	count := len(proc.Instructions)
	if count > 0 {
		returnRegister = proc.Instructions[count-1].Out
	}

InstructionLoop:
	for _, inst := range proc.Instructions {
		switch inst.Op {
		case bc.NOOP:
			continue
		case bc.COPY_VALUE:
			registers[inst.Out] = registers[inst.Left]
		case bc.LOAD_CONST:
			registers[inst.Out] = proc.Program.Constants[inst.Left]
		case bc.CALL:
			registers[inst.Out] = Evaluate(proc.Program.Procedures[inst.Left])
		case bc.CALL_ASM:
			asm := proc.Program.Metadata[inst.Left].(*bc.Assembly)
			CallAsm(asm, registers)
		case bc.RETURN:
			returnRegister = inst.Out
			break InstructionLoop

		// signed 64-bit arithmetic
		case bc.I64_ADD:
			left := bc.ToI64(registers[inst.Left])
			right := bc.ToI64(registers[inst.Right])
			registers[inst.Out] = bc.FromI64(left + right)
		case bc.I64_SUBTRACT:
			left := bc.ToI64(registers[inst.Left])
			right := bc.ToI64(registers[inst.Right])
			registers[inst.Out] = bc.FromI64(left - right)
		case bc.I64_MULTIPLY:
			left := bc.ToI64(registers[inst.Left])
			right := bc.ToI64(registers[inst.Right])
			registers[inst.Out] = bc.FromI64(left * right)
		case bc.I64_DIVIDE:
			left := bc.ToI64(registers[inst.Left])
			right := bc.ToI64(registers[inst.Right])
			registers[inst.Out] = bc.FromI64(left / right)

		// unsigned 64-bit arithmetic
		case bc.U64_ADD:
			left := bc.ToU64(registers[inst.Left])
			right := bc.ToU64(registers[inst.Right])
			registers[inst.Out] = bc.FromU64(left + right)
		case bc.U64_SUBTRACT:
			left := bc.ToU64(registers[inst.Left])
			right := bc.ToU64(registers[inst.Right])
			registers[inst.Out] = bc.FromU64(left - right)
		case bc.U64_MULTIPLY:
			left := bc.ToU64(registers[inst.Left])
			right := bc.ToU64(registers[inst.Right])
			registers[inst.Out] = bc.FromU64(left * right)
		case bc.U64_DIVIDE:
			left := bc.ToU64(registers[inst.Left])
			right := bc.ToU64(registers[inst.Right])
			registers[inst.Out] = bc.FromU64(left / right)

		// floating-point 64-bit arithmetic
		case bc.F64_ADD:
			left := bc.ToF64(registers[inst.Left])
			right := bc.ToF64(registers[inst.Right])
			registers[inst.Out] = bc.FromF64(left + right)
		case bc.F64_SUBTRACT:
			left := bc.ToF64(registers[inst.Left])
			right := bc.ToF64(registers[inst.Right])
			registers[inst.Out] = bc.FromF64(left - right)
		case bc.F64_MULTIPLY:
			left := bc.ToF64(registers[inst.Left])
			right := bc.ToF64(registers[inst.Right])
			registers[inst.Out] = bc.FromF64(left * right)
		case bc.F64_DIVIDE:
			left := bc.ToF64(registers[inst.Left])
			right := bc.ToF64(registers[inst.Right])
			registers[inst.Out] = bc.FromF64(left / right)

		// conversions
		case bc.CONVERT_F64_TO_I64:
			value := bc.ToF64(registers[inst.Left])
			registers[inst.Out] = bc.FromI64(int64(value))
		case bc.CONVERT_F64_TO_U64:
			value := bc.ToF64(registers[inst.Left])
			registers[inst.Out] = bc.FromU64(uint64(value))
		case bc.CONVERT_I64_TO_F64:
			value := bc.ToI64(registers[inst.Left])
			registers[inst.Out] = bc.FromF64(float64(value))
		case bc.CONVERT_U64_TO_F64:
			value := bc.ToU64(registers[inst.Left])
			registers[inst.Out] = bc.FromF64(float64(value))

		default:
			panic("TODO: Unhandled opcode")
		}
	}

	if returnRegister >= 0 {
		return registers[returnRegister]
	} else {
		return 0
	}
}
