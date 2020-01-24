#include "HelloWorld.h"

#include <string>
#include <vector>
namespace io {
  #include stdio
}


void HelloWorld::Generated$Init() {
 for (auto number = 0; number < 100; number++) {
  io::Printf("%d\n", number);
 }
io::Puts("Hello, World!");
}

std::vector<int> HelloWorld::Range(int begin, int end) {
 std::vector<int> $yield;
 for (auto number = begin; number < end; number++) {
  $yield.push_back(number);
 }

 return $yield;
}

HelloWorld::HelloWorld() {
 Generated$Init();
 {

 }
}
