using NUnit.Framework;
using Strict.Parsing;
using Strict.Statements;
using Strict.Tests.Statements;

namespace Strict.Tests.Parsing
{
	public class TypeFileParserTests : TestWithContext
	{
		[Test]
		public void TypeFileParserNeedsFilenameEndingWithType()
		{
			Assert.Throws<TypeFileParser.FilenameNeedsToBeValidAndEndWith>(
				() => new TypeFileParser(TestContext, ""));
			Assert.Throws<TypeFileParser.FilenameNeedsToBeValidAndEndWith>(
				() => new TypeFileParser(TestContext, "abc"));
			Assert.Throws<TypeFileParser.FilenameNeedsToBeValidAndEndWith>(
				() => new TypeFileParser(TestContext, "C:\\abc"));
		}

		[Test]
		public void FilePathMustBeInADirectoryForTheContext()
		{
			Assert.Throws<TypeFileParser.FileMustBeInFolderForItsContext>(
				() => new TypeFileParser(TestContext, "invalidEmpty.type"));
		}

		[Test]
		public void DifferentStartDirectoryIsNotAllowedForTheContext()
		{
			Assert.Throws<TypeFileParser.FilePathMustBeBelowStrictExecutableFolder>(
				() => new TypeFileParser(TestContext, "c:\\invalidEmpty.type"));
		}

		[Test]
		public void ContentAndFilePathMustMatch()
		{
			Assert.Throws<TypeFileParser.ContextMustMatchPath>(
				() => new TypeFileParser(TestContext, "OtherContext\\invalidEmpty.type"));
		}

		[Test]
		public void EmptyTypeFileIsNotAllowed()
		{
			Assert.Throws<TypeFileParser.ContentCannotBeEmpty>(
				() => new TypeFileParser(TestContext, "TestContext\\invalidEmpty.type"));
		}
		
		[Test]
		public void TypeFileMustStartWithUpperLetter()
		{
			Assert.Throws<Type.TypesMustStartWithLowercaseLetter>(
				() => new TypeFileParser(TestContext, "TestContext\\InvalidFilename.type"));
		}
		
		[Test]
		public void MemberMustStartWithLowerLetter()
		{
			Assert.Throws<NamedStatement.MemberOrMethodMustStartWithLowerCaseLetter>(
				() => new TypeFileParser(TestContext, "TestContext\\invalidName.type"));
		}
		
		[Test]
		public void LoadSimpleTypeFile()
		{
			var simpleType = new TypeFileParser(TestContext, "TestContext\\simple.type");
			Assert.That(simpleType.Name, Is.EqualTo("simple"));
			Assert.That(simpleType.Members, Has.Count.EqualTo(2));
			Assert.That(simpleType.Members[0].Name, Is.EqualTo("number"));
			Assert.That(simpleType.Members[1].Name, Is.EqualTo("text"));
		}
		
		[Test]
		public void LoadTestTypeFile()
		{
			var simpleType = new TypeFileParser(TestContext, "TestContext\\test.type");
			Assert.That(simpleType.Name, Is.EqualTo("test"));
			Assert.That(simpleType.Members, Has.Count.EqualTo(4));
			Assert.That(simpleType.Members[0].Name, Is.EqualTo("first"));
			Assert.That(simpleType.Members[1].Name, Is.EqualTo("second"));
			Assert.That(simpleType.Members[2].Name, Is.EqualTo("text"));
			Assert.That(simpleType.Members[3].Name, Is.EqualTo("pi"));
		}

		[Test]
		public void TypeFileParserNeedsValidName()
		{
			Assert.Throws<Type.TypesMustStartWithLowercaseLetter>(
				() => new TypeFileParser(TestContext, "", new[] { "number a" }));
		}

		[Test]
		public void TypeFileParserNeedsValidContext()
		{
			Assert.Throws<Type.CannotCreateTypeInBaseContext>(
				() => new TypeFileParser(Base.Context, "abc", new[] { "number a" }));
		}

		[Test]
		public void TypeFileParserContentCannotBeEmpty()
		{
			Assert.Throws<TypeFileParser.ContentCannotBeEmpty>(
				() => new TypeFileParser(TestContext, "abc", new string[0]));
		}

		[Test]
		public void TypeCanBeCreatedWithContextAndMemberLines()
		{
			using (var type = new TypeFileParser(TestContext, "abc", new[] { "number value", "text" }))
			{
				Assert.That(type.Context.Name, Is.EqualTo("TestContext"));
				Assert.That(type.Name, Is.EqualTo("abc"));
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
				() => new TypeFileParser(TestContext, "abc", new[] { "number", "" }));
		}

		[Test]
		public void TypeMembersCanBeSetToInitialValue()
		{
			using (var type = new TypeFileParser(TestContext, "point", new[] { "x = 1", "y = 2" }))
			{
				Assert.That(type.Name, Is.EqualTo("point"));
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
				() => new TypeFileParser(TestContext, "assignFromOtherMember", new[] { "x = y" }));
			Assert.Throws<TypeFileParser.InvalidTypeMemberSyntax>(
				() => new TypeFileParser(TestContext, "assignTwice", new[] { "x = 1 = 2" }));
			Assert.Throws<TypeFileParser.InvalidTypeMemberSyntax>(
				() => new TypeFileParser(TestContext, "invalidSyntax", new[] { "x 1 2" }));
			Assert.Throws<Context.TypeNotFound>(
				() => new TypeFileParser(TestContext, "justNumber", new[] { "2385" }));
		}
	}
}