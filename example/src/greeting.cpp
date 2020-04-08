#include "greeting.h"

extern "C" {
#include "../lib/libhello.h"
}

auto hi(std::string name) -> std::string {
    if (name == "")
        return "Hello, there";
    return "Hello, " + name;
}

auto main(int argc, char* argv[]) -> int {
    if (argc < 2)
        std::cout << hi("") << std::endl;
    else
        std::cout << hi(argv[1]) << std::endl;
	
	std::cout << (std::string)HiFromGo() << std::endl;
	return 0;
}
