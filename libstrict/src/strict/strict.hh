#ifndef STRICT_STRICT_HH_
#define STRICT_STRICT_HH_

#include <string>
#include <vector>

namespace strict {

using Text = std::string;
using Number = double;

template<class T>
using List = std::vector<T>;

} // namespace strict

#endif // STRICT_STRICT_HH_
