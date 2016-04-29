package interpreter

import (
	"unsafe"

	bc "github.com/kestred/philomath/code/bytecode"
	"github.com/kestred/philomath/code/utils"
)

func Run(prog *bc.Program) []byte {
	start := prog.Procedures[prog.Text["start_"]]
	out := bc.Rg(start.AssignLocation(), start.PrevResult.Typ)
	proc := prog.Procedures[prog.Text["main"]]
	call := bc.Inst(bc.CALL, bc.Proc(proc, out))
	start.Instructions = append(start.Instructions, call)
	return Evaluate(start)
}

// HACK: for now, Evaluate will return whatever the result of the last instruction is
func Evaluate(proc bc.Procedure) []byte {
	registers := make([][]byte, uint(proc.NextFree))
	returnRegister := bc.Rg(-1, bc.None)
	count := len(proc.Instructions)
	if count > 0 {
		inst := proc.Instructions[count-1]
		switch args := inst.Args.(type) {
		case bc.UnaryArgs:
			returnRegister = args.Out
		case bc.BinaryArgs:
			returnRegister = args.Out
		case bc.ConstantArgs:
			returnRegister = args.Out
		}
	}

InstructionLoop:
	for _, inst := range proc.Instructions {
		switch inst.Op {
		case bc.NOOP:
			continue
		case bc.COPY:
			args := inst.Args.(bc.UnaryArgs)
			registers[args.Out.Loc] = registers[args.In.Loc]
		case bc.LOAD:
			switch args := inst.Args.(type) {
			case bc.ConstantArgs:
				if args.Ptr {
					ptr := &proc.Program.Data[args.Name][0]
					raw := uint64(uintptr(unsafe.Pointer(ptr)))
					registers[args.Out.Loc] = bc.Pack(raw)
				} else {
					registers[args.Out.Loc] = proc.Program.Data[args.Name]
				}
			case bc.UnaryArgs:
				panic("TODO: Load data from pointer")
				// registers[args.Out] = proc.Program.Constants[args.Left]
			}
		case bc.CALL:
			args := inst.Args.(bc.ProcedureArgs)
			registers[args.Out.Loc] = Evaluate(args.Proc)
		case bc.CALL_ASM:
			asm := inst.Args.(*bc.AssemblyArgs)
			CallAsm(asm, registers)
		case bc.RETURN:
			returnRegister = inst.Args.(bc.NullaryArgs).Rg
			break InstructionLoop

		// signed 64-bit arithmetic
		case bc.ADD:
			args := inst.Args.(bc.BinaryArgs)
			switch args.Left.Typ {
			case bc.Int64:
				var left, right int64
				unpackRegister(inst, registers, args.Left.Loc, &left)
				unpackRegister(inst, registers, args.Right.Loc, &right)
				registers[args.Out.Loc] = bc.Pack(left + right)
			case bc.Uint64:
				var left, right uint64
				unpackRegister(inst, registers, args.Left.Loc, &left)
				unpackRegister(inst, registers, args.Right.Loc, &right)
				registers[args.Out.Loc] = bc.Pack(left + right)
			case bc.Float64:
				var left, right float64
				unpackRegister(inst, registers, args.Left.Loc, &left)
				unpackRegister(inst, registers, args.Right.Loc, &right)
				registers[args.Out.Loc] = bc.Pack(left + right)
			}
		case bc.SUBTRACT:
			args := inst.Args.(bc.BinaryArgs)
			switch args.Left.Typ {
			case bc.Int64:
				var left, right int64
				unpackRegister(inst, registers, args.Left.Loc, &left)
				unpackRegister(inst, registers, args.Right.Loc, &right)
				registers[args.Out.Loc] = bc.Pack(left - right)
			case bc.Uint64:
				var left, right uint64
				unpackRegister(inst, registers, args.Left.Loc, &left)
				unpackRegister(inst, registers, args.Right.Loc, &right)
				registers[args.Out.Loc] = bc.Pack(left - right)
			case bc.Float64:
				var left, right float64
				unpackRegister(inst, registers, args.Left.Loc, &left)
				unpackRegister(inst, registers, args.Right.Loc, &right)
				registers[args.Out.Loc] = bc.Pack(left - right)
			}
		case bc.MULTIPLY:
			args := inst.Args.(bc.BinaryArgs)
			switch args.Left.Typ {
			case bc.Int64:
				var left, right int64
				unpackRegister(inst, registers, args.Left.Loc, &left)
				unpackRegister(inst, registers, args.Right.Loc, &right)
				registers[args.Out.Loc] = bc.Pack(left * right)
			case bc.Uint64:
				var left, right uint64
				unpackRegister(inst, registers, args.Left.Loc, &left)
				unpackRegister(inst, registers, args.Right.Loc, &right)
				registers[args.Out.Loc] = bc.Pack(left * right)
			case bc.Float64:
				var left, right float64
				unpackRegister(inst, registers, args.Left.Loc, &left)
				unpackRegister(inst, registers, args.Right.Loc, &right)
				registers[args.Out.Loc] = bc.Pack(left * right)
			}
		case bc.DIVIDE:
			args := inst.Args.(bc.BinaryArgs)
			switch args.Left.Typ {
			case bc.Int64:
				var left, right int64
				unpackRegister(inst, registers, args.Left.Loc, &left)
				unpackRegister(inst, registers, args.Right.Loc, &right)
				registers[args.Out.Loc] = bc.Pack(left / right)
			case bc.Uint64:
				var left, right uint64
				unpackRegister(inst, registers, args.Left.Loc, &left)
				unpackRegister(inst, registers, args.Right.Loc, &right)
				registers[args.Out.Loc] = bc.Pack(left / right)
			case bc.Float64:
				var left, right float64
				unpackRegister(inst, registers, args.Left.Loc, &left)
				unpackRegister(inst, registers, args.Right.Loc, &right)
				registers[args.Out.Loc] = bc.Pack(left / right)
			}

		// conversions
		case bc.CAST_I64:
			args := inst.Args.(bc.UnaryArgs)
			switch args.In.Typ {
			case bc.Uint64:
				var in uint64
				unpackRegister(inst, registers, args.In.Loc, &in)
				registers[args.Out.Loc] = bc.Pack(int64(in))
			case bc.Float64:
				var in float64
				unpackRegister(inst, registers, args.In.Loc, &in)
				registers[args.Out.Loc] = bc.Pack(int64(in))
			}
		case bc.CAST_U64:
			args := inst.Args.(bc.UnaryArgs)
			switch args.In.Typ {
			case bc.Int64:
				var in int64
				unpackRegister(inst, registers, args.In.Loc, &in)
				registers[args.Out.Loc] = bc.Pack(uint64(in))
			case bc.Float64:
				var in float64
				unpackRegister(inst, registers, args.In.Loc, &in)
				registers[args.Out.Loc] = bc.Pack(uint64(in))
			}
		case bc.CAST_F64:
			args := inst.Args.(bc.UnaryArgs)
			switch args.In.Typ {
			case bc.Int64:
				var in int64
				unpackRegister(inst, registers, args.In.Loc, &in)
				registers[args.Out.Loc] = bc.Pack(float64(in))
			case bc.Uint64:
				var in uint64
				unpackRegister(inst, registers, args.In.Loc, &in)
				registers[args.Out.Loc] = bc.Pack(float64(in))
			}

		default:
			panic("TODO: Unhandled opcode")
		}
	}

	if returnRegister.Loc >= 0 {
		return registers[returnRegister.Loc]
	} else {
		return nil
	}
}

func unpackRegister(inst bc.Instruction, registers [][]byte, loc bc.Location, ptr interface{}) {
	err := bc.Unpack(registers[loc], ptr)
	utils.Assert(err == nil, `%v (at %v)`, err, inst)
}
