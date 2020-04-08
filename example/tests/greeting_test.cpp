#include "../src/greeting.h"

#include "catch.hpp"

auto hi(std::string name) -> std::string {
    if (name == "")
        return "Hello, there";
    return "Hello, " + name;
}

TEST_CASE("Greeting works as expected", "[greeting]") {
    REQUIRE(hi("") == "Hello, there");
    REQUIRE(hi("Damien") == "Hello, Damien");
}