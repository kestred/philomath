## Use Cases
This doc describes which "not in C" features I attribute to supporting each of
the use cases described in the [readme](../readme.md).

Like most of the docs at this point, this is more for brainstorming than
for use as a reference material.

### Machine Learning

Helpful:
 1. Catch-all for making it easy to translate formal math and logic into code

### Evolutionary Computing

Helpful:
 1. Catch-all for making it easy to dynamic loading and linking
 2. Being interpreted makes it easy to run generated code (compiled lets it run fast)

### Self-hosting

Helpful:
 1. Unicode strings by default
 2. Union types (eg. int|text) and the "empty" type allow succinct tree definitions
 3. More-powerful enums are useful for succinct Tokens, etc, etc, etc
 4. Operator overloading w/ unicode math could be useful for formal logic

Painful:
 1. Language accepts utf8 characters

### Game Programming

Helpful:
 1. Catch-all for making it easy to dynamic loading and linking
 2. Operator overloading is great for Vectors, et. al.
 3. Being interpreted and compiled means you can also write "scripts" for easy
    to iterate and update behavior, but without any of the drawbacks:
      There doesn't need to be a "glue" or "wrapper" layer.
      The scripts can be compiled for better performance after they stablize.

### Embedded Programming

Helpful:
 1. Builtin linting tools (for eg. MISRA)
 2. Unrestricted inline assembly (any instruction set, specified assembler)
 3. Ability to substitue builtin operator behavior to be processor-specific (see 4)
 4. Ability to substitue builtin memory manipulators (eg. malloc implementation)
 5. Compilation to C and/or (optional?) calling-conventional compatiblity
 6. Short names for fixed-width builtin types (i8, u16, f32, etc)

Painful:
 1. Use of platform-width types are conventional and encouraged (int, float)
 2. Using Go for bootstraping reduces portability compared to C
 3. I do embedded programming infrequently, so I'll make bad decisions unintentionally
