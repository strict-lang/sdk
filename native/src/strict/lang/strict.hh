#ifndef STRICT_STRICT_HH_
#define STRICT_STRICT_HH_

#include <string>
#include <vector>

namespace strict {

using Text = std::string;
using Number = double;

template<class T>
using List = std::vector<T>;

// Logs the message followed by a linefeed character.
inline void Log(const strict::Text &message);

// Formats and logs the message using the passed arguments. The |format| string
// is using the same format as the c-function 'printf' does.
inline void Logf(const char *format, ...);

inline strict::Number InputNumber(const strict::Text &message);

} // namespace strict

#endif // STRICT_STRICT_HH_
