# cm

## **C**`ompiler` **M**`anager`

### a dumb wrapper around your C++ compiler to make it feel a little more like `go build`

#### Installation

```sh
$ go get github.com/damienstanton/cm
# installs to your Go bin dir
```

#### Usage

```sh
# given a hello.h and hello.cpp file in a directory `hello_world`
# and expecting our C++ program to print "OK":

$ cm .
2020/04/06 15:11:37 ✅ binary name: "hello_world"
2020/04/06 15:11:37 ✅ binary output path: "./hello_world"
2020/04/06 15:11:37 ✅ maximum optimization? false
2020/04/06 15:11:37 ✅ compiling project...
2020/04/06 15:11:38 🎉 compilation succeeded with no errors

$ ./hello_world
OK
```

Nice, right? Didn't have to think of anything. Probably could've just been a zsh alias, but hey, this is more fun. I do
intend to expand the feature set to include automated C++ tests and benches, and JSON configuration.

#### Help

See `cm -help` for all options:

```sh
Usage of ./cm:
  -I string
        path to header files (optional, default is target dir)
  -compiler string
        c++ compiler to use (default "g++")
  -max
        maximum optimization (optional, default is off)
  -o string
        name of the output binary (optional, default is name of target dir)
  -std string
        c++ standard library to use (default "c++2a")
```

#### Considerations

- This toy is made for my own purposes iterating small C++ programs (e.g. code interviews, competitive exercises, etc.)
- I would not use it for anything "serious", but feel free to use it as a starting point for a better automation tool.

#### TODO

- [ ] Unit tests (for `cm` itself)
- [ ] JSON config
- [ ] C++ unit test automation
- [ ] C++ benchmark automation

#### License

© 2020 Damien Stanton

Distributed under the Apache 2.0 license; See LICENSE for details.
