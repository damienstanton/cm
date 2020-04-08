package main

import (
	"io"
	"log"
	"strings"
)

// compile executes the compilation process with the given compiler and arguments
func compile(includepath *string, targetpath string, extra ...string) {
	var libPath string
	if *testMode {
		targetpath = targetpath + "/tests"
		libPath = strings.Replace(targetpath, "/tests", "/lib", 1)
	} else {
		targetpath = targetpath + "/src"
		libPath = strings.Replace(targetpath, "/src", "/lib", 1)
	}
	globs := []string{
		"*.cpp",
		"*.cxx",
		"*.cc",
	}
	targets, err := findAll(targetpath, globs)
	if err != nil {
		log.Fatalf("could not find target files: %+v", err)
	}
	binaryOutput := strings.Replace(targetpath, "/src", "/bin", 1) + "/" + *name

	optLevel := "0"
	if *optimize {
		optLevel = "3"
	}

	cArgs := []string{
		"-std=" + *std,
		"-Wall",
		"-O" + optLevel,
		"-o" + binaryOutput,
	}
	cArgs = append(cArgs, targets...)
	cArgs = append(cArgs, extra...)

	if includepath == nil {
		includepath = &targetpath
	} else {
		if *includepath == "" {
			includepath = &targetpath
		}
		cArgs = append(cArgs, "-I"+*includepath)
		log.Printf("checking for shared objects in %s...", libPath)
		libs, err := dirIsEmpty(libPath)
		if err != nil && err != io.EOF {
			log.Fatalf("error reading lib dir: %+v", err)
		}
		if libs {
			link := "-L" + libPath
			log.Printf("found objects in %s. linking", libPath)
			cArgs = append(cArgs, link)
			libs, err := linkLibs(libPath)
			if err != nil {
				log.Fatalf("error reading shared libraries: %+v", err)
			}
			for _, l := range libs {
				cArgs = append(cArgs, "-l:"+l)
			}
		} else {
			log.Println("none found")
		}
	}
	// fmt.Println("DEBUG:", cArgs)
	out, err := wrap(*compiler, cArgs)
	if err != nil {
		log.Fatalf("wrap error: %+v, %v", string(out), err)
	}
	if len(out) == 0 {
		log.Println("🎉 compilation succeeded with no errors")
	} else {
		log.Print(string(out))
	}
}
