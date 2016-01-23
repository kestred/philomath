### Motivation
Global variables are almost universally known to be bad, but most languages
make it really easy to define a global variable and very difficult to pass
around a set of shared variables for an execution instance or library

### Feature
The language supports defining a "context" (and procedures using that context)
that is implicitly passed to procedures that use it.
Values in a context can be re-assigned before or when calling a procedure to get
a fresh or modified context for that procedure invocation.

### Machine Implementation
**Problem:** Passing many context w/o passing many pointers  
**Problem:** Remembering a specific context in a nested procedure when called
from a procedure not using that context  

> Both of these problems can be solved by using something like a
"context table".  To pass the context to any procedure, you pass a pointer
to the context table; then, a procedure within the originally compiled executable
knows at compile time what the offset for a given context is in the table.

**Problem:** Passing a context to a linked libary / shared object

> Context offsets, at least those used when compiling a library, should be
stored in the .DATA or .BSS sections of the binary.  At link time, this value
must be overridden with the actual offset of the context in the executable's
context table.  The library in the looks up the offset when it needs to access
values from the context.

**Problem:** Creating/using a context within a linked library

> If a library uses a context that was not defined by the linking executable at
compile time, the linker can assign values past the end of the compiled offsets,
and set the offset values to refer to those locations.  Optionally the last
value in the context table (or some other location) could be a pointer to an
extended or dynamic context table. The size of a context may need to be
specified in the BSS/DATA of the library.

**Problem:** Optimizing "leaf" procedure calls

> If a procedure does not use any contexts and does not make any procedure calls,
it would be possible to elide storing the context argument in a register or on
the stack.  If the procedureal call semantics has the context table in a known
location, then that register or stack location does not need to be set.
(You could potentially use that location for something else, but for shared
objects / dynamically loaded code I'd avoid it because it complicates the ABI).

**Problem:** Using a context with threads
> A thread creation routine may need to allocate and initialize, possibly by
copying the current values, a new context table and new contexts.
> If context support multi-threading (and all features must in this world), then
it cannot use a predefined global location for the context table.  Instead, the
context table COULD be stored in thread-local storage (using segment registers).
In the general case, if the context table is always stored on the heap and
passed as a pointer, this wouldn't be a problem; however, consider the case
of calling into C code and calling back a procedure which requires a context,
then the table would need to be at a knowable (not necessarily static) location.
