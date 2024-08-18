package main

import (
	"flag"
	"log"
	"os"

	"lc3-vm/vm"
)

func main() {
	log.Println("lc3-vm starting")
	cpu := vm.NewCPU()

	imagePath := getImagePath()
	log.Printf("Loading image: %s", imagePath)
	err := cpu.LoadImage(imagePath)
	if err != nil {
		panic(err)
	}

	cpu.Run()

	log.Println("lc3-vm stopped")
}

func getImagePath() string {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Printf("Usage: lc3 image-file")
		os.Exit(2)
	}
	if info, err := os.Stat(args[0]); err != nil {
		log.Printf("Image not found")
		os.Exit(1)
	} else if info.IsDir() {
		log.Printf("Image must be file")
		os.Exit(1)
	}
	return args[0]
}
