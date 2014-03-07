package easy

import (
	"fmt"
	"os"
	"sort"
)

type Cmd struct {
	Brief  string // One line brief description of what the command does.
	Detail string // Detailed description of the command.
	Action func(args []string)
}

// SubCmd selects the given subcommand named by the first command line
// argument next to the command name and runs the command by passing
// the rest of the command line arguments. When no subcommand is
// given, a list of available subcommands are printed to stderr. There
// is also a built-in "help" command that either lists the available
// subcommands or describes a subcommand in more detail. The "help"
// command can be over-ridden supplying a "help" command in m.
func SubCmd(m map[string]Cmd) {
	prog := os.Args[0]

	// Supply the default "help" command if there is not one.
	_, ok := m["help"]
	if !ok {
		m["help"] = newHelp(prog, m)
	}

	// Show list of available commands when there is no argument.
	if len(os.Args) <= 1 {
		printUsage(prog, m)
		os.Exit(1)
	}

	// Find and run the command.
	cmd := os.Args[1]
	sub, ok := m[cmd]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unrecognized subcommand: %q\nRun without arguments to see available subcommands.\n", cmd)
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
		fmt.Fprintf(os.Stderr, "  %s: %s\n", name, m[name].Brief)
	}
}

func describeCommand(cmd string, sub Cmd) {
	fmt.Printf("%s: %s\n\n%s\n", cmd, sub.Brief, sub.Detail)
}

func newHelp(prog string, m map[string]Cmd) Cmd {
	return Cmd{
		"lists available subcommands or describes a subcommand in detail",
		`Without an argument, "help" lists all available subcommands. Otherwise, it describes the subcommand specified by the first argument.`,
		func(args []string) {
			if len(args) == 0 {
				printUsage(prog, m)
			} else {
				cmd := args[0]
				sub, ok := m[cmd]
				if !ok {
					fmt.Fprintf(os.Stderr, "help: unrecognized subcommand: %q; run without arguments to see available subcommands.\n", cmd)
					os.Exit(1)
				} else {
					describeCommand(cmd, sub)
				}
			}
		},
	}
}
