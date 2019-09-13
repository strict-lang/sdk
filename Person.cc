#include "Person.h"

#include <string>
namespace  {
  #include "Date.h"
}

namespace  {
  #include "Name.h"
}


Person::Person(Name name, Date birthDate) {
 this->Name = name;
 this->BirthDate = birthDate;
}

