using NUnit.Framework;
using Strict.Parsing;
using Strict.Tests.Statements;

namespace Strict.Tests.Parsing
{
	public class MethodFileParserTests : TestWithContext
	{
		[Test]
		public void MethodFilenameMustContainBrackets()
		{
			Assert.Throws<MethodFileParser.MethodFilenameMustContainBrackets>(
				() => new MethodFileParser(TestContext, "TestContext\\invalidEmpty.method"));
		}

		[Test]
		public void MethodFileParserNeedsFilenameEndingWithMethod()
		{
			Assert.Throws<TypeFileParser.FilenameNeedsToBeValidAndEndWith>(
				() => new MethodFileParser(TestContext, "()"));
			Assert.Throws<TypeFileParser.FilenameNeedsToBeValidAndEndWith>(
				() => new MethodFileParser(TestContext, "abc()"));
			Assert.Throws<TypeFileParser.FilenameNeedsToBeValidAndEndWith>(
				() => new MethodFileParser(TestContext, "C:\\abc()"));
		}

		[Test]
		public void FilePathMustBeInADirectoryForTheContext()
		{
			Assert.Throws<TypeFileParser.FileMustBeInFolderForItsContext>(
				() => new MethodFileParser(TestContext, "invalidEmpty().method"));
		}

		[Test]
		public void DifferentStartDirectoryIsNotAllowedForTheContext()
		{
			Assert.Throws<TypeFileParser.FilePathMustBeBelowStrictExecutableFolder>(
				() => new MethodFileParser(TestContext, "c:\\invalidEmpty().method"));
		}

		[Test]
		public void ContentAndFilePathMustMatch()
		{
			Assert.Throws<TypeFileParser.ContextMustMatchPath>(
				() => new MethodFileParser(TestContext, "OtherContext\\invalidEmpty().method"));
		}


		[Test]
		public void EmptyMethodFileIsNotAllowed()
		{
			Assert.Throws<TypeFileParser.ContentCannotBeEmpty>(
				() => new MethodFileParser(TestContext, "TestContext\\invalidEmpty().method"));
		}

		/*TODO
		[Test]
		public void TypeFileMustStartWithUpperLetter()
		{
			Assert.Throws<Type.TypesMustStartWithUppercaseLetter>(
				() => new TypeFileParser(TestContext, "TestContext\\invalidFilename.type"));
		}
		
		[Test]
		public void MemberMustStartWithLowerLetter()
		{
			Assert.Throws<NamedStatement.MemberOrMethodMustStartWithLowerCaseLetter>(
				() => new TypeFileParser(TestContext, "TestContext\\InvalidName.type"));
		}
		
		[Test]
		public void LoadSimpleTypeFile()
		{
			var simpleType = new TypeFileParser(TestContext, "TestContext\\Simple.type");
			Assert.That(simpleType.Name, Is.EqualTo("Simple"));
			Assert.That(simpleType.Members, Has.Count.EqualTo(2));
			Assert.That(simpleType.Members[0].Name, Is.EqualTo("number"));
			Assert.That(simpleType.Members[1].Name, Is.EqualTo("text"));
		}
		
		[Test]
		public void LoadTestTypeFile()
		{
			var simpleType = new TypeFileParser(TestContext, "TestContext\\Test.type");
			Assert.That(simpleType.Name, Is.EqualTo("Test"));
			Assert.That(simpleType.Members, Has.Count.EqualTo(4));
			Assert.That(simpleType.Members[0].Name, Is.EqualTo("first"));
			Assert.That(simpleType.Members[1].Name, Is.EqualTo("second"));
			Assert.That(simpleType.Members[2].Name, Is.EqualTo("text"));
			Assert.That(simpleType.Members[3].Name, Is.EqualTo("pi"));
		}

		[Test]
		public void TypeFileParserNeedsValidContext()
		{
			Assert.Throws<Type.CannotCreateTypeInBaseContext>(
				() => new TypeFileParser(Base.Context, "", new[] { "number a" }));
		}

		[Test]
		public void TypeFileParserNeedsValidName()
		{
			Assert.Throws<Type.TypesMustStartWithUppercaseLetter>(
				() => new TypeFileParser(TestContext, "", new[] { "number a" }));
		}

		[Test]
		public void TypeFileParserContentCannotBeEmpty()
		{
			Assert.Throws<TypeFileParser.ContentCannotBeEmpty>(
				() => new TypeFileParser(TestContext, "Abc", new string[0]));
		}

		[Test]
		public void TypeCanBeCreatedWithContextAndMemberLines()
		{
			using (var type = new TypeFileParser(TestContext, "Abc", new[] { "Number value", "text" }))
			{
				Assert.That(type.Context.Name, Is.EqualTo("TestContext"));
				Assert.That(type.Name, Is.EqualTo("Abc"));
				Assert.That(type.Members, Has.Count.EqualTo(2));
				Assert.That(type.Members[0].Name, Is.EqualTo("value"));
				Assert.That(type.Members[0].ReturnType, Is.EqualTo(Base.Number));
				Assert.That(type.Members[1].Name, Is.EqualTo("text"));
				Assert.That(type.Members[1].ReturnType, Is.EqualTo(Base.Text));
			}
		}

		[Test]
		public void EmptyLinesAreNotAllowedInTypeFile()
		{
			Assert.Throws<TypeFileParser.EmptyLinesAreNotAllowedInTypeFile>(
				() => new TypeFileParser(TestContext, "Abc", new[] { "number", "" }));
		}

		[Test]
		public void TypeMembersCanBeSetToInitialValue()
		{
			using (var type = new TypeFileParser(TestContext, "Point", new[] { "x = 1", "y = 2" }))
			{
				Assert.That(type.Name, Is.EqualTo("Point"));
				Assert.That(type.Members, Has.Count.EqualTo(2));
				Assert.That(type.Members[0].Name, Is.EqualTo("x"));
				Assert.That(type.Members[0].ReturnType, Is.EqualTo(Base.Number));
				Assert.That(type.Members[0].Value, Is.EqualTo(new Number(1)));
				Assert.That(type.Members[1].Name, Is.EqualTo("y"));
				Assert.That(type.Members[1].ReturnType, Is.EqualTo(Base.Number));
				Assert.That(type.Members[1].Value, Is.EqualTo(new Number(2)));
			}
		}

		[Test]
		public void InvalidMemberAssignmentSyntaxIsNotAllowed()
		{
			Assert.Throws<Value.TypeParsingIsNotYetSupported>(
				() => new TypeFileParser(TestContext, "AssignFromOtherMember", new[] { "x = y" }));
			Assert.Throws<TypeFileParser.InvalidTypeMemberSyntax>(
				() => new TypeFileParser(TestContext, "AssignTwice", new[] { "x = 1 = 2" }));
			Assert.Throws<TypeFileParser.InvalidTypeMemberSyntax>(
				() => new TypeFileParser(TestContext, "InvalidSyntax", new[] { "x 1 2" }));
			Assert.Throws<Context.TypeNotFound>(
				() => new TypeFileParser(TestContext, "JustNumber", new[] { "2385" }));
		}
	*/
	}
}