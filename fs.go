package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rakyll/statik/fs"
)

// dirIsEmpty returns true if the given path is empty, otherwise false
func dirIsEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	r, err := f.Readdirnames(1)
	if err != nil {
		return false, err
	}
	if len(r) < 1 || io.EOF == err {
		return false, err
	}
	return true, err
}

// linkLibs returns a slice of the discovered shared object filkes in the given path
func linkLibs(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	r, err := f.Readdirnames(1)
	if err != nil {
		return []string{}, err
	}
	return r, err
}

// wrap calls a given command and args, returning the raw byte slice of the combined stdout and stderr outputs, as well
// as any encountered errors. The command is invoked using a context timer, so any compilation options that run for
// longer than compileTimeout (defined in config.go) will be killed using os.Process.Kill.
func wrap(cmd string, args []string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), compileTimeout)
	defer cancel()
	command := exec.CommandContext(ctx, cmd, args...)
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("could not get working directory: %+v", err)
		return nil, err
	}
	command.Dir = wd
	out, err := command.CombinedOutput()
	if err != nil {
		log.Printf("os/exec error: %+v", err)
		return out, err
	}
	return out, nil
}

// findAll takes a given list of file extensions and a target dir and returns all the files with the right extensions.
func findAll(target string, extensions []string) ([]string, error) {
	res := make([]string, 0)
	for _, g := range extensions {
		globs, err := find(target, g)
		if err != nil {
			return res, fmt.Errorf("could not locate source files: %+w", err)
		}
		res = append(res, globs...)
	}
	return res, nil
}

// find returns a slice corresponding to all files in the given path that have the right extension
func find(target, extension string) ([]string, error) {
	res := make([]string, 0)
	err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(extension, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			res = append(res, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// copyTestFramework copies the given Catch2 header and test_main impl from the statik filesystem to the target dir
func copyTestFramework(catchFile, hostFile, testPath string) {
	filesys, err := fs.New()
	if err != nil {
		log.Fatalf("could not find embedded catch2 header: %+v", err)
	}
	b, err := filesys.Open(catchFile)
	if err != nil {
		log.Fatalf(" ould not find embedded catch2 header: %+v", err)
	}
	contents, err := ioutil.ReadAll(b)
	if err != nil {
		log.Fatalf("error reading catch.hpp file: %+v", err)
	}
	ioutil.WriteFile(testPath+"/tests/"+catchFile, contents, 0664)

	h, err := filesys.Open(hostFile)
	if err != nil {
		log.Fatalf("could not find embedded catch2 header: %+v", err)
	}
	hostContents, err := ioutil.ReadAll(h)
	if err != nil {
		log.Fatalf("error reading host.hpp file: %+v", err)
	}
	ioutil.WriteFile(testPath+"/tests"+hostFile, hostContents, 0664)
}
