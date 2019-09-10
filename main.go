package main

import (
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "mdtocgen"
	app.Usage = "Generate beautiful table of contents for your markdown files"
	app.UseShortOptionHandling = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "source, s",
			Value: "-",
			Usage: "path to markdown file to parse for generating table of contents or '-' to read from stdin",
		},
		cli.StringFlag{
			Name:  "destination, d",
			Value: "-",
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
		cli.IntFlag{
			Name:  "max-depth, m",
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
		depth := context.Int("max-depth")
		skipHeaders := context.Int("skip-headers")

		// get reader
		reader, err := getSourceReader(source)
		if err != nil {
			return err
		}
		defer reader.Close()

		// get source input
		input, err := readFromSource(reader)
		if err != nil {
			return err
		}

		// generate toc
		output, err := generateToc(input, depth, skipHeaders)
		if err != nil {
			return err
		}

		if outputCompleteFile {
			// replace
		}

		// get writer
		var writer io.WriteCloser
		if inline {
			writer, err = getSourceWriter(source)
		} else {
			writer, err = getSourceWriter(destination)
		}
		if err != nil {
			return err
		}
		defer writer.Close()

		// write output
		err = writeToDestination(output, writer)
		if err != nil {
			return err
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getSourceReader(source string) (io.ReadCloser, error) {
	if source == "-" {
		return os.Stdin, nil
	} else {
		return os.Open(source)
	}
}

func getSourceWriter(destination string) (io.WriteCloser, error) {
	if destination == "-" {
		return os.Stdout, nil
	} else {
		return os.Open(destination)
	}
}

func readFromSource(reader io.Reader) ([]byte, error) {
	return ioutil.ReadAll(reader)
}

func writeToDestination(output []byte, writer io.Writer) error {
	_, err := writer.Write(output)
	return err
}
