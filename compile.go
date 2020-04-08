package main

import (
	"fmt"
	"io"
	"log"
	"runtime"
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
	binaryPath := strings.Replace(targetpath, "/src", "/bin", 1) + "/"
	binaryNameFQ := binaryPath + *name

	optLevel := "0"
	if *optimize {
		optLevel = "3"
	}

	cArgs := []string{
		"-std=" + *std,
		"-Wall",
		"-O" + optLevel,
		"-o" + binaryNameFQ,
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
			cArgs = append(cArgs, link)
			libs, err := linkLibs(libPath)
			if err != nil {
				log.Fatalf("error reading shared libraries: %+v", err)
			}
			for _, l := range libs {
				if strings.HasSuffix(l, ".so") {
					striplib := strings.Replace(l, "lib", "", 1)
					trimmedLib := strings.Replace(striplib, ".so", "", 1)
					log.Printf("linking shared object: %s", trimmedLib)
					cArgs = append(cArgs, "-l"+trimmedLib)
				}
				if runtime.GOOS == "darwin" {
					log.Println("compiling on macos, so rpath linking will be done after compilation")
				} else {
					cArgs = append(cArgs, "-Wl,-rpath,"+binaryPath)
				}
			}
		} else {
			log.Println("none found")
		}
	}
	out, err := wrap(*compiler, cArgs)
	if err != nil || *debug {
		fmt.Println("printing wrapped call for debugging:")
		fmt.Println("═════════════════════════════════════")
		fmt.Printf("\n%s %s", *compiler, strings.Join(cArgs, " "))
		fmt.Printf("\n\n")
		fmt.Println("═════════════════════════════════════")
		if err != nil {
			log.Fatalf("reason: %+v, %v", string(out), err)
		}
	}
	if len(out) == 0 {
		log.Println("🎉 compilation succeeded with no errors")
		if runtime.GOOS == "darwin" {
			libs, err := dirIsEmpty(libPath)
			if err != nil && err != io.EOF {
				log.Fatalf("error reading lib dir: %+v", err)
			}
			if libs {
				libs, err := linkLibs(libPath)
				if err != nil {
					log.Fatalf("error reading shared libraries: %+v", err)
				}
				for _, l := range libs {
					// -id "@loader_path/lib/libhello.so" bin/example
					rPathArgs := []string{
						"-id",
						"\"@loader_path/lib/" + l + "\"",
						binaryNameFQ,
					}
					out, err := wrap(
						"install_name_tool",
						rPathArgs,
					)
					if err != nil {
						log.Printf("error with mac rpath tool: %v", err)
						fmt.Printf("\n%s %s", "install_name_tool", strings.Join(rPathArgs, " "))
						fmt.Printf("\n\n")
						log.Println(string(out))
					}
					if len(out) == 0 {
						log.Println("🎉 dynamic linking succeeded with no errors")
					}

				}
			}
		}
	} else {
		log.Print(string(out))
	}
}
