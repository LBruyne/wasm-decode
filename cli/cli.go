package cli

import (
	"flag"
	"fmt"
	"github.com/LBruyne/wasm-decode/decode"
	"os"
)

func main() {
	dumpFlag := flag.Bool("d", false, "dump")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Println("Usage: wasm-decode [-d] -filename")
		os.Exit(1)
	}

	module, err := decode.DecodeFile(flag.Args()[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if *dumpFlag {
		dump(module)
	}
}
