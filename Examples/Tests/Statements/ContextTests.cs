using System.Linq;
using NUnit.Framework;
using Strict.Statements;
using Text = Strict.Statements.Text;

namespace Strict.Tests.Statements
{
	public class ContextTests : TestWithContext
	{
		[Test]
		public void ContextNameShouldNotContainDots()
		{
			Assert.Throws<Context.NameShouldNotContainDots>(() => new Context(Base.Context, "a.b"));
		}
		
		[Test]
		public void CreatingAContextWithAnEmptyNameIsNotAllowed()
		{
			Assert.Throws<Context.ContextNameCannotBeEmpty>(() => new Context(Base.Context, ""));
		}
		
		[Test]
		public void TypeFileParserNameMustStartWithUpperCaseLetter()
		{
			Assert.Throws<Context.NameMustStartWithUpperCaseLetter>(
				() => new Context(Base.Context, "abc"));
		}

		[Test]
		public void EachContextNeedsAUniqueName()
		{
			Assert.Throws<Context.ContextAlreadyExists>(() => new Context(Base.Context, "TestContext"));
		}

		[Test]
		public void CannotCreateContextWhenTypeWithSameNameAlreadyExists()
		{
			new Type(TestContext, "abc");
			Assert.Throws<Context.TypeWithThisContextNameAlreadyExists>(
				() => new Context(TestContext, "Abc"));
		}

		[Test]
		public void NestedContextCanHaveTheSameName()
		{
			using (var nestedContext = new Context(TestContext, "TestContext"))
				Assert.That(nestedContext.ToString(), Is.EqualTo("TestContext.TestContext"));
		}

		[Test]
		public void NumberTypeIsAlwaysTheSame()
		{
			Assert.That(Base.Context.GetType("number"), Is.EqualTo(Base.Number));
		}

		[Test]
		public void BaseContextHasBuildInTypes()
		{
			Assert.That(Base.Context.NumberOfTypes, Is.EqualTo(7));
			Assert.That(Base.Number.Name, Is.EqualTo("number"));
			Assert.That(Base.Text.Name, Is.EqualTo("text"));
			Assert.That(Base.Bool.Name, Is.EqualTo("bool"));
			Assert.That(Base.List.Name, Is.EqualTo("list"));
			Assert.That(Base.Map.Name, Is.EqualTo("map"));
			Assert.That(Base.Anything.Name, Is.EqualTo("anything"));
			Assert.That(Base.Void.Name, Is.EqualTo("void"));
		}

		[Test]
		public void EachTypeIsInAContext()
		{
			Assert.That(TestContext.NumberOfTypes, Is.EqualTo(0));
			new Type(TestContext, "test");
			Assert.That(TestContext.NumberOfTypes, Is.EqualTo(1));
		}

		[Test]
		public void CreateInnerContext()
		{
			using (var method = new Method(TestContext, Base.Number, "testAdd",
				new Parameter(Base.Number, "first"), new Parameter(Base.Number, "second")))
			{
				var memberInMethod = new Member("testNumber", new Number(0));
				method.Add(memberInMethod);
				using (var innerMethod = new Method(method.Scope, Base.Number, "innerAdd"))
				{
					Assert.That(innerMethod.Scope.AllAccessibleMembersRecursively.Last(), Is.EqualTo(memberInMethod));
					Assert.That(innerMethod.Scope.GetMethod("innerAdd", null), Is.EqualTo(innerMethod));
					Assert.That(innerMethod.Scope.GetMethod("testAdd", new Number(1), new Number(2)),
						Is.EqualTo(method));
				}
				Assert.Throws<Context.MethodNotFound>(() => method.Scope.GetMethod("innerAdd", null));
			}
		}

		[Test]
		public void FindMethod()
		{
			using (var method = new Method(TestContext, Base.Number, "testAdd",
				new Parameter(Base.Number, "first"), new Parameter(Base.Number, "second")))
			{
				Assert.That(method.Scope.AllAccessibleMembersRecursively, Is.Empty);
				var memberInMethod = new Member("testNumber", new Number(0));
				method.Add(memberInMethod);
				Assert.That(method.Scope.AllAccessibleMembersRecursively, Is.Not.Empty);
				Assert.That(method.Scope.GetMethod("testAdd", new Number(1), new Number(2)),
					Is.EqualTo(method));
				Assert.Throws<Context.MethodNotFound>(
					() => method.Scope.GetMethod("testAdd", new Number(1)));
				Assert.Throws<Context.MethodNotFound>(
					() => method.Scope.GetMethod("testAdd", new Text(""), new Text("")));
			}
		}

		[Test]
		public void CreateTypeInNestedContext()
		{
			var allTypesCount = TestContext.GetAllTypesRecursively().Count();
			new Type(TestContext, "outerType");
			Assert.That(TestContext.GetAllTypesRecursively().Count(), Is.EqualTo(allTypesCount + 1));
			using (var nestedContext = new Context(TestContext, "Nested"))
			{
				Assert.That(nestedContext.NumberOfTypes, Is.EqualTo(0));
				new Type(nestedContext, "innerType");
				Assert.That(nestedContext.NumberOfTypes, Is.EqualTo(1));
			}
			Assert.That(TestContext.GetAllTypesRecursively().Count(), Is.EqualTo(allTypesCount + 1));
		}

		[Test]
		public void RemoveType()
		{
			var testTypes = TestContext.NumberOfTypes;
			using (new Type(TestContext, "testType"))
				Assert.That(TestContext.NumberOfTypes, Is.EqualTo(testTypes + 1));
			Assert.That(TestContext.NumberOfTypes, Is.EqualTo(testTypes));
		}

		[Test]
		public void GetContext()
		{
			Assert.That(Base.GetContext("TestContext"), Is.EqualTo(TestContext));
			Assert.Throws<Context.ChildContextNotFound>(() => Base.GetContext("Unknown"));
		}
	}
}