/*
Copyright 2011 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Author: lode.vandevenne@gmail.com (Lode Vandevenne)
Author: jyrki.alakuijala@gmail.com (Jyrki Alakuijala)
*/

/*
Zopfli compressor program. It can output gzip-, zlib- or deflate-compatible
data. By default it creates a .gz file. This tool can only compress, not
decompress. Decompression can be done by any standard gzip, zlib or deflate
decompressor.
*/

package main

import (
	"flag"
	"fmt"
	"github.com/foobaz/go-zopfli/zopfli"
	"io"
	"io/ioutil"
	"os"
	"runtime/pprof"
)

// outfilename: filename to write output to, or 0 to write to stdout instead
func compressFile(options *zopfli.Options, outputType int,
	inFileName, outFileName string) error {
	in, inErr := ioutil.ReadFile(inFileName)
	if inErr != nil {
		return inErr
	}

	var out io.WriteCloser
	if outFileName == "" {
		out = os.Stdout
	} else {
		var outErr error
		out, outErr = os.Create(outFileName)
		if outErr != nil {
			return outErr
		}
		defer out.Close()
	}

	compressErr := zopfli.Compress(options, outputType, in, out)
	if compressErr != nil {
		return compressErr
	}

	return nil
}

func main() {
	options := zopfli.DefaultOptions()

	flag.BoolVar(&options.Verbose, "v", options.Verbose, "verbose mode")
	flag.BoolVar(&options.VerboseMore, "vv", options.VerboseMore, "more verbose mode")
	outputToStdout := flag.Bool("c", false, "write the result on standard output, instead of disk")
	deflate := flag.Bool("deflate", false, "output to deflate format instead of gzip")
	zlib := flag.Bool("zlib", false, "output to zlib format instead of gzip")
	gzip := flag.Bool("gzip", true, "output to gzip format")
	flag.BoolVar(&options.BlockSplittingLast, "splitlast", options.BlockSplittingLast, "do block splitting last instead of first")
	flag.IntVar(&options.NumIterations, "i", options.NumIterations, "perform # iterations (default 15). More gives more compression but is slower. Examples: -i=10, -i=50, -i=1000")
	var cpuProfile string
	flag.StringVar(&cpuProfile, "cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if options.VerboseMore {
		options.Verbose = true
	}
	var outputType int
	if *deflate && !*zlib && !*gzip {
		outputType = zopfli.FORMAT_DEFLATE
	} else if *zlib && !*deflate && !*gzip {
		outputType = zopfli.FORMAT_ZLIB
	} else {
		outputType = zopfli.FORMAT_GZIP
	}

	if options.NumIterations < 1 {
		fmt.Fprintf(os.Stderr, "Error: must have 1 or more iterations")
		return
	}

	var allFileNames []string
	if *outputToStdout {
		allFileNames = append(allFileNames, "")
	} else {
		allFileNames = flag.Args()
	}
	if len(allFileNames) <= 0 {
		fmt.Fprintf(os.Stderr, "Please provide filename\n")
	}
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err == nil {
			pprof.StartCPUProfile(f)
			defer f.Close()
			defer pprof.StopCPUProfile()
		}
	}
	for _, fileName := range allFileNames {
		var outFileName string
		if *outputToStdout {
			outFileName = ""
		} else {
			switch outputType {
			case zopfli.FORMAT_GZIP:
				outFileName = fileName + ".gz"
			case zopfli.FORMAT_ZLIB:
				outFileName = fileName + ".zlib"
			case zopfli.FORMAT_DEFLATE:
				outFileName = fileName + ".deflate"
			default:
				panic("Unknown output type")
			}
			if options.Verbose {
				fmt.Fprintf(os.Stderr, "Saving to: %s\n", outFileName)
			}
		}
		compressErr := compressFile(&options, outputType, fileName, outFileName)
		if compressErr != nil {
			fmt.Fprintf(os.Stderr, "could not compress %s: %v\n", fileName, compressErr)
		}
	}
}
