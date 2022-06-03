package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"
)

// Example go-license output
// github.com/tinylib/msgp/msgp,https://github.com/tinylib/msgp/blob/v1.1.6/LICENSE,MIT
type goLicenseOutput struct {
	lib        string
	licenseUrl string
	license    string
}

// Example solicitor input
// org.eclipse.xtend;xtend;2.2.0;MIT;https://spdx.org/licenses/MIT#licenseText
type solLicenseInput struct {
	lib        string
	version    string
	license    string
	licenseUrl string
}

func main() {

	target := flag.String("t", "license.csv", "target csv from go-license")
	output := flag.String("o", "sol_input.csv", "target output file for solicitor")
	flag.Parse()

	f, err := os.Open(*target)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	slis := make([]solLicenseInput, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		glo := NewGoLicenseOutput(line)

		sli := NewSolLicenseInput(glo)

		slis = append(slis, sli)

	}

	fo, err := os.Create(*output)
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	WriteFromSliSlice(slis, fo)

}

func WriteFromSliSlice(slis []solLicenseInput, out io.Writer) int {

	//create a buffer for writing to file
	w := bufio.NewWriter(out)
	overall := 0
	//write all objects to the buffer
	for i := range slis {

		buf, err := w.WriteString(slis[i].String() + "\n")

		if err != nil {
			panic(err)
		}
		overall += buf

	}
	//write to the file
	w.Flush()
	return overall

}

func NewGoLicenseOutput(raw string) goLicenseOutput {

	rawSplit := strings.Split(raw, ",")
	glo := goLicenseOutput{

		lib:        rawSplit[0],
		licenseUrl: rawSplit[1],
		license:    rawSplit[2],
	}

	return glo

}

func (sli *solLicenseInput) String() string {

	return ";" + sli.lib + ";" + sli.version + ";" + sli.licenseUrl + ";" + sli.license

}

func NewSolLicenseInput(glo goLicenseOutput) solLicenseInput {

	v := extractVersionFromLicenseUrl(glo.licenseUrl)

	sli := solLicenseInput{

		lib:        glo.lib,
		licenseUrl: glo.licenseUrl,
		license:    glo.license,
		version:    v,
	}

	return sli

}

// Example LicenseURL
// https://github.com/tinylib/msgp/blob/v1.1.6/LICENSE
func extractVersionFromLicenseUrl(lurl string) string {

	// remove "LICENSE"
	out := lurl[:len(lurl)-8]
	ind := strings.LastIndex(out, "/")

	out = out[ind+1:]

	return out

}
