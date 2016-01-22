 Phi (.φ) - Philomath
----------------------
Phi is a toy language implementation to learn/play-with programming language
design, compiler architectures, byte code generation, virtual machines and
byte code interpretation, byte-code to assembly conversion, and optimizing
compiler techniques.

There are a few primary use cases that the language is being designed for
(primarily focusing on my educational benefit), but ideally it would be suitable
for most computing from embedded programming upto use as an interpreted REPL

 * Machine learning / Neural networks
 * Evolutionary computing / Genetic programming
 * Self hosting compilation
 * Embedded Programming
 * Game development
 * General computing

### Implementation
The bootstrapping compiler is being written in Go.  
See the [compiler overview](docs/compiler.md) for the compiler "architecture".  

### Design
See the [language design notes](docs/notes.md) for the reasoning behind many design decisions.  
See the [grammar](grammar.ebnf) for the currently planned language syntax.  
See the [examples folder](examples) for source code examples of how the language
might be used (I recommend the [feedforward network](examples/feedforward.phi)).  

### License
Primarily for educational use,
but see the [license](license.txt)
