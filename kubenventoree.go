package kubenventoree

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type cmdoptions struct {
	AllTests     bool   `help:"Run all queries"`
	Output       string `short:"o" required help:"output file name" type:"path"`
	OutputFormat string `short:"f" default:text help:"format of the output (json, yaml, text)"`
}

var cliOptions cmdoptions

func main() {
	ctx := kong.Parse(&cliOptions)
	fmt.Printf("Command: %s\n", ctx.Command())
	fmt.Printf("Struct: %s %s", cliOptions.Output, cliOptions.OutputFormat)
}
