package main

import (
	"Sano/compiler"
	"Sano/cpu"
	"Sano/linker"
	"Sano/parser"
	"Sano/vera"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"image"
	_ "image/png"
	_ "image/gif"
)

var ConvertImage = &cli.Command{
	Name:    "convert-image",
	Aliases: []string{"ci"},
	Usage:   "convert an indexed-colour image to VERA's VRAM format",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "mode",
			Aliases:  []string{"m"},
			Required: true,
			Usage:    "bitmap or tiled output",
		},
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Required: true,
			Usage:    "input file",
		},
		&cli.StringFlag{
			Name:     "output",
			Aliases:  []string{"o"},
			Required: true,
			Usage:    "output file",
		},
		&cli.IntFlag{
			Name:     "bpp",
			Aliases:  []string{"b"},
			Required: true,
			Usage:    "the bits per pixel of the output (1, 2, 4, or 8)",
		},
		&cli.IntFlag{
			Name:    "tile-width",
			Aliases: []string{"tw"},
			Value:   8,
			Usage:   "width of tiles for tiled output (8, 16, 32, or 64)",
		},
		&cli.IntFlag{
			Name:    "tile-height",
			Aliases: []string{"th"},
			Value:   8,
			Usage:   "height of tiles for tiled output (8, 16, 32, or 64)",
		},
	},
	Action: func(ctx *cli.Context) error {
		file, err := os.Open(ctx.String("input"))
		defer file.Close()
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}

		img, _, err := image.Decode(file)
		if err != nil {
			return fmt.Errorf("failed to decode image: %w", err)
		}

		pal, ok := img.(image.PalettedImage)
		if !ok {
			return fmt.Errorf("image is not in indexed colour")
		}

		outputFile, err := os.OpenFile(ctx.String("output"), os.O_RDWR|os.O_CREATE, 0o660)
		defer outputFile.Close()
		if err != nil {
			return fmt.Errorf("failed to open output file: %w", err)
		}

		bpp, ok := vera.ToBPP(ctx.Int("bpp"))
		if !ok {
			return fmt.Errorf("bpp must be one of 1, 2, 4, or 8, but it was %d", ctx.Int("bpp"))
		}

		mode := ctx.String("mode")
		if mode == "tiled" || mode == "t" {
			h, ok := vera.ToTileSize(ctx.Int("tile-height"))
			if !ok {
				return fmt.Errorf("tile-height must be one of 8, 16, 32, or 64, but it was %d", ctx.Int("tile-height"))
			}
			w, ok := vera.ToTileSize(ctx.Int("tile-width"))
			if !ok {
				return fmt.Errorf("tile-width must be one of 8, 16, 32, or 64, but it was %d", ctx.Int("tile-width"))
			}
			err = vera.ExportTile(pal, bpp, w, h, outputFil)egi
			if err != nil {
				return fmt.Errorf("failed to export tile data: %w", err)
			}
		} else if mode == "bitmap" || mode == "b" {
			err = vera.ExportBitmap(pal, bpp, file)
			if err != nil {
				return fmt.Errorf("failed to export bitmap data: %w", err)
			}
		} else {
			return fmt.Errorf("unknown mode: %s", mode)
		}

		return nil
	},
}

var Assembler = &cli.Command{
	Name:  "assembler",
	Usage: "WIP assembler and linker",
	Action: func(ctx *cli.Context) error {
		data, err := os.ReadFile(ctx.Args().Get(0))
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}
		g, err := parser.Parser.ParseBytes(ctx.Args().Get(0), data)
		if err != nil {
			return fmt.Errorf("failed to parse input file: %w", err)
		}

		cx16 := cpu.Base6502Opcodes.And(cpu.WDC65C02ExtensionOpcodes)
		c := compiler.Compiler{Instructions: cx16}
		obj, errors := c.Compile(g)
		if len(errors) > 0 {
			for _, err := range errors {
				println(err.String())
			}
			os.Exit(1)
		}

		prg, err := linker.LinkToPrg([]*linker.Object{obj})
		if err != nil {
			return fmt.Errorf("failed to link file into prg: %w", err)
		}

		err = os.WriteFile(ctx.Args().Get(1), prg, 0660)
		if err != nil {
			return fmt.Errorf("failed to write prg file: %w", err)
		}

		return nil
	},
}

func main() {
	app := &cli.App{
		Name:  "sano",
		Usage: "commander x16 developer's toolbox",
		ExitErrHandler: func(cCtx *cli.Context, err error) {
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
		Commands: []*cli.Command{ConvertImage, Assembler},
	}
	app.Run(os.Args)
}
