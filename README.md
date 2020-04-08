# **C**`ompiler` **M**`anager`

_`cm` is a dumb wrapper around your C++ compiler to make it feel a little more like `go build`_

## Installation

```sh
$ go get github.com/damienstanton/cm
# or, download from the release tab, add to your PATH

```

## Directory Layout Expectations

The tool is naive, and so a certain dir tree is required to find dynamic libraries and tests.
Here is the directory layout for a smoke test program:

```sh
example
├── bin # where executables are dropped
│   └── example
├── lib # where to put shared objects
│   ├── hello_cgo.go
│   ├── libhello.h
│   └── libhello.so
├── src # where to put headers & impls
│   ├── greeting.cpp
│   └── greeting.h
└── tests # where to put unit tests
    └── greeting_test.cpp
```

## Compiling

```console
# suppose we are in the example dir in this repo
$ cm
╔═════════════════════════╗
║ Compiler Manager v0.1.0 ║
╚═════════════════════════╝
╠ 2020/04/08 13:07:39 binary name: "example"
╠ 2020/04/08 13:07:39 binary output path: "/Users/damien/oss/cm/example/bin/example"
╠ 2020/04/08 13:07:39 maximum optimization? false
╠ 2020/04/08 13:07:39 compiling project...
╠ 2020/04/08 13:07:39 checking for shared objects in /Users/damien/oss/cm/example/lib...
╠ 2020/04/08 13:07:39 none found
╠ 2020/04/08 13:07:40 🎉 compilation succeeded with no errors
$ bin/example
Hello, there
$ bin/example Damien
Hello, Damien
```

Nice, right? Didn't have to think of anything. Probably could've just been a zsh alias, but hey, this is more fun. I do intend to expand the feature set (see [Features & TODOs](#features--todos)).

## Testing

`cm` comes with a bundled C++ test framework, [Catch2](https://github.com/catchorg/Catch2). This is embedded in the application binary and is removed when tests pass. All you need to do is `#include "catch.hpp"` and follow the Catch macro/guidelines for testing and the tool does the rest. Neat!

A failing test:

```console
$ cm -test
╔═════════════════════════╗
║ Compiler Manager v0.1.0 ║
╚═════════════════════════╝
╠ 2020/04/08 13:10:12 entering test mode...
╠ 2020/04/08 13:10:12 compiling catch2 and tests (this may take a while)...
╠ 2020/04/08 13:10:12 checking for shared objects in /Users/damien/oss/cm/example/lib...
╠ 2020/04/08 13:10:12 none found
╠ 2020/04/08 13:10:18 🎉 compilation succeeded with no errors
╠ 2020/04/08 13:10:18 running /Users/damien/oss/cm/example/tests/example tests using catch v2.11.3
╠ 2020/04/08 13:10:18 wrap error: exit status 1
╠ 2020/04/08 13:10:18
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
example is a Catch v2.11.3 host application.
Run with -? for options

-------------------------------------------------------------------------------
Greeting with no args
-------------------------------------------------------------------------------
/Users/damien/oss/cm/example/tests/greeting_test.cpp:11
...............................................................................

/Users/damien/oss/cm/example/tests/greeting_test.cpp:12: FAILED:
  REQUIRE( hi("") == "Hello, there." )
with expansion:
  "Hello, there" == "Hello, there."

===============================================================================
test cases: 2 | 1 passed | 1 failed
assertions: 2 | 1 passed | 1 failed


╠ 2020/04/08 13:10:18 cleaning up test framework...
╠ 2020/04/08 13:10:18 exited test mode
```

Once the test is fixed:

```console
$ cm -test
╔═════════════════════════╗
║ Compiler Manager v0.1.0 ║
╚═════════════════════════╝
╠ 2020/04/08 13:11:03 entering test mode...
╠ 2020/04/08 13:11:03 compiling catch2 and tests (this may take a while)...
╠ 2020/04/08 13:11:03 checking for shared objects in /Users/damien/oss/cm/example/lib...
╠ 2020/04/08 13:11:03 none found
╠ 2020/04/08 13:11:09 🎉 compilation succeeded with no errors
╠ 2020/04/08 13:11:09 running /Users/damien/oss/cm/example/tests/example tests using catch v2.11.3
╠ 2020/04/08 13:11:09 ===============================================================================
All tests passed (2 assertions in 2 test cases)


╠ 2020/04/08 13:11:09 cleaning up test framework...
╠ 2020/04/08 13:11:09 exited test mode
```

### Geez, tests are really slow

[I know.](#features--todos)

## Help

See `cm -help` for options, all of which are optional:

```console
Usage of cm:

```

## Considerations

- This is very poorly tested, there are likely lots of edge cases
- This toy is made for my own purposes iterating small C++ programs (e.g. code interviews, competitive exercises, etc.)
- I would not use it for anything "serious", but feel free to use it as a starting point for a better automation tool

## Features & TODOs

Not an exhaustive list, will probably use the Issues/Project tab if I end up using this more broadly.

- [x] C++ compilation automation
- [x] C++ shared object import
- [x] C++ unit test automation
- [ ] C++ benchmark automation
- [ ] Unit tests (for `cm` itself)
- [ ] JSON config
- [ ] Make test compilation less brutally slow (linking against already-compiled test main, should be easy)

## Contributing

Sure, just fill out an issue/PR.

## License

© 2020 Damien Stanton

Distributed freely under the Apache 2.0 license. See LICENSE file for details.
