package easy

import (
	"fmt"
	"os"
	"sort"
)

type Cmd struct {
	Help   string
	Action func(args []string)
}

// SubCmd selects the given subcommand named by the first command line
// argument next to the command name and runs the command by passing
// the rest of the command line arguments. The subcommand name "help"
// is reserved for the built-in subcommand for listing available
// subcommands and should not be used (m["help"] is silently
// over-ridden if a user-specified "help" subcommand is given).
func SubCmd(m map[string]Cmd) {
	prog := os.Args[0]
	m["help"] = Cmd{"list available subcommands", func(_ []string) { printUsage(prog, m) }}
	if len(os.Args) <= 1 {
		printUsage(prog, m)
		os.Exit(1)
	}
	cmd := os.Args[1]
	sub, ok := m[cmd]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unrecognized subcommand: %q\n", cmd)
		os.Exit(1)
	}
	sub.Action(os.Args[2:])
}

func printUsage(prog string, m map[string]Cmd) {
	fmt.Fprintf(os.Stderr, "Available subcommands of %s:\n", prog)
	names := []string{}
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintf(os.Stderr, "  %s: %s\n", name, m[name].Help)
	}
}
