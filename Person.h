#pragma once

#include <string>
namespace  {
  #include "Date.h"
}

namespace  {
  #include "Name.h"
}


class Person {
 public:
  explicit Person();
  explicit Person(Name name, Date birthDate);

  Name Name;
  Date BirthDate;
  std::vector<Person> Friends;

 private:
  void Generated$Init();
};

