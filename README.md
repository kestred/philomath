 Phi (.φ) - Philomath
----------------------
Phi is a toy programming language that was built to teach myself more about
compiler architecture, byte-code generation, virtual machines, assembly, and
optimizing-compiler techniques.

The code architecture is briefly described in the [implementation](Implementation.md) doc.

### What does it do?

Phi was intended as both an interpreted and compiled language, but only the interperter was implemented.

See the [example](examples/example1.phi) which demonstrates (most of?) the functionality that was implemented:

  - Basic conditional control flow
  - Basic procedure/function invocation
  - Inline GNU Assembly code with JIT execution
      - __(a.k.a.)__ a humorous substitute for a missing standard library 😁;  
        perhapy the only __cool__ (☃️) language feature.

### Public Domain

To the extent possible under law, I waive all copyright and related or neighboring rights to Philomath.

![Public Domain](http://i.creativecommons.org/p/zero/1.0/88x31.png)
