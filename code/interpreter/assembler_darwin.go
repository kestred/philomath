package interpreter

import "github.com/kestred/philomath/code/utils"

// +build darwin

func formatAssembly(label, source string) string {
  utils.NotImplemented("inline assembly on OS X / Darwin")
  return ""
}
