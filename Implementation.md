 Implementation
----------------
The Philomath interpreter is organized as a re-entrant multi-pass compiler;
that is, it may process different parts of the source, tokens, or other
intermediaries out of order with each part moving through a pass independently.
Additionally, it should be possible to introduce a new code section to an earlier
pass while processing code sections in a later pass.

In general, parts are partially processed until they run into a missing
dependency and then the compiler continues on with the next part until can
fulfill the dependency.

### Compiler Passes and Sub-passes

1. Lexical Analysis - Transform source characters to tokens
2. Syntax Analysis - Transform tokens into an AST and dependency graph
3. Semantic Analysis - Transform an AST into metadata and control flow graphs
   1. Declaration uniqueness - Verify each declaration name is unique in its scope
   2. Declaration ordering - Verify that variables are declared before use
   3. Type inference - Infer the type of each untyped declaration
   4. Type checking - Check that expressions evaluate into the correct type
   5. ... - There are probably some sub-passses missing here
4. Bytecode Generation - Transform a control flow graph into bytecode
   1. (Optional) CFG optimizations - Data-flow, recursion, loop, and etc optimizations
   2. SSA conversion - Convert to single static assignment form
   3. (Optional) SSA optimizations - Optimiztions which benefit strongly from SSA
   4. Intermediate bytecode - Generate a platform-agnostic bytecode representation
5. (Optional) Interpretation - Run bytecode, and possible re-enter an earlier pass
6. (Optional) Code generation - Transform bytecode to platform-specific assembly
