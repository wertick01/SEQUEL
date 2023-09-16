package pkg

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

func Xopen(fName string) (*os.File, *bufio.Reader) {
	var (
		fos *os.File
		err error
	)

	// Open the file (unless we are reading from stdin
	if fName != "-" {
		fos, err = os.Open(fName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't open %s: error: %s\n", fName, err)
			os.Exit(1)
		}
	} else {
		fos = os.Stdin
	}

	// Deal with the compression of the file
	var newReader io.Reader
	err = nil
	if strings.HasSuffix(fName, ".gz") {
		newReader, err = gzip.NewReader(fos)
	} else if strings.HasSuffix(fName, ".bz2") {
		newReader = bzip2.NewReader(fos)
	} else {
		newReader = fos // No compression, stdin
	}

	if err != nil {
		fmt.Fprintf(os.Stderr,
			"%s file is not in the corRect format: error: %s\n", fName, err)
		os.Exit(1)
	}

	// We want a buffered stream to increase performance
	return fos, bufio.NewReader(newReader)
}

type Record struct {
	Name, Seq, Qual string
}

type FqReader struct {
	R               *bufio.Reader
	Last, Seq, Qual []byte
	Finished        bool
	Rec             Record
}

func (fq *FqReader) IterLines() ([]byte, bool) {
	line, err := fq.R.ReadSlice('\n')
	if err != nil {
		if err == io.EOF {
			return line, true
		} else {
			panic(err)
		}
	}
	return line, false
}

var space = []byte(" ")

func (fq *FqReader) Iter() (Record, bool) {
	if fq.Finished {
		return fq.Rec, fq.Finished
	}
	// Read the seq id (fasta or fastq)
	if fq.Last == nil {
		for l, done := fq.IterLines(); !done; l, done = fq.IterLines() {
			if l[0] == '>' || l[0] == '@' { // read id
				fq.Last = l[0 : len(l)-1]
				break
			}
		}
		if fq.Last == nil { // We couldn't find a valid Record, no more data in file
			fq.Finished = true
			return fq.Rec, fq.Finished
		}
	}
	fq.Rec.Name = string(bytes.SplitN(fq.Last, space, 1)[0])
	fq.Last = nil

	// Now read the sequence
	fq.Seq = fq.Seq[:0]
	for l, done := fq.IterLines(); !done; l, done = fq.IterLines() {
		c := l[0]
		if c == '+' || c == '>' || c == '@' {
			fq.Last = l[0 : len(l)-1]
			break
		}
		fq.Seq = append(fq.Seq, l[0:len(l)-1]...)
	}
	fq.Rec.Seq = string(fq.Seq)

	if fq.Last == nil { // Reach EOF
		fq.Finished = true
	}
	if fq.Last[0] != '+' { // fasta Record; set sequence
		return fq.Rec, fq.Finished
	}
	leng := 0
	fq.Qual = fq.Qual[:0]
	for l, done := fq.IterLines(); !done; l, done = fq.IterLines() {
		fq.Qual = append(fq.Qual, l[0:len(l)-1]...)
		leng += len(l)
		if leng >= len(fq.Seq) { // we have read enough quality
			fq.Last = nil
			fq.Rec.Qual = string(fq.Qual)
			return fq.Rec, fq.Finished
		}
	}
	fq.Finished = true
	fq.Rec.Qual = string(fq.Qual)
	return fq.Rec, fq.Finished // incomplete fastq quality, return what we have
}
