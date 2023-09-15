package main

import (
	"flag"
)

type fileFlag []string

func (f *fileFlag) String() string {
	str := ""
	for _, s := range *f {
		str += s
		str += ", "
	}
	return str[:]
}

func (f *fileFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func initArgs() ([]string, string) {
	var files fileFlag
	files = []string{}
	flag.Var(&files, "f", "Specify input files (must be in csv format!)")
	outputFile := ""
	flag.StringVar(&outputFile, "o", "", "output file")
	flag.Parse()

	return files, outputFile
}
