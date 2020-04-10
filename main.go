// Program `cm` is a tiny tool for automating modern C++ projects.
package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/damienstanton/cm/statik"
)

var (
	debug       = flag.Bool("debug", false, "print the wrapped command for inspection")
	name        = flag.String("o", "", "name of the output binary")
	includepath = flag.String("I", "", "path to header files")
	optimize    = flag.Bool("max", false, "maximum optimization")
	std         = flag.String("std", "c++2a", "c++ standard library to use")
	compiler    = flag.String("compiler", "clang++", "c++ compiler to use")
	testMode    = flag.Bool("test", false, "run tests using Catch2")
	initF       = flag.Bool("init", false, "scaffold & .gitkeep the required dirs")
	run         = flag.Bool("run", false, "execute the successfully compiled binary, like go run")
)

func main() {
	log.SetPrefix("╠ ")
	flag.Parse()
	target, err := os.Getwd()
	if err != nil {
		log.Fatal("could not determine current directory (are you in a symlink?)")
	}

	if *name == "" {
		path, err := filepath.Abs(target)
		if err != nil {
			log.Fatalf("path error: %+v", err)
		}
		split := strings.Split(path, "/")
		*name = split[len(split)-1]
	}

	printBanner()
	if *initF {
		err := mkScaffoldDirs()
		if err != nil {
			log.Fatalf("dir write error: %v", err)
		}
		log.Printf("init completed successfully for %s\n", target)
		os.Exit(0)
	}
	if *testMode {
		runTests(target)

	}
	if !*testMode {
		runCompile(target)
	}
}

// runCompile executes the given compiler config
func runCompile(target string, args ...string) {
	binary := target + "/bin/" + *name
	log.Printf("binary name: \"%s\"", *name)
	log.Printf("binary output path: \"%s\"", binary)
	log.Printf("maximum optimization? %v", *optimize)
	log.Printf("compiling project...\n")
	compile(includepath, target, args...)
	wrap(binary, []string{})
}

// runCompile executes the given compiler config (like runCompile), but with extra operations around unit tests
func runTests(target string, args ...string) {
	catchFile := "/catch.hpp"
	hostFile := "/test_main.cpp"

	log.Println("entering test mode...")
	copyTestFramework(catchFile, hostFile, target)
	log.Printf("compiling catch2 and tests (this may take a while)...\n")
	compile(includepath, target, args...)

	testBinary := target + "/tests/" + *name
	log.Printf("running %s tests using catch %s", testBinary, catchVersion)
	out, _ := wrap(testBinary, []string{}) // ignore this error as it just indicates test failures
	log.Println(string(out))

	log.Println("cleaning up test framework...")
	err := os.Remove(target + "/tests/" + *name)
	err = os.Remove(target + "/tests/" + hostFile)
	err = os.Remove(target + "/tests/" + catchFile)
	if err != nil {
		log.Fatalf("cleanup error: %+v", err)
	}
	log.Println("exited test mode")
}
