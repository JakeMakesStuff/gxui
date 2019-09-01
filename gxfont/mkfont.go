// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// Small program to generate roboto_regular.go and droid_sans_mono.go.
//
// TODO(crawshaw): italics, bold.
package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var urls = map[string]string{
	"roboto_regular":  "https://github.com/jakemakesstuff/fonts/raw/master/apache/roboto/Roboto-Regular.ttf",
	"droid_sans_mono": "https://github.com/jakemakesstuff/fonts/raw/master/apache/droidsansmono/DroidSansMono.ttf",
}

func write(name, url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading %q: %v", url, err)
	}
	data, err = deflate(data)
	if err != nil {
		log.Fatalf("error deflating %q: %v", url, err)
	}

	w := new(bytes.Buffer)
	fmt.Fprint(w, "// GENERATED BY mkfont.go\n\n")
	fmt.Fprint(w, "package gxfont\n\n")
	fmt.Fprintf(w, "var %s = []byte{\n", name)

	for len(data) > 0 {
		n := 16
		if n > len(data) {
			n = len(data)
		}
		for _, c := range data[:n] {
			fmt.Fprintf(w, "0x%02x,", c)
		}
		fmt.Fprintf(w, "\n")
		data = data[n:]
	}

	fmt.Fprint(w, "}")
	wbytes := w.Bytes()

	b, err := format.Source(wbytes)
	if err != nil {
		os.Stderr.Write(wbytes)
		log.Fatalf("error formatting: %v", err)
	}

	f, err := os.Create(name + ".go")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(b); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	for name, url := range urls {
		write(name, url)
	}
}

func deflate(src []byte) (dst []byte, err error) {
	buf := new(bytes.Buffer)
	w, err := flate.NewWriter(buf, flate.DefaultCompression)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(src); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
