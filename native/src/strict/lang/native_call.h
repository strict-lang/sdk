#ifndef STRICT_NATIVE_CALL_H_
#define STRICT_NATIVE_CALL_H_

#define STRICT_NATIVE_CALL(Name, ReturnType) \
  NATIVE_##returnType Name

namespace strict {
namespace runtime {

STRICT_NATIVE_CALL(logf, void) (const char *message, ...) {}

STRICT_NATIVE_CALL(log, void) (const std::string &message) {
  std::cout << message << '\n';
}

STRICT_NATIVE_CALL(flushStdout, void) () {
  std::clout.flush();
}

} // namespace runtime
} // namespace strict


#endif // STRICT_NATIVE_CALL_H_