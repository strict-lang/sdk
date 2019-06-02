using NUnit.Framework;
using Strict.Statements;

namespace Strict.CollapsedGrammar.Tests
{
	public class TestWithParser
	{
		[SetUp]
		public void CreateParser()
		{
			TestContext = new Context(Base.Context, "Test");
			testMethod = new Method(TestContext, Base.Void, "testMethod");
			parser = new Parser(testMethod.Scope);
		}

		protected Context TestContext { get; private set; }
		protected Parser parser;
		private Method testMethod;

		[TearDown]
		public void DestroyTestContext()
		{
			testMethod.Dispose();
      TestContext.Dispose();
		}
	}
}