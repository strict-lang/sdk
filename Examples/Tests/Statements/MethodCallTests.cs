using NUnit.Framework;
using Strict.Statements;

namespace Strict.Tests.Statements
{
	public class MethodCallTests : TestWithContext
	{
		[SetUp]
		public void CreateMethodCall()
		{
			var method = new Method(TestContext, Base.Void, "sayHi");
			call = new MethodCall(method);
		}

		private MethodCall call;

		[TearDown]
		public void RemoveMethod()
		{
			call.Method.Dispose();
		}

		[Test]
		public void CheckValues()
		{
			Assert.That(call.ReturnType, Is.EqualTo(Base.Void));
			Assert.That(call.Method.Name, Is.EqualTo("sayHi"));
			Assert.That(call.Arguments, Is.Empty);
		}

		[Test]
		public void MethodCallToString()
		{
			Assert.That(call.ToString(), Is.EqualTo("sayHi()"));
		}

		[Test]
		public void MethodCallEquals()
		{
			Assert.That(call, Is.EqualTo(call));
			Assert.That(call, Is.EqualTo(new MethodCall(call.Method)));
			Assert.That(call, Is.Not.EqualTo(call.ReturnType));
		}
	}
}