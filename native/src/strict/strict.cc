#include "strict/strict.hh"

#include <cstdarg>
#include <cstdio>
#include <iostream>

namespace strict {

inline strict::Number InputNumber(const strict::Text &message) {
  std::cout << message << '\n';
  strict::Number number;
  scanf("%lf", &number);
  return number;
}

inline void Log(const strict::Text &message) {
  std::cout << message << '\n';
}

inline void Logf(const char *format, ...) {
  va_list arguments;
  va_start(arguments, format);
  std::vfprintf(stdout, format, arguments);
  va_end(arguments);
}

} // namespace strict