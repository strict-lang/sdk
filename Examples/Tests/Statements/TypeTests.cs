using NUnit.Framework;
using Strict.Statements;

namespace Strict.Tests.Statements
{
	public class TypeTests : TestWithContext
	{
		[Test]
		public void CreateBuildInTypeIsNotAllowed()
		{
			Assert.That(Base.Number, Is.EqualTo(Base.Number));
			Assert.Throws<Type.CannotCreateAlreadyExistingType>(() => new Type(Base.Context, "number"));
		}

		[Test]
		public void TypesCannotBeCreatedInBaseScope()
		{
			Assert.Throws<Type.CannotCreateTypeInBaseContext>(
				() => new Type(Base.Context, "simple"));
		}

		[Test]
		public void TypesMustStartWithUppercaseLetter()
		{
			Assert.Throws<Type.TypesMustStartWithLowercaseLetter>(
				() => new Type(TestContext, "InvalidName"));
		}

		[Test]
		public void CannotCreateTypeWhenContextWithSameNameAlreadyExists()
		{
			new Context(TestContext, "Abc");
			Assert.Throws<Context.TypeWithThisContextNameAlreadyExists>(
				() => new Type(TestContext, "abc"));
		}
		
		[Test]
		public void DisposingABaseTypeIsNotAllowed()
		{
			Assert.Throws<Type.BaseContextTypesShouldNotBeDisposed>(() => Base.Number.Dispose());
		}

		[Test]
		public void CreateSimpleType()
		{
			var simple = new Type(TestContext, "simple");
			Assert.That(simple.Name, Is.EqualTo("simple"));
			Assert.That(simple.UsedBy, Is.Empty);
		}

		[Test]
		public void CreatingTheSameTypeTwiceIsNotAllowed()
		{
			var simpleType1 = new Type(TestContext, "simple");
			Assert.That(simpleType1.Name, Is.EqualTo("simple"));
			Assert.Throws<Type.CannotCreateAlreadyExistingType>(() => new Type(TestContext, "simple"));
		}

		[Test]
		public void TypeToString()
		{
			var simple = new Type(TestContext, "simple");
			Assert.That(simple.ToString(), Is.EqualTo("TestContext.simple"));
		}

		[Test]
		public void UsingATypeAutomaticallyAddsItToTheUsedByList()
		{
			var simpleType = new Type(TestContext, "simple");
			Assert.That(simpleType.UsedBy, Is.Empty);
			var value = new Value(simpleType);
			Assert.That(simpleType.UsedBy, Has.Count.EqualTo(1));
			new Member("usingSimple", value);
			Assert.That(simpleType.UsedBy, Has.Count.EqualTo(2));
			new Method(TestContext, simpleType, "returnSimple");
			Assert.That(simpleType.UsedBy, Has.Count.EqualTo(3));
		}
	}
}