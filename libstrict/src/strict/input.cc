#include "strict/input.hh"

#include <stdio.h>
#include <iostream>

namespace strict {

inline strict::Number InputNumber(const strict::Text &message) {
  std::cout << message << '\n';
  strict::Number number;
  scanf("%lf", &number);
  return number;
}

} // namespace strict
