package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const csvPath = "KEN_ALL.CSV"
const utf8Path = "utf-8.txt"
const zipAddressPath = "zip_address.txt"

type ZipCode string
type Address string

func main() {
	sjisToUtf8()
	utf, err := os.Open(utf8Path)
	if err != nil {
		log.Fatalln(err)
	}
	defer utf.Close()

	f, err := os.Create(zipAddressPath)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(utf)
	for scanner.Scan() {
		removedQuote := strings.ReplaceAll(scanner.Text(), `"`, "")
		z, a := extractZipCodeAddress(removedQuote)
		za := fmt.Sprintf("('%s', '%s')\n", z, a)
		fmt.Fprint(f, za)
	}
}

func sjisToUtf8() {
	csv, err := os.Open(csvPath)
	if err != nil {
		log.Fatalln("ファイルオープンでエラー %v\n", err)
	}
	defer csv.Close()

	reader := transform.NewReader(csv, japanese.ShiftJIS.NewDecoder())

	f, err := os.Create(utf8Path)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	tee := io.TeeReader(reader, f)
	scanner := bufio.NewScanner(tee)

	for scanner.Scan() {
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}

func extractZipCodeAddress(line string) (ZipCode, Address) {
	s := strings.Split(line, ",")
	z := ZipCode(s[2])
	a := Address(fmt.Sprintf("%s%s%s", s[6], s[7], s[8]))
	return z, a
}
