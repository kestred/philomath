 Axioms
--------
See [The Zen of Python](https://www.python.org/dev/peps/pep-0020/)

 Decision Register
-------------------
List of current language design decisions, with reasoning and concerns.
Decisions are not necessarily in a particular order, except where earlier
decisions heavily impact later decisions.

### [0] Code modules/namespaces
Decisions:
 * Modules and files are completely unrelated
 * Modules provide a namespace which doesn't conflict with the global, or other
   module namespaces
 * A name can be accessed in a module via "." syntax (eg. `module.name`)
 * A module's symbols can be included wholesale in the current scope via
   a `using mymodule` declaration (also allows aliasing and selective imports)
 * A using declaration at file scope doesn't propogate included names when
   that file is loaded in another file, unless declared "propogate" (or "reexport"?)
 * A using declaration at module scope doesn't propogate included names when
   that module is used in another module, unless declared "propogate"
 * The builtin library is provided as the module "builtin" (or similar),
   with submodules for as appropriate (eg. "os", "math", "time", "net", etc)

Reasoning:
 * If modules and files are tied together (either by directories, filenames, etc)
   then it becomes much harder to create easy to distribute libraries with
   beautiful APIs. By not associating these two, it is possible to distribute an
   entire library even as a single source file
 * File inclusion order is less significant if declarations can be out of order
 * Namespacing provides the basic building blocks of making beautiful apis and
   can significantly increase code reuse and identifier readability
 * Seasoned programmers are used to `.` access things like class variables in
   C++, or modules in scripting languages.  It is more beautiful than the `::`
   syntax in C++.  Unless there is some compelling not to use `.`, then it is
   the easiest choice
 * When heavily utilizing symbols from a given namespace, it frequently improves
   code readability to be able to elide the name qualification, so the language
   should support that (for example via "using")
 * If a module is used with "using" in one file, it shouldn't introduce names
   from that module into other files, because the module may introduce some
   functions that aren't used and have separate local implementations
   (for example, it may have a function inverseSqrt which is fast while the
    same function is defined in the current project implemented to be precise)
 * Similarly, names introduced from "using" in other files can make code harder
   to understand, because it is not obvious where a given name is defined
 * An exception to the previous two arguments is: when a module is used
   pervasively throughout the entire project (like a builtin or support library),
   it should be possible to explicitly declare that the symbols should be propogated
 * If a module is used with "using" in one module, using _that_ module shouldn't
   introduce names from the included module. In most cases "using" would be used
   to improve readability in the source code of a module, but the used modules
   are not intended to be exposed in the "public api" of the using module
 * The exception to the previous argument is: when a module is extending,
   wrapping, or borrowing the functionality of another module, it should be
   possible to explicitly declare that the symbols should be propogated
 * The builtin library (if there is one), should not automatically pollute a
   new program's namespace, because a user may want to introduce their own
   alternative or extension to a builtin library (for example, a custom "math"
   library).  Alternatively, it should be possible to explicitly exclude the
   whole or parts of the builtin library, but that option can increase complexity
   when creating a binary w/o the builtin library, and when working with multiple
   files or 3rd party modules.  On the other hand, it is fairly easy to include
   the statement `using builtin propogate` at the top of a single source file
   (probably the file which contains `main ()`)

### [1] Boolean semantics
Decisions:
 * Booleans in a struct are treated as 1-bit bit-fields

Reasoning:
 * Allows for effortless use of flags
 * Easy to remember / intuitive to reason about
 * Can be easier to save memory compared to built-in 8-bit bools

Concerns:
 * Accessing booleans in structs almost always requires masking

### [2] Truthy and Falsy
Decisions:
 * Built-in integer types are not truthy/falsy
 * Null pointers are falsy, other points are truthy

Reasoning:
 * Reduces programming errors from changed types, non-zero "error" conditions, etc
 * Pointer return types don't have the same issues with error return values, so
   there isn't any reason to prevent checking them for truthiness
 * Additionally, implicit bool to pointer conversion allows for more legible
   built-in support for "option types"

Concerns:
 * Increases programming friction when dealing with built-in int types
 * Maybe support defining truthiness for a type (boolean syntactic sugar)?

### [3] Implementation of built-in types
Decisions:
 * The default integer size is machine word size
 * The default floating-point size is machine word size
 * Arrays and strings are (length, pointer) pairs and not zero terminated
 * Strings/text are utf8
 * The byte type doesn't support arithmetic operators,
   and may be the only type which supports bitwise operators

Reasoning:
 * Machine size calculations are typically fastest on the architecture
 * Machine size prevents proliferation of types like uintptr, size_t, etc
 * Arrays with lengths are faster and safer to operate on (at the expense
 * of slightly higher memory usage)
 * Using a byte type for raw binary is more explicit than an integer type

### [4] Names of built-in types
Decisions:
 * bool
 * byte
 * int  i8  i16 i32 i64
 * uint u8  u16 u32 u64
 * real r32 r64
 * text rune

Reasoning:
 * Concise names are less to type
 * Built-in types are easy to learn and remember
 * Integer types i8, u16, etc are less ambiguous than char, short, etc
 * Real is explicitly not "float" to prevent seasoned programmers from
   accidently assuming 32-bit (see decision [3])
 * Real is more intuitive to students than real
 * Real is a little shorter than real
 * Text is explicitly not "string" to prevent seasoned programmers from
   accidently assuming null-termination or 8-bit characters
 * Text is more intuitive to students than string
 * Text is a little shorter than string
 * Rune is explicitly not "char" to prevent seasoned programmers from
   accidently assuming an 8-bit integer
 * Names are lower case because these should be most familiar
   for int, bool, etc to seasoned programmers and there is no
   compelling reason to change those names
 * Types are all lower case because lower-case characters catch less
   attention and types are typically less important to a reader
   (users generally prefer that they be can elided and inferred)

Concerns:
 * Using an unusual name for string is very uncommon and may be
   hard to switch to from other languages using "string"
 * Char is more intuitive to students than rune

### [5] Line comments
Decisions:
 * Line comments start with "#"

Reasoning:
 * Allows for use as an interpreted shell "#!" script
 * Using "#" is shorter than "//" (both should be familiar)
 * Using "#" with directives (as in rust) feels strange to
   me except at the beginning of a line (as in the C preprocessor)

### [6] Block comments
Decisions:
 * Block comments start with "#-" and end with "-#"
 * Block comments are nestable

Reasoning:
 * Feels similar to `//` vs `/* */`, so should be more familiar
 * Looks better than for example `#*` and `*#`
 * Using `#-` instead of `/*`, `"""`, etc because it reduces the
   number of symbols we use as part of the primary language syntax,
   allowing those to be used elsewhere for something else
 * Its easy to turn a line comment into a block comment
 * Nesting allows the programmer to comment out large sections
   even when a smaller section was already block commented
   (for example for documentation, or to temporarily disable a change)

### [7] Variable initialization
Decisions:
 * Variables are initialized to 0 by default
 * Variables can be explicitly declared uninitialized

Reasoning:
 * Zero initialization is the most common use case
 * Reduces programming errors from use of uninitialized variables
 * Expert users can declare a value uninitialized for performance
   within tight loops, although the compiler should also attempt
   to detect and optimize away unnecessary initialization

### [8] Enumerations
Decisions:
 * Enum definitions have the syntax `IDENT :: enum { ITEMS }`
   or `IDENT :: enum TYPE { ITEMS }`
 * An item is an identifier, optionally followed by a value and/or string literal
 * Item values increment automatically from the previously mentioned value
   (the first item having the value 0 by default)
 * Items may be converted into a string, using the provided string or defaulting
   to the text of the identifier
 * Items may be separated by `> IDENT` (or other syntax?), which allows the
   user to check whether a given enum value is in that set of following values
   (up to the next separator)
 * Enums will not be intentionally designed extra functionality to support flag
   values like `0x2`, `0x4`, etc

Reasoning:
 * Enum values very frequently have custom stringified names, and the language
   should support defining these names effortlessly (in addition to expected
   introspection capabilities)
 * Enums are frequently created containing many categories of values, and
   the language should support defining these categories effortlessly (examples
   include "emoji themes", "asset types", and "token categories")
 * Flags are already well supported by the 1-bit booleans w/ struct parameters

Examples:
 * An example with emoticons

   ```
   emoticons :: enum {
       Invalid 0 ""

       > Emotions
       Smiling  ":-)"
       Winking  ";-)"
       Laughing "X-D"
       Frowning ":-("
       Crying   ":'("

       > Things
       Heart       "<3"
       BrokenHeart "</3"
       Rose        "@~)~~~~"
       Bicycle     "(*)/(*)"
   }
   ```

 * An example with game asset types

   ```
   assetType :: enum {
       Invalid 0

       > Sounds
       Attack
       Yelp
       Growl
       Wind
       Chime

       > Bitmaps
       Hero
       Sword
       Monster
   }
    ```

### [9] User-defined definitions
Decisions:
 * User-defined definitions have the syntax `IDENT₁ :: IDENT₂ { DSL_SOURCE }`
 * Parsing the source should be deferred until the IDENT₂ identifier is resolved
 * Allow users to define constant defintions similar to "struct" or "enum"
   declarations with a custom grammar as a way of expressing compile-time DSLs

Reasoning:
 * The recommend syntax is context-free
 * Allows third-party libraries to be as easy to use as native features
 * Code generation libraries do not add file-system dependencies
 * Supports fast type-safe embedded XML, HTML, etc (see X-expressions)
 * Supports mathematics as code (see "formula" example)
 * Supports Backus-Naur, or Lex/Yacc, or etc without generating another source file

Examples:
 * It should be possible to define "formula", such that when used
   elsewhere in the source, the following sentences:

   ```
   meanSqError :: formula { 0.5 Σ[o in |outputs|] (targetsₒ - outputsₒ)² }
   meanSqError :: formula { 0.5 Sum[o in |outputs|] Square[targets[o] - outputs[o]] }
   ```

   are syntactically valid, checkable for correct semantics, with behavior similar to example python:

   ```
   def meanSqError(outputs, targets):
      return 0.5 * sum([(targets[o] - outputs[o])**2 for o in range(len(outputs))])
   ```

### [10] User defined context
Decisions:
 * Allow users to define a "context" (and functions within that context) that is
   implicitly passed to each function call
 * Variables defined in a context should be allocatable on the heap, so that if
   a dynamically loaded library is reloaded, state pointers are not invalidated
 * A context can be re-assigned before or when calling a function to get a fresh
   or modified context for that library

Reasoning:
 * Global variables are almost universally known to be bad, but most languages
   make it really easy to define a global variable and very difficult to pass
   around a set of shared variables for an execution instance or library
 * User defined contexts allow global state to be captured in little to no effort

Examples:
 * Four ways of defining functions that expect a context

   ```
   # Some file
   game_context :: context {
       debugLogger  : logger
       screenBuffer : buffer2d
       audioBuffer  : audioloop
   }

   # Another file (or a module)
   using game_context

   runGameStep :: () {
      logTo(debugLogger, "Ran game")
   }

   # A third option
   updateAudio :: () {
       using game_context

       audioBuffer.samples[0] = 1.0
   }

   # A fourth file/option
   ast_context :: context {
       nullNode : *node

       parseFile :: (path: text) -> root { ... }
   }
   ```

### [A] Naming conventions for types
Decisions:
 * Prefer snake_case for types

Reasoning:
 * Custom defined types will appear similar to built-in
   types, allowing 3rd-party developers to write libraries
   and language extensions that feel like built-ins
 * Similarly, prevents a firstclass / SecondClass dichotomy

### [B] Naming conventions for values
Decisions:
 * Prefer camelCase for values (variables, functions, constants, etc)

Reasoning:
 * Its nice to have some distinction between types and other identifiers
 * Frequently, this is accomplished by using PascalCase and snake_case for one
   or the other; this is where user preference comes in, I like PascalTypes and
   snake_values, except that I dislike Int, Float, etc, dislike inverting the
   cases, and dislike builtins being inconsistent

Compromise:
 * Unlike C/C++, Philomath is a context free grammar;
   with a context free grammar, it is easy to add syntax highlighting which
   correctly highlights any part of the code the programmer deems relevant
 * Ergo, it isn't as necessary to distinguish between types and variables by
   naming convention (ie. naming convention can be used to annotate other info)

### [C] Naming conventions for module names
Decisions:
 * Prefer lowercase without underscores for modules

Reasoning:
 * I personally feel it looks more beautiful as a prefix,
   although I'm on the fence about underscores

   ```
   somemodule.type_name
   somemodule.functionName()
   ```

### [D] Naming conventions for acronyms in identifiers
Decisions:
 * Don't uppercase acronyms in mixedCase or PascalCase names (ie. use Http over HTTP)
 * Don't mix schemes (eg. don't use XMLHttpRequest)

Reasoning:
 * Consistency is helpful to speed code comprehension and guess future names
 * All-caps acronyms are hard to read when 2+ acronyms are adjacent in an identifier
