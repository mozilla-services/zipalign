package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var bias int
	var verbose = flag.Bool("v", false, "Verbose output")
	var alignment = flag.Int("a", 4, "Alignment in bytes, e.g. '4' provides 32-bit alignment")
	var inputFile = flag.String("i", "", "Input ZIP file to be aligned")
	var outputFile = flag.String("o", "", "Output aligned ZIP file")
	var overwrite = flag.Bool("f", false, "Overwrite existing outfile.zip")
	var help = flag.Bool("h", false, "Print this help")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" || *help {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *inputFile == *outputFile && !*overwrite {
		log.Fatalf("Refusing to overwrite output file %q without -f being set", *outputFile)
	}
	if *verbose {
		log.Printf("Aligning %q on %d bytes and writing out to %q", *inputFile, *alignment, *outputFile)
	}
	// Open a zip archive for reading.
	r, err := zip.OpenReader(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Iterate through the files in the archive,
	for _, f := range r.File {
		if *verbose {
			log.Printf("Processing %q from input archive", f.Name)
		}
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		var padlen int
		if f.CompressedSize64 != f.UncompressedSize64 {
			// File is compressed, copy the entry without padding
			if *verbose {
				log.Printf("--- %s: len %d (compressed)", f.Name, f.UncompressedSize64)
			}
		} else {
			// source: https://android.googlesource.com/platform/build.git/+/android-4.2.2_r1/tools/zipalign/ZipAlign.cpp#76
			newOffset := len(f.Extra) + bias
			padlen = (*alignment - (newOffset % *alignment)) % *alignment
			if *verbose && padlen > 0 {
				log.Printf(" --- %s: padding %d bytes", f.Name, padlen)
			}
		}

		fwhead := &zip.FileHeader{
			Name:   f.Name,
			Method: zip.Deflate,
		}
		// add padlen number of null bytes to the extra field of the file header
		// in order to align files on 4 bytes
		for i := 0; i < padlen; i++ {
			fwhead.Extra = append(fwhead.Extra, '\x00')
		}

		fw, err := w.CreateHeader(fwhead)
		if err != nil {
			log.Fatal(err)
		}
		data, err := ioutil.ReadAll(rc)
		if err != nil {
			log.Fatal(err)
		}
		_, err = fw.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		rc.Close()
		bias += padlen
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Write the aligned zip file
	err = ioutil.WriteFile(*outputFile, buf.Bytes(), 0744)
	if err != nil {
		log.Fatal(err)
	}
}
