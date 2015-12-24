 Axioms
--------
See [The Zen of Python](https://www.python.org/dev/peps/pep-0020/)

 Decision Register
-------------------
List of current language design decisions, with reasoning and concerns

### [1] Boolean semantics
Decision:
  Booleans in a struct are treated as 1-bit bit-fields

Reasoning:
  Allows for very easy implementation of flags
  Easy to remember / intuitive to reason about
  Can be easier to save memory compared to built-in 8-bit bools

Concerns:
  Accessing booleans in structs almost always requires masking

### [2] Truthy and Falsy
Decision:
  Built-in integer types are not truthy/falsy

Reasoning:
  Reduces programming errors from changed types, non-zero error returns, etc

Concerns:
  Increases programming friction when dealing with built-in int types.
  Maybe support defining truthiness for a type (boolean syntactic sugar)?

### [3] Implementation of built-in types
Decision:
  The default integer size is machine word size
  The default realing-point size is machine word size
  Arrays and strings are (length, pointer) pairs and not zero terminated
  Strings are utf8

Reasoning:
  Machine size calculations are typically fastest on the architecture
  Machine size prevents proliferation of types like uintptr, size_t, etc
  Arrays with lengths are faster and safer to operate on (at the expense
  of slightly higher memory usage)

### [4] Names of built-in types
Decision:
  bool byte
  int  i8  i16 i32 i64
  uint u8  u16 u32 u64
  real r32 r64
  text rune

Reasoning:
  Concise names are less to type
  Built-in types are easy to learn and remember
  Integer types i8, u16, etc are less ambiguous than char, short, etc
  Using "byte" for raw binary is more explicit than an integer type
  Real is explicitly not "real" to prevent seasoned programmers from
  accidently assuming 32-bit (see decision [3])
  Real is more intuitive to students than real
  Real is a little shorter than real
  Text is explicitly not "string" to prevent seasoned programmers from
  accidently assuming null-termination or 8-bit characters
  Text is more intuitive to students than string
  Text is a little shorter than string
  Rune is explicitly not "char" to prevent seasoned programmers from
  accidently assuming an 8-bit integer
  Names are lower case because these should be most familiar
  for int, bool, etc to seasoned programmers and there is no
  compelling reason to change those names
  Types are all lower case because lower-case characters catch less
  attention and types are typically less important to a reader
  (users generally prefer that they be can elided and inferred)

Concerns:
  Using an unusual name for string is very uncommon and may be
  hard to switch to from other languages using "string"
  Char is more intuitive to students than rune

### [5] Line comments
Decision:
  Line comments start with "#"

Reasoning:
  Allows for use as an interpreted shell "#!" script
  Using "#" is shorter than "//" (both should be familiar)
  Using "#" with directives (as in rust) feels strange to
  me except at the beginning of a line (as in the C preprocessor)

### [6] Block comments
Decision:
  Block comments start with "#-" and end with "-#"
  Block comments are nestable

Reasoning:
  Feels similar to // vs /* */, so should be more familiar
  Looks better than for example "#*" and "*#"
  Using `#-` instead of `/*`, `"""`, etc because it reduces the
  number of symbols we use as part of our primary PL syntax,
  allowing those to be used elsewhere for something else
  Its easy to turn a line comment into a block comment
  Nesting allows the programmer to comment out large sections
  even when a smaller section was already block commented
  (for example for documentation, or to temporarily disable a change)

### [7] Variable initialization
Decision:
  Variables are initialized to 0 by default
  Variables can be explicitly declared uninitialized

Reasoning:
  Zero initialization is the most common use case
  Reduces programming errors from use of uninitialized variables
  Expert users can declare a value uninitialized for performance
  within tight loops, although the compiler should also attempt
  to detect and optimize away unnecessary initialization

### [8] Macros, Meta-programming, DSLs, etc
Decision:
  The compiler is self-extensible from within a project's source
  Extensions can generate/modify code at multiple stages possibly
  including source, ast, control flow graph, and/or assembly.

Reasoning:
  Allows third-party libraries to be as easy to use as native features
  Code generation libraries do not add file-system dependencies
  Supports fast type-safe embedded XML, HTML, etc (see X-expressions)
  Supports mathematics as code (see "formula" example)
  Supports Backus-Naur or Lex/Yacc without an external file

Examples:
  It should be possible to define "formula", such that when used
  elsewhere in the source, the following sentences:

      MeanSqError := formula { 0.5 Σ[o in |outputs|] (targetsₒ - outputsₒ)² }
      MeanSqError := formula { 0.5 Sum[o in |outputs|] Square[targets[o] - outputs[o]] }

  are syntactically valid, checkable for correct semantics, and each
  compile to the same (abstract) behavior as the python (3.5+) example:

      def MeanSqError(targets, outputs):
        return 0.5 * sum([(targets[o] - outputs[o])**2 for o in range(len(outputs))])

### [A] Naming conventions for types
Decision:
  Prefer snake_case for types

Reasoning:
  Custom defined types will appear similar to built-in
  types, allowing 3rd-party developers to write libraries
  and language extensions that feel like built-ins
  Similarly, prevents a firstclass / SecondClass dichotomy
  Types are all lowercase because lowercase characters catch less
  attention and types are typically less important to a reader
  (users generally prefer that they be can elided and inferred)

Concerns:
  Maybe mixedCase would look more beautiful

### [B] Naming conventions for functions
Decision:
  Prefer PascalCase for function names

Reasoning:
  Functions start with a capital character to appear more
  noticable in source code.  Functions describe a significant
  portion of most applications behavior and generally should
  be more noticable.

### [C] Naming conventions for module names
Decision:
  Prefer lowercase without underscores for modules

Reasoning:
  I personally feel it looks better as a prefix,
  although I'm on the fence about underscores

      somemodule.type_name
      somemodule.FunctionName()

### [D] Naming conventions for struct/enum members and module globals
Decision:
  Prefer PascalCase for scoped members

Reasoning:
  When accessing struct members, having identifiers begin with
  a capital letter helps distinguish between accessing a type
  from a module vs accessing a member from a struct:

      httprouter.router
      self.Value

  This assumes that modules are expanded with ".", which may change.

  I also intuitively feel this decision will work well, but don't yet have
  another rational explanation

### [E] Naming conventions for function parameters and locals
Decision:
  Prefer mixedCase for local variables

Reasoning:
  I intuitively feel this decision will work well, but don't yet have a rational explanation

### [F] Naming conventions for acronyms in identifiers
Decision:
  Don't uppercase acronyms in mixedCase or PascalCase names (ie. use Http over HTTP)
  Don't mix schemes (eg. don't use XMLHttpRequest)

Reasoning:
  Consistency is helpful to speed code comprehension and guess future names
  All-caps acronyms are hard to read when they are adjacent within an identifier
