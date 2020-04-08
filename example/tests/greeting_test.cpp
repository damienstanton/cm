#include "../src/greeting.h"

#include "catch.hpp"

auto hi(std::string name) -> std::string {
    if (name == "")
        return "Hello, there";
    return "Hello, " + name;
}

TEST_CASE("Greeting with no args", "[greeting]") {
    REQUIRE(hi("") == "Hello, there");
}

TEST_CASE("Greeting with an arg", "[greeting]") {
    REQUIRE(hi("Damien") == "Hello, Damien");
}
