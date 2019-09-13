#include <string>
namespace  {
  #include "Date.h"
}

namespace  {
  #include "Name.h"
}


class Person {
 public:
 Name Name;
 Date BirthDate;
 std::vector<Person> Friends;
 explicit Person();
 private:
 void Generated$Init();
}

