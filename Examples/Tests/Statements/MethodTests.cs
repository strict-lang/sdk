using NUnit.Framework;
using Strict.Statements;
using Text = Strict.Statements.Text;

namespace Strict.Tests.Statements
{
	public class MethodTests : TestWithContext
	{
		[Test]
		public void MethodsMustStartWithLowercaseLetter()
		{
			Assert.Throws<NamedStatement.MemberOrMethodMustStartWithLowerCaseLetter>(
				() => new Method(Base.Context, Base.Void, "SayHi"));
		}

		[Test]
		public void MethodsCannotBeCreatedInBaseContext()
		{
			Assert.Throws<Method.CannotCreateMethodInBaseContext>(
				() => new Method(Base.Context, Base.Void, "sayHi"));
		}

		[Test]
		public void CreateMethod()
		{
			using (var method = new Method(TestContext, Base.Void, "sayHi"))
			{
				Assert.That(method.ReturnType, Is.EqualTo(Base.Void));
				Assert.That(method.Name, Is.EqualTo("sayHi"));
				Assert.That(method.ToString(), Is.EqualTo("void TestContext.sayHi()"));
			}
		}

		[Test]
		public void CreateMethodsWithSameNameButDifferentTypes()
		{
			var numberParameter = new Parameter(Base.Number, "value");
			var textParameter = new Parameter(Base.Text, "value");
			using (var method1 = new Method(TestContext, Base.Void, "test", numberParameter))
			using (var method2 = new Method(TestContext, Base.Void, "test", textParameter))
			{
				Assert.That(method1.Parameters[0].ReturnType, Is.EqualTo(Base.Number));
				Assert.That(method2.Parameters[0].ReturnType, Is.EqualTo(Base.Text));
				Assert.That(TestContext.GetMethod("test", new Number(5)), Is.EqualTo(method1));
				Assert.That(TestContext.GetMethod("test", new Text("")), Is.EqualTo(method2));
			}
		}
		
		[Test]
		public void DisposingTypeKillsAllMethodsDefinedInItToo()
		{
			var abcType = new Type(TestContext, "abc");
			using (new Method(TestContext, abcType, "test"))
			{
				Assert.That(TestContext.GetMethod("test", null), Is.Not.Null);
				abcType.Dispose();
				Assert.Throws<Context.MethodNotFound>(() => TestContext.GetMethod("test", null));
			}
		}

		[Test]
		public void OverridingMethodIsDoneManually()
		{
			var aType = new Type(TestContext, "a");
			var abType = new Type(TestContext, "ab");
			abType.Add(new Member("a", new Value(aType)));
			using (var parentMethod = new Method(TestContext, Base.Void, "test",
				new Parameter(aType, "value")))
			using (
				var overriddenMethod = new Method(TestContext, Base.Void, "test",
					new Parameter(abType, "value")))
			{
				//overriddenMethod.Add(new MethodCall(parentMethod,
				//TODO: we still need a way to access members:	new GetMember(overriddenMethod.Parameters[0], "a")));
				//TODO: this kind of test is better written as a file and then checked here!
			}
		}

		//TODO: add missing Condition.cs and Loop.cs with tests first
		//TODO: Add method file saving and loading: https://strict.fogbugz.com/default.asp?W10
	}
}