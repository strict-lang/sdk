#ifndef STRICT_LOG_HH_
#define STRICT_LOG_HH_

#include "strict/strict.hh"

namespace strict {

// Logs the message followed by a linefeed character.
inline void Log(const strict::Text &message);

// Formats and logs the message using the passed arguments. The |format| string
// is using the same format as the c-function 'printf' does.
inline void Logf(const strict::Text &format, ...);

} // namespace strict

#endif // STRICT_LOG_HH_