cmake_minimum_required(VERSION 2.8)

project(libstrict)

add_library(libstrict STATIC
  src/strict/all.hh
  src/strict/log.cc
  src/strict/log.hh
  src/strict/strict.cc
  src/strict/strict.hh)

install(TARGETS libstrict DESTINATION include)

include_directories(src/)
