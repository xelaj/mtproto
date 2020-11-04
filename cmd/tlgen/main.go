package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xelaj/mtproto/cmd/tlgen/gen"
	"github.com/xelaj/mtproto/cmd/tlgen/typelang"
)

const helpMsg = `generate-tl-files
usage: tlgen input_file.tl output_dir/

THIS TOOL IS USING ONLY FOR AUTOMATIC CODE
GENERATION, DO NOT GENERATE FILES BY HAND!

No, seriously. Don't. go generate is amazing. You
are amazing too, but lesser 😏
`

func main() {
	if len(os.Args) != 3 {
		fmt.Println(helpMsg)
		return
	}

	if err := run(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(tlfile, outdir string) error {
	b, err := ioutil.ReadFile(tlfile)
	if err != nil {
		return fmt.Errorf("read schema file: %w", err)
	}

	schema, err := typelang.ParseSchema(string(b))
	if err != nil {
		return fmt.Errorf("parse schema file: %w", err)
	}

	g, err := gen.NewGenerator(schema, outdir)
	if err != nil {
		return err
	}

	return g.Generate()
}
