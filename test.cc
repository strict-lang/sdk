namespace io {
  #include "stdio.h"
}

void greet() {
  io::puts("Hello, World!");
}

int main(int argc, char **argv) {
  greet();
}
