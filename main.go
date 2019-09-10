package main

import (
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "mdtocgen"
	app.Usage = "Generate beautiful table of contents for your markdown files"
	app.UseShortOptionHandling = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "source, s",
			Value: "",
			Usage: "path to markdown file to parse for generating table of contents or '-' to read from stdin",
		},
		cli.StringFlag{
			Name:  "destination, d",
			Value: "",
			Usage: "path to markdown file to write after table of contents has been generated or '-' to write to stdout",
		},
		cli.BoolTFlag{
			Name:  "output-complete-file, c",
			Usage: "entire source file is written as output including the table of contents. if this flag is set to false, inline option is ignored",
		},
		cli.BoolFlag{
			Name:  "inline, i",
			Usage: "appends the source file with the table of contents without creating a new file",
		},
		cli.StringSliceFlag{
			Name:  "bullets, b",
			Usage: "bullets to be used for items generated in the table of contents. multiple bullets also supported which would be used based on the depth of header",
			Value: &cli.StringSlice{"*"},
		},
		cli.IntFlag{
			Name:  "max-depth, d",
			Usage: "use headings with nested depth `n` to generate table of contents items",
			Value: 6,
		},
		cli.IntFlag{
			Name:  "skip-headers, k",
			Usage: "skip first `n` headers when generating the table of contents",
		},
	}

	app.Action = func(context *cli.Context) error {
		// parse options
		source := context.String("source")
		destination := context.String("destination")
		outputCompleteFile := context.Bool("output-complete-file")
		inline := context.Bool("inline")
		bullets := context.StringSlice("bullets")
		depth := context.Int("max-depth")
		skipHeaders := context.Int("skip-headers")

		// read source

		// generate toc

		// write output
	}
}