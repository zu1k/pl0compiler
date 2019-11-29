package fileio

import (
	"bufio"
	"log"
	"os"
)

type File struct {
	file    *os.File
	scanner *bufio.Scanner
}

func (f *File) Open(filepath string) {
	var err error
	f.file, err = os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	f.scanner = bufio.NewScanner(f.file)
	f.scanner.Split(bufio.ScanRunes)
}

func (f *File) ReadRune() (r rune, end bool) {
	if f.scanner.Scan() {
		r = []rune(f.scanner.Text())[0]
		end = false
	} else {
		end = true
	}
	return
}
