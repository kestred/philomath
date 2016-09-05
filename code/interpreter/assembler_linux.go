package interpreter

// +build linux

// formatAssembly targeting GNU Assembler w/ intel syntax
func formatAssembly(label, source string) string {
  return fmt.Sprintf(`
  .intel_syntax
  .global %s
  .section .text

  %s:
  %s
  ret
  `, label, label, source)
}
