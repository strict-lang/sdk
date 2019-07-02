#include "strict/log.hh"

#include <iostream>
#include <cstdio>

namespace strict {

// Logs the message followed by a linefeed character.
inline void Log(const strict::Text &message) {
  std::cout << message << '\n';
}

// Formats and logs the message using the passed arguments. The |format| string
// is using the same format as the c-function 'printf' does.
inline void Logf(const strict::Text &format, ...) {
  va_list arguments;
  va_start(arguments, format);
  vfprintf(stdout, format, arguments);
  va_end(arguments);
}

} // namespace strict