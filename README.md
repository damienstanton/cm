# cm

## **C**`ompiler` **M**`anager`

### a dumb wrapper around your C++ compiler to make it feel a little more like `go build`

#### Installation

```sh
$ go get github.com/damienstanton/cm
# installs to your Go bin dir
```

#### Help

See `cm -help` for options:

```sh
Usage of ./cm:
  -I string
        path to header files
  -compiler string
        c++ compiler to use (default "g++")
  -max
        maximum optimization
  -o string
        name of the output binary
  -std string
        c++ standard library to use (default "c++2a")
```

#### Considerations

- This toy is made for my own purposes iterating small C++ programs (e.g. code interviews, competitive exercises, etc.)
- I would not use it for anything "serious", but feel free to use it as a starting point for a better automation tool.

#### TODO

- [ ] Unit tests
- [ ] JSON config

#### License

© 2020 Damien Stanton

Distributed under the Apache 2.0 license; See LICENSE for details.
