// Cm is a tiny tool for automating modern C++ projects.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	name        = flag.String("o", "", "name of the output binary")
	includepath = flag.String("I", "", "path to header files")
	optimize    = flag.Bool("max", false, "maximum optimization")
	std         = flag.String("std", "c++2a", "c++ standard library to use")
	compiler    = flag.String("compiler", "g++", "c++ compiler to use")
)

func wrap(cmd string, args []string) ([]byte, error) {
	command := exec.Command(cmd, args...)
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("💥 compiler error: %+v", err)
		return nil, err
	}
	command.Dir = wd
	out, err := command.CombinedOutput()
	if err != nil {
		log.Printf("💥 compiler error: %+v", err)
		return out, err
	}
	return out, nil
}

func findAll(target string, extensions []string) ([]string, error) {
	res := make([]string, 0)
	for _, g := range extensions {
		globs, err := find(target, g)
		if err != nil {
			return res, fmt.Errorf("💥 could not locate source files: %+w", err)
		}
		res = append(res, globs...)
	}
	return res, nil
}

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

func compile(includepath *string, targetpath string, extra ...string) {
	globs := []string{
		"*.cpp",
		"*.cxx",
		"*.cc",
	}
	targets, err := findAll(targetpath, globs)
	if err != nil {
		log.Fatalf("💥 could not find target files: %+v", err)
	}

	debugArgs := []string{
		"-std=" + *std,
		"-Wall",
		"-O0",
		"-o" + targetpath + "/" + *name,
	}
	debugArgs = append(debugArgs, targets...)
	debugArgs = append(debugArgs, extra...)
	optimizeArgs := []string{
		"-std=" + *std,
		"-Wall",
		"-O3",
		"-o" + targetpath + "/" + *name,
	}
	optimizeArgs = append(optimizeArgs, targets...)
	optimizeArgs = append(optimizeArgs, extra...)

	if *optimize {
		if includepath == nil {
			includepath = &targetpath
		} else {
			if *includepath == "" {
				includepath = &targetpath
			}
			optimizeArgs = append(optimizeArgs, "-I"+*includepath)
			if _, err := os.Stat(*includepath + "/lib"); os.IsExist(err) {
				optimizeArgs = append(optimizeArgs, "-L"+*includepath+"/lib")
			}
		}
		out, err := wrap(*compiler, optimizeArgs)
		if err != nil {
			log.Fatalf("💥 wrap error: %+v, %v", string(out), err)
		}
		if len(out) == 0 {
			log.Println("🎉 compilation succeeded with no errors")
		} else {
			log.Print(string(out))
		}
	} else {
		if includepath != nil {
			if *includepath == "" {
				includepath = &targetpath
			}
			debugArgs = append(debugArgs, "-I"+*includepath)
			if _, err := os.Stat(*includepath + "/lib"); os.IsExist(err) {
				debugArgs = append(debugArgs, "-L"+*includepath+"/lib")
			}
		}
		out, err := wrap(*compiler, debugArgs)
		if err != nil {
			log.Fatalf("💥 wrap error: %+v, %v", string(out), err)
		}
		if len(out) == 0 {
			log.Println("🎉 compilation succeeded with no errors")
		} else {
			log.Print(string(out))
		}
	}
}

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		log.Fatal("💥 missing target dir (where should cm be building?)")
	}
	target := os.Args[len(os.Args)-1]

	if *name == "" {
		path, err := filepath.Abs(target)
		if err != nil {
			log.Fatalf("💥 path error: %+v", err)
		}
		split := strings.Split(path, "/")
		*name = split[len(split)-1]
	}

	log.Printf("✅ binary name: \"%s\"", *name)
	log.Printf("✅ binary output path: \"%s\"", target+"/"+*name)
	log.Printf("✅ maximum optimization? %v", *optimize)
	log.Printf("✅ compiling project...\n")

	compile(includepath, target)
}
