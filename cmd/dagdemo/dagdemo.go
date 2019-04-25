// Copyright (c) 2018-2019 The Soteria Engineering developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/soteria-dag/soterd/soterutil"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/soteria-dag/soterd/chaincfg"
	"github.com/soteria-dag/soterd/integration/rpctest"
	"github.com/soteria-dag/soterd/rpcclient"
)

// runNet runs a network of miners, generates some blocks on them, renders the dag as html, and returns the file name
// of the rendered html.
func runNet(minerCount, blockCount int, output string) (string, error) {
	var miners []*rpctest.Harness

	// Spawn miners
	for i := 0; i < minerCount; i++ {
		miner, err := rpctest.New(&chaincfg.SimNetParams, nil, nil, false)
		if err != nil {
			return "", fmt.Errorf("unable to create mining node %d: %s", i, err)
		}

		if err := miner.SetUp(false, 0); err != nil {
			return "", fmt.Errorf("unable to complete mining node %d setup: %s", i, err)
		}

		miners = append(miners, miner)
	}
	// NOTE(cedric): We'll call defer on a single anonymous function instead of minerCount times in the above loop
	defer func() {
		for _, miner := range miners {
			_ = (*miner).TearDown()
		}
	}()

	// Connect the nodes to one another
	err := rpctest.ConnectNodes(miners)
	if err != nil {
		return "", fmt.Errorf("unable to connect nodes: %s", err)
	}

	// Generate blocks on each miner.
	var futures []*rpcclient.FutureGenerateResult
	for _, miner := range miners {
		future := miner.Node.GenerateAsync(uint32(blockCount))
		futures = append(futures, &future)
	}

	// Wait for block generation to finish
	for i, future := range futures {
		_, err := (*future).Receive()
		if err != nil {
			return "", fmt.Errorf("failed to wait for blocks to generate on node %d: %s", i, err)
		}
	}

	// Render the dag in graphviz DOT file format
	dot, err := rpctest.RenderDagsDot(miners)
	if err != nil {
		return "", fmt.Errorf("failed to render dag in graphviz DOT format: %s", err)
	}

	// Convert DOT file contents to an SVG image
	svg, err := soterutil.DotToSvg(dot)
	if err != nil {
		return "", fmt.Errorf("failed to convert DOT file to SVG: %s", err)
	}
	
	// We're going to embed the SVG image in HTML, so strip out the xml declaration
	svgEmbed, err := soterutil.StripSvgXmlDecl(svg)
	if err != nil {
		return "", fmt.Errorf("failed to strip xml declaration from SVG image: %s", err)
	}
	
	// Render the dag in an HTML document
	h, err := soterutil.RenderSvgHTML(svgEmbed, "dag")
	if err != nil {
		return "", fmt.Errorf("failed to render SVG image as HTML: %s", err)
	}

	// Determine where to save HTML document
	var fh *os.File
	pattern := "dag_*.html"
	if len(output) == 0 {
		// Save to randomly-named file in the system's tempdir
		fh, err = ioutil.TempFile("", pattern)
	} else {
		info, pathErr := os.Stat(output)
		if pathErr == nil && info.IsDir() {
			// Save to randomly-named file in provided path
			fh, err = ioutil.TempFile(output, pattern)
		} else {
			// Save to provided file name
			fh, err = os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		}
	}

	if err != nil {
		return "", fmt.Errorf("failed to create output file-handle: %s", err)
	}

	// Save the HTML document
	err = save(h, fh)
	if err != nil {
		return "", fmt.Errorf("failed to save HTML file: %s", err)
	}

	return fh.Name(), nil
}

// save bytes to a file descriptor
func save(bytes []byte, fh *os.File) error {
	_, err := fh.Write(bytes)
	if err != nil {
		return err
	}

	err = fh.Close()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var output string
	flag.StringVar(&output, "o", "", "Where to save the rendered dag")
	flag.Parse()

	fmt.Println("Generating dag")
	htmlFile, err := runNet(4, 50, output)
	if err != nil {
		fmt.Println(err)
		syscall.Exit(1)
	}
	fmt.Println("Saved dag to", htmlFile)
}

