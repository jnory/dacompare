package main

import (
	"flag"
	"log"
	"os"

	"./common"
	"./try1"
)

func build(path string) error {
	reader, err := common.NewOrderedReader(path)
	if err != nil {
		return err
	}

	log.Print("# of entries: ", reader.Size())

	da, err := try1.NewDA(reader)
	if err != nil {
		return err
	}

	log.Print("Size of DA: ", len(da))

	return nil
}

func main() {
	path := flag.String("data", "", "path to entry file")

	flag.Parse()

	if *path == "" {
		flag.Usage()
		os.Exit(1)
	}

	err := build(*path)
	if err != nil {
		log.Print(err)
	}
}
