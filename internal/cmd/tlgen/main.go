package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xelaj/mtproto/internal/cmd/tlgen/gen"
	"github.com/xelaj/mtproto/internal/cmd/tlgen/tlparser"
)

const helpMsg = `tlgen
usage: tlgen input_file.tl output_dir/

THIS TOOL IS USING ONLY FOR AUTOMATIC CODE
GENERATION, DO NOT GENERATE FILES BY HAND!

No, seriously. Don't. go generate is amazing. You
are amazing too, but lesser 😏
`
const license = `Copyright (c) 2020 KHS Films

This file is a part of mtproto package.
See https://github.com/xelaj/mtproto/blob/master/LICENSE for details
`

func main() {
	if len(os.Args) != 3 {
		fmt.Println(helpMsg)
		return
	}

	if err := root(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func root(tlfile, outdir string) error {
	b, err := ioutil.ReadFile(tlfile)
	if err != nil {
		return fmt.Errorf("read schema file: %w", err)
	}

	schema, err := tlparser.ParseSchema(fmt.Sprintf("%s", b))
	if err != nil {
		return fmt.Errorf("parse schema file: %w", err)
	}

	g, err := gen.NewGenerator(schema, license, outdir)
	if err != nil {
		return err
	}

	return g.Generate()
}
