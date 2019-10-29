#ifndef STRICT_TESTING_H_
#define STRICT_TESTING_H_

#include <string>
#include <functional>

namespace strict {
namespace testing {

// Location within the source-code at which some test-assertion is written.
// Used to report the specific position of a failed test-assertion.
struct CodePosition final {
  int line_index;
  int column;
};

// Member of the tested class, that is tested by a MethodTesting.
// Captures the full description of a method, which contains name,
// return type, parameters and the position of the methods
// declaration, which is used to provide richer failure reports.
struct TestedMethod final {
  std::string name;
  std::string descriptor;
  std::string return_type_name;
  CodePosition declaration_position;
};

// The full report created by a testing after all of the classes
// methods, with tests, have been tested. If the vector of test
// reports is empty, all tests have been successful.
struct TestReport final {
  std::vector<MethodTestReport> method_tests;
};

// Report created by a method testing, that captures the possible
// errors encountered while running the method test and a description
// of the method that has been tested. If the entries are empty, the
// test has been run successful.
struct MethodTestReport final {
  struct Entry final {
    std::string message;
    CodePosition position;
  };
  TestedMethod method;
  std::vector<Entry> entries;
};

class Testing final {
 public:
  explicit Testing(const std::string &type_name);

  TestReport CreateTestReport();
 protected:
  void AppendMethodTestReport(const MethodTestReport &report);
 private:
  std::string type_name_;
  std::vector<MethodTestReport> method_reports;
};

class MethodTesting final {
 public:
  explicit MethodTesting(Testing *testing, const TestedMethod &method);

  ~MethodTesting();

  void ReportFailedAssertion(
    const CodePosition &position, const std::string &message);

 private:
  Testing *testing_;
  TestedMethod method_;
  std::vector<MethodTestReport::Entry> entries_;
};

} // namespace testing
} // namespace strict

#endif // STRICT_TESTING_H_