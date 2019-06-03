using NUnit.Framework;
using Strict.Statements;

namespace Strict.Tests.Statements
{
	/// <summary>
	/// Helps testing strict statements in a local test context as a child of Base.Context, which
	/// cannot be directly used to add members and members as they would never be cleaned up. The
	/// context used here is properly cleaned up in TearDown. Also provides a TestMethod for members.
	/// </summary>
	public class TestWithContext
	{
		[SetUp]
		public void CreateTestContext()
		{
			TestContext = new Context(Base.Context, "TestContext");
			TestMethod = new Method(TestContext, Base.Void, "testMethod");
		}

		protected Context TestContext { get; private set; }
		protected Method TestMethod { get; private set; }

		[TearDown]
		public void DestroyTestContext()
		{
			TestMethod.Dispose();
			foreach (var type in TestContext.GetAllTypesRecursively())
				Assert.That(type.RegisteredMethods, Is.Empty);
			TestContext.Dispose();
			Assert.That(TestContext.ChildContexts, Is.Empty);
			Assert.That(TestContext.NumberOfTypes, Is.EqualTo(0));
		}
	}
}