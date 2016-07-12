package main

import (
	"encoding/xml"
	"os"
	"fmt"
	"io/ioutil"
	"bufio"
	"bytes"
	"strings"
)

type Resources struct {
	XMLName xml.Name	`xml:"resources"`
	Strings []StringRes	`xml:"string"`
}

type StringRes struct {
	Name string		`xml:"name,attr"`
	Value string		`xml:",chardata"`
}

func main() {
	if len(os.Args) < 3 {
		printHelp()
		os.Exit(0)
	}
	var outputFile string
	inputFile := os.Args[1]
	localizableFile := os.Args[2]
	if len(os.Args) > 3 {
		outputFile = os.Args[3]
	} else {
		outputFile = inputFile + ".strings"
	}
	fmt.Printf("Input: %s, output: %s\n", inputFile, outputFile)

	// read XML
	res := Resources{}
	xmlContent,_ := ioutil.ReadFile(inputFile)
	var err = xml.Unmarshal(xmlContent, &res)
	if err != nil { panic(err) }

	androidStrings := res.Strings
	fmt.Printf("Found: %d Android string resources!\n", len(androidStrings))

	// put string resources to map keyed by name
	aMap := make(map[string]StringRes)
	for _, i := range androidStrings {
		aMap[i.Name] = i
	}

	// read Localizable.strings
	format := "\"%s\" = \"%s\";\n"
	locBytes, _ := ioutil.ReadFile(localizableFile)
	scanner := bufio.NewScanner(bytes.NewReader(locBytes))
	locKeys := []string{}
	locValues := make(map[string]StringRes)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "\"") {
			// fmt.Println(line)
			elements := strings.Split(line, "=")
			// fmt.Println(elements)
			name := strings.TrimSpace(elements[0])
			name = strings.TrimSuffix(name, "\"")
			name = strings.TrimPrefix(name, "\"")

			value := strings.TrimSpace(elements[1])
			value = strings.TrimPrefix(value, "\"")
			value = strings.TrimSuffix(value, ";")
			value = strings.TrimSuffix(value, "\"")

			locValues[name] = StringRes{ Name: name, Value: value }
			locKeys = append(locKeys, name)
		}

	}
	fmt.Printf("Found: %d ObjC Localizable string resources!\n", len(locKeys))


	// write to file
	f, err := os.Create(outputFile)
	if err != nil { panic(err) }
	w := bufio.NewWriter(f)
	for _, k := range locKeys {
		element, ok := aMap[k]
		if ok {
			_, err = fmt.Fprintf(w, format, element.Name, element.Value)
			if err != nil { panic(err) }
		} else {
			fmt.Printf("Failed to find Android resource with key: %s!\n", k)
			_, err = fmt.Fprintln(w, "\n/** Missed translation */")
			if err != nil { panic(err) }

			_, err = fmt.Fprintf(w, format, k, locValues[k].Value)
			if err != nil { panic(err) }

			_, err = fmt.Fprintln(w, "")
		}
	}

	w.Flush()
}

func printHelp()  {
	fmt.Println("Arguments:")
	fmt.Println("inputFile - the path to the Android String XML file")
	fmt.Println("locFile - the path to the Localizable.strings file which has ethalon keys")
	fmt.Println("outputFile - the path to the output file (optional). If missed than name of input will be used by appending .strings suffix")
}
