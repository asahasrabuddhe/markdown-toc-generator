package main

import (
	"io"
	"io/ioutil"
	"os"
)

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
		return os.Create(destination)
	}
}

func readFromSource(reader io.Reader) ([]byte, error) {
	return ioutil.ReadAll(reader)
}

func writeToDestination(output []byte, writer io.Writer) error {
	_, err := writer.Write(output)
	return err
}
