// Program `cm` is a tiny tool for automating modern C++ projects.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/damienstanton/cm/statik"
	"github.com/rakyll/statik/fs"
)

const catchVersion = "v2.11.3"

var (
	name        = flag.String("o", "", "name of the output binary")
	includepath = flag.String("I", "", "path to header files")
	optimize    = flag.Bool("max", false, "maximum optimization")
	std         = flag.String("std", "c++2a", "c++ standard library to use")
	compiler    = flag.String("compiler", "g++", "c++ compiler to use")
	testMode    = flag.Bool("test", false, "run tests using Catch2")
)

func wrap(cmd string, args []string) ([]byte, error) {
	command := exec.Command(cmd, args...)
	wd, err := os.Getwd()
	if err != nil {
		if err.Error() == "exit status 1" {
			log.Println("test failure!")
		}
		log.Printf("💥 wrap error: %+v", err)
		return nil, err
	}
	command.Dir = wd
	out, err := command.CombinedOutput()
	if err != nil {
		if err.Error() == "exit status 1" {
			log.Println("test failure!")
		}
		log.Printf("💥 wrap error: %+v", err)
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
	if *testMode {
		targetpath = targetpath + "/tests"
	} else {
		targetpath = targetpath + "/src"
	}
	globs := []string{
		"*.cpp",
		"*.cxx",
		"*.cc",
	}
	targets, err := findAll(targetpath, globs)
	if err != nil {
		log.Fatalf("💥 could not find target files: %+v", err)
	}

	binaryOutput := strings.Replace(targetpath, "/src", "/bin", 1) + "/" + *name
	debugArgs := []string{
		"-std=" + *std,
		"-Wall",
		"-O0",
		"-o" + binaryOutput,
	}
	debugArgs = append(debugArgs, targets...)
	debugArgs = append(debugArgs, extra...)
	optimizeArgs := []string{
		"-std=" + *std,
		"-Wall",
		"-O3",
		"-o" + binaryOutput,
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
			if _, err := os.Stat(targetpath + "/lib"); os.IsExist(err) {
				optimizeArgs = append(optimizeArgs, "-L"+targetpath+"/lib")
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
			if _, err := os.Stat(targetpath + "/lib"); os.IsExist(err) {
				debugArgs = append(debugArgs, "-L"+targetpath+"/lib")
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

func copyTestFramework(catchFile, hostFile, testPath string) {
	filesys, err := fs.New()
	if err != nil {
		log.Fatalf("💥 could not find embedded catch2 header: %+v", err)
	}
	b, err := filesys.Open(catchFile)
	if err != nil {
		log.Fatalf("💥 could not find embedded catch2 header: %+v", err)
	}
	contents, err := ioutil.ReadAll(b)
	if err != nil {
		log.Fatalf("💥 error reading catch.hpp file: %+v", err)
	}
	ioutil.WriteFile(testPath+"/tests/"+catchFile, contents, 0664)

	h, err := filesys.Open(hostFile)
	if err != nil {
		log.Fatalf("💥 could not find embedded catch2 header: %+v", err)
	}
	hostContents, err := ioutil.ReadAll(h)
	if err != nil {
		log.Fatalf("💥 error reading host.hpp file: %+v", err)
	}
	ioutil.WriteFile(testPath+"/tests"+hostFile, hostContents, 0664)
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

	if *testMode {
		catchFile := "/catch.hpp"
		hostFile := "/test_main.cpp"

		log.Println("🔍 entering test mode...")
		copyTestFramework(catchFile, hostFile, target)
		log.Printf("🔍 compiling catch2 and tests (this may take a while)...\n")
		compile(includepath, target)

		testBinary := target + "/tests/" + *name
		log.Printf("🔍 running %s tests using catch %s", testBinary, catchVersion)
		out, _ := wrap(testBinary, []string{}) // ignore this error as it just indicates test failures
		log.Println(string(out))

		log.Println("🔍 cleaning up test framework...")
		err := os.Remove(target + "/tests/" + *name)
		err = os.Remove(target + "/tests/" + hostFile)
		err = os.Remove(target + "/tests/" + catchFile)
		if err != nil {
			log.Fatalf("💥 cleanup error: %+v", err)
		}
		log.Println("🔍 exited test mode")
	} else {
		log.Printf("✅ binary name: \"%s\"", *name)
		log.Printf("✅ binary output path: \"%s\"", target+"/bin/"+*name)
		log.Printf("✅ maximum optimization? %v", *optimize)
		log.Printf("✅ compiling project...\n")
		compile(includepath, target)
	}
}
