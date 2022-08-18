package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bearmini/sexp"
	"github.com/pkg/errors"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

func run() error {
	var (
		wat2Wasm      string
		inputWastFile string
		outputDir     string
	)

	flag.StringVar(&wat2Wasm, "wat2wasm", "", "")
	flag.StringVar(&inputWastFile, "input", "", "One of .wast files in github.com/WebAssembly/spec/test/core")
	flag.StringVar(&outputDir, "output-dir", "", "")
	flag.Parse()

	if wat2Wasm == "" {
		return errors.New("-wat2wasm must be specified")
	}

	if inputWastFile == "" {
		return errors.New("-input must be specified")
	}

	if outputDir == "" {
		return errors.New("-output-dir must be specified")
	}

	f, err := os.Open(inputWastFile)
	if err != nil {
		return err
	}

	r := NewWastReader(f)
	i := 0

	for {
		s, err := r.NextSexp()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if s == nil {
			break
		}

		if len(s.Children) == 0 {
			continue
		}

		if s.Children[0].Atom == nil {
			continue
		}

		if s.Children[0].Atom.Type != sexp.TokenTypeSymbol {
			continue
		}

		if s.Children[0].Atom.Value != "module" {
			continue
		}

		log.Printf("s-expression: %s\n", s.String())

		if (strings.HasSuffix(inputWastFile, "ref_func.wast") && i == 2) || (strings.HasSuffix(inputWastFile, "unreached-valid.wast") && i == 1) {
			i++
			continue
		}

		b, err := compileModule(wat2Wasm, s.String())
		if err != nil {
			return errors.Wrapf(err, "error while compiling wat to wasm: input file = %s, i = %d, s-expression= %s", inputWastFile, i, s.String())
		}

		outputFileName := getOutputFileName(inputWastFile, i)

		path := filepath.Join(outputDir, outputFileName)
		err = ioutil.WriteFile(path, b, 0644)
		if err != nil {
			return err
		}

		i++
	}

	return nil
}

func compileModule(compiler, s string) ([]byte, error) {
	d, err := os.MkdirTemp("", "fq-wasm-test-module-generator")
	if err != nil {
		return nil, err
	}

	fi, err := os.CreateTemp(d, "wat")
	if err != nil {
		return nil, err
	}
	defer os.Remove(fi.Name())

	_, err = fi.WriteString(s)
	if err != nil {
		return nil, err
	}
	fi.Close()

	fo, err := os.CreateTemp(d, "wasm")
	if err != nil {
		return nil, err
	}
	defer os.Remove(fo.Name())
	fo.Close()

	cmd := exec.Command(compiler, fi.Name(), "-o", fo.Name())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "error while compiling wat to wasm: %q, output = %s", cmd.Args, string(out))
	}

	b, err := ioutil.ReadFile(fo.Name())
	if err != nil {
		return nil, err
	}

	return b, nil
}

func getOutputFileName(inputWastFile string, n int) string {
	_, f := filepath.Split(inputWastFile)
	ext := filepath.Ext(f)
	return fmt.Sprintf("%s-%d.wasm", strings.TrimSuffix(f, ext), n)
}
