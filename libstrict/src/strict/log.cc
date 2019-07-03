#include "strict/log.hh"

#include <stdarg.h>
#include <stdio.h>

#include <iostream>

namespace strict {

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