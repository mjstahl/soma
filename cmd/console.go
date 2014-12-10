package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/socialmachines/soma/file"
	"github.com/socialmachines/soma/parse"
	"github.com/socialmachines/soma/rt"
)

var ConsoleUsage = `Usage:
    soma console

    Start a Social Machines read-eval-print loop
    allowing the user to evaluate expressions.

Example:
    $ soma console
    >>> + True not -> False.
	=== True
    >>> True not
    === False

The Social Machines console supports commands
that are evaluated differently than Social
Machines expressions.

The commands are:
   :exit 	exits the runtime
   :info	displays runtime information
   :objs	list of objects and behaviors
   :quit	exits the runtime
`

func StartConsole(ver string) {
	scope := rt.NewScope(nil)

	pwd, _ := os.Getwd()
	pd := file.ProjDirFrom(pwd)
	if pd == "" {
		fmt.Printf("Social Machines (v%s). Type ':exit' to exit.\n", ver)

		startREPL(scope)
	} else {
		ps, err := LoadProjectDir(pd, scope)
		if err != nil {
			displayConsoleError("failed to load project directory", err)
		}

		fmt.Printf("%s // Social Machines (v%s). Type ':exit' to exit.\n", filepath.Base(pd), ver)
		startREPL(ps)
	}
}

func startREPL(s *rt.Scope) {
	if len(rt.RT.Peers) > 0 {
		ln, port := rt.StartListening(10810)
		rt.RT.Port = port

		go http.Serve(ln, nil)
	}

	for {
		fmt.Printf(">>> ")
		reader := bufio.NewReader(os.Stdin)
		raw, _ := reader.ReadString('\n')

		input := strings.TrimSpace(raw)
		if input != "" {
			if isConsoleCmd(input) {
				evalConsoleCmd(input)
			} else {
				evaluateInput(input, s)
			}
		}
	}
}

func isConsoleCmd(input string) bool {
	if input[0] == ':' {
		return true
	}

	return false
}

func evalConsoleCmd(input string) {
	lower := strings.ToLower(input)
	switch lower {
	case ":exit":
		os.Exit(0)
	case ":info":
		printRuntimeInfo()
		printMemoryInfo()
	case ":objs":
		printObjects()
	}
}

func evaluateInput(input string, scope *rt.Scope) {
	expr, err := parse.ParseExpr(input)
	if err != nil {
		fmt.Println("!!!", err)
	} else {
		if len(expr) > 0 {
			fmt.Println("===", expr[len(expr)-1].Visit(scope))
		}
	}
}

func printRuntimeInfo() {
	fmt.Println(" + Runtime")
	fmt.Printf(" |   ID: 0x%x\n", rt.RT.ID>>31)

	named := len(rt.RT.Globals.Values)
	heap := len(rt.RT.Heap.Values) + len(rt.RT.Peers)
	goroutines := runtime.NumGoroutine()
	fmt.Printf(" |   Objects (Named/Lang/Sys): %d/%d/%d\n", named, heap, goroutines)

	fmt.Printf(" |   Cores Used: %d\n", runtime.NumCPU())
}

func printMemoryInfo() {
	mem := new(runtime.MemStats)
	runtime.ReadMemStats(mem)

	fmt.Println(" + Memory")
	avg := mem.PauseTotalNs / uint64(mem.NumGC) / 1.0e3
	fmt.Printf(" |   Avg. GC Pause: %d \u03BCs\n", avg)
	fmt.Printf(" |   In Use: %d KB\n", mem.Alloc/1024)
	fmt.Printf(" |   Total Allocated: %d KB\n", mem.TotalAlloc/1024)
}

func printObjects() {
	for index, name := range rt.RT.Globals.Order {
		fmt.Printf(" + %s\n", name)

		oid := rt.RT.Globals.Values[index]
		obj := rt.RT.Heap.Values[oid]
		for behave := range obj.(*rt.Object).Behaviors {
			fmt.Printf(" |   %s\n", behave)
		}
	}
}

func displayConsoleError(msg string, err error) {
	fmt.Printf("soma console: %s: %s\n", msg, err)
	os.Exit(1)
}
