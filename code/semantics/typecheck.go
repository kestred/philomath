package semantics

import "github.com/kestred/philomath/code/utils"

// TODO: implement most types of type checking (so far I only check # of arguments in procedure calls)

func CheckTypes(cs *Section) []error {
	utils.Assert(!cs.DidSteps(Step_CheckTypes), "Tried to run type-checking twice on the same code section")
	utils.Assert(cs.DidSteps(Step_InferTypes), "Tried to run type-checking before type inference")
	utils.Assert(cs.DidSteps(Step_ResolveNames), "Tried to run type-checking before name resolution")

	// TODO: move error printing out from parser and into its own module; then add source ranges to AST nodes

	cs.StepsCompleted |= Step_CheckTypes
	return nil
}
