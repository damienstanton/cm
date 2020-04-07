#include "greeting.hpp"

fn hi(std::string name)->std::string {
    if (name == "")
        return "Hello, there";
    return "Hello, " + name;
}

// TODO: generate a Go string from cgo to test dylib stuff

fn main(int argc, char* argv[])->int {
    if (argc < 2)
        std::cout << hi("") << std::endl;
    else
        std::cout << hi(argv[1]) << std::endl;
    return 0;
}