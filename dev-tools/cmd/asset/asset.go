// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// +build ignore

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"strings"

	"github.com/elastic/beats/libbeat/asset"
	"github.com/elastic/beats/licenses"
)

var (
<<<<<<< HEAD
	pkg     string
	input   string
	output  string
	name    string
	license = "ASL2"
=======
	pkg    string
	input  string
	output string
>>>>>>> aa82756e2ff04055bd5c7678a03abc815bec4b32
)

func init() {
	flag.StringVar(&pkg, "pkg", "", "Package name")
	flag.StringVar(&input, "in", "-", "Source of input. \"-\" means reading from stdin")
	flag.StringVar(&output, "out", "-", "Output path. \"-\" means writing to stdout")
<<<<<<< HEAD
	flag.StringVar(&license, "license", "ASL2", "License header for generated file.")
	flag.StringVar(&name, "name", "", "Asset name")
=======
>>>>>>> aa82756e2ff04055bd5c7678a03abc815bec4b32
}

func main() {
	flag.Parse()
	args := flag.Args()

	var (
		file, beatName string
		data           []byte
		err            error
	)
	if input == "-" {
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "File path must be set")
			os.Exit(1)
		}
		file = args[0]
		beatName = args[1]

		r := bufio.NewReader(os.Stdin)
		data, err = ioutil.ReadAll(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while reading from stdin: %v\n", err)
			os.Exit(1)
		}
	} else {
		file = input
		beatName = args[0]
		data, err = ioutil.ReadFile(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid file path: %s\n", input)
			os.Exit(1)
		}
<<<<<<< HEAD
	}

	// Depending on OS or tools configuration, files can contain carriages (\r),
	// what leads to different results, remove them before encoding.
	encData, err := asset.EncodeData(strings.Replace(string(data), "\r", "", -1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding the data: %s\n", err)
		os.Exit(1)
	}

	licenseHeader, err := licenses.Find(license)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid license: %s\n", err)
=======
	}

	// Depending on OS or tools configuration, files can contain carriages (\r),
	// what leads to different results, remove them before encoding.
	encData, err := asset.EncodeData(strings.Replace(string(data), "\r", "", -1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding the data: %s\n", err)
>>>>>>> aa82756e2ff04055bd5c7678a03abc815bec4b32
		os.Exit(1)
	}
	var buf bytes.Buffer
	if name == "" {
		name = file
	}
	asset.Template.Execute(&buf, asset.Data{
		Beat:    beatName,
		Name:    name,
		Data:    encData,
<<<<<<< HEAD
		License: licenseHeader,
=======
>>>>>>> aa82756e2ff04055bd5c7678a03abc815bec4b32
		Package: pkg,
	})

	bs, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	if output == "-" {
		os.Stdout.Write(bs)
	} else {
		ioutil.WriteFile(output, bs, 0640)
	}
}
