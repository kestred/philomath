 Philomath (.φ)
----------------
Philomath is a toy language implementation to learn/play-with programming
language design, compiler architectures, byte code generation, virtual machines
and byte code interpretation, byte-code to assembly conversion, and optimizing
compiler techniques.

There are a few primary examples that the language is being designed for
(primarily focusing on my educational benefit), but ideally it would be fairly
suitable as both a low-level systems language and in an interpreted REPL.

 * Machine learning
 * Self hosting compilation
 * Game development
 * General computing

### Implementation
The bootstrapping compiler is being written in Go.  
See the [compiler overview](compiler.md) for the compiler "architecture".  

### Design
See the [language design notes](notes.md) for the reasoning behind each design decision.  
See the [grammar](grammar.ebnf) for the currently planned language syntax.  
See the [examples folder](examples) for source code examples of how the language
mght be used (I recommend the [feedforward network](examples/feedforward.phi)).  

### License
Primarily for educational use,
but see the [license](license.txt)
