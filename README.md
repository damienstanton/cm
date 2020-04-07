# **C**`ompiler` **M**`anager`

_`cm` is a dumb wrapper around your C++ compiler to make it feel a little more like `go build`_

## Installation

```sh
$ go get github.com/damienstanton/cm
# or, download from the release tab and add to your PATH

```

## Directory Layout Expectations

The tool is naive, and so a certain dir tree is required to find dynamic libraries and tests.

- `src`: where non-test C++ header and impl files live
- `tests`: where, uh, the tests live
- `bin`: where binaries are dropped.

Here is the directory layout for a smoke test program:

```sh
example
├── bin
├── src
│   ├── greeting.cpp
│   └── greeting.hpp
└── tests
    └── greeting_test.cpp
```

## Compiling

```sh
# suppose we are in the example dir in this repo
$ cm .
2020/04/07 14:24:14 ✅ binary name: "example"
2020/04/07 14:24:14 ✅ binary output path: "./bin/example"
2020/04/07 14:24:14 ✅ maximum optimization? false
2020/04/07 14:24:14 ✅ compiling project...
2020/04/07 14:24:14 🎉 compilation succeeded with no errors
$ bin/example
Hello, there
$ bin/example Damien
Hello, Damien
```

Nice, right? Didn't have to think of anything. Probably could've just been a zsh alias, but hey, this is more fun. I do intend to expand the feature set (see [Features & TODOs](#Features-&-TODOs)).

## Testing

`cm` comes with a bundled C++ test framework, [Catch2](https://github.com/catchorg/Catch2). This is embedded in the application binary and is removed when tests pass. All you need to do is `#include "catch.hpp"` and follow the Catch macro/guidelines for testing and the tool does the rest. Neat!

A failing test:

```console
$ cm -test .
2020/04/07 14:35:35 🔍 entering test mode...
2020/04/07 14:35:35 🔍 compiling catch2 and tests (this may take a while)...
2020/04/07 14:35:42 🎉 compilation succeeded with no errors
2020/04/07 14:35:42 🔍 running ./tests/example tests using catch v2.11.3
2020/04/07 14:35:42 test failure!
2020/04/07 14:35:42 💥 wrap error: exit status 1
2020/04/07 14:35:42
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
example is a Catch v2.11.3 host application.
Run with -? for options

-------------------------------------------------------------------------------
Greeting works as expected
-------------------------------------------------------------------------------
tests/greeting_test.cpp:11
...............................................................................

tests/greeting_test.cpp:12: FAILED:
  REQUIRE( hi("") == "Hello, there." )
with expansion:
  "Hello, there" == "Hello, there."

===============================================================================
test cases: 1 | 1 failed
assertions: 1 | 1 failed


2020/04/07 14:35:42 🔍 cleaning up test framework...
2020/04/07 14:35:42 🔍 exited test mode
```

Once the test is fixed:

```console
$ cm -test .
2020/04/07 14:37:34 🔍 entering test mode...
2020/04/07 14:37:34 🔍 compiling catch2 and tests (this may take a while)...
2020/04/07 14:37:40 🎉 compilation succeeded with no errors
2020/04/07 14:37:40 🔍 running ./tests/example tests using catch v2.11.3
2020/04/07 14:37:40 ===============================================================================
All tests passed (2 assertions in 1 test case)


2020/04/07 14:37:40 🔍 cleaning up test framework...
2020/04/07 14:37:40 🔍 exited test mode
```

### Geez, tests are really slow

[I know.](#Features-&-TODOs)

## Help

See `cm -help` for options, all of which are optional:

```console
Usage of cm:
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
  -test
    	run tests using Catch2
```

## Considerations

- This is very poorly tested, there are likely lots of edge cases
- This toy is made for my own purposes iterating small C++ programs (e.g. code interviews, competitive exercises, etc.)
- I would not use it for anything "serious", but feel free to use it as a starting point for a better automation tool

## Features & TODOs

Not an exhaustive list, will probably use the Issues/Project tab if I end up using this more broadly.

- [x] C++ compilation automation
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
