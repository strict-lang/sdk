using NUnit.Framework;

namespace Strict.Parsing.Tests
{
	public class TypeFileParserTests
	{
		[Test]
		public void TypeFileParserNeedsFilenameEndingWithType()
		{
			Assert.Throws<TypeFileParser.FilenameNeedsToBeValidAndEndWithType>(
				() => new TypeFileParser(""));
			Assert.Throws<TypeFileParser.FilenameNeedsToBeValidAndEndWithType>(
				() => new TypeFileParser("abc"));
			Assert.Throws<TypeFileParser.FilenameNeedsToBeValidAndEndWithType>(
				() => new TypeFileParser("C:\\"));
		}

		[Test]
		public void FilePathMustBeInADirectoryForTheNamespace()
		{
			Assert.Throws<TypeFileParser.TypeFileMustBeInFolderForItsNamespace>(
				() => new TypeFileParser("Empty.type"));
		}

		[Test]
		public void DifferentStartDirectoryIsNotAllowedForTheNamespace()
		{
			Assert.Throws<TypeFileParser.FilePathMustBeBelowStrictExecutableFolder>(
				() => new TypeFileParser("c:\\Empty.type"));
		}

		[Test]
		public void EmptyTypeFileIsNotAllowed()
		{
			Assert.Throws<TypeFileParser.ContentCannotBeEmpty>(
				() => new TypeFileParser("TestNamespace\\Empty.type"));
		}

		[Test]
		public void TypeFileParserNeedsValidNamespace()
		{
			Assert.Throws<TypeFileParser.NamespaceNameMustNotBeEmpty>(
				() => new TypeFileParser("", "", new[] { "number a" }));
		}

		[Test]
		public void TypeFileParserNeedsValidName()
		{
			Assert.Throws<TypeFileParser.TypeNameCannotBeEmpty>(
				() => new TypeFileParser("Test", "", new[] { "number a" }));
		}

		[Test]
		public void TypeFileParserNameMustStartWithUpperCaseLetter()
		{
			Assert.Throws<TypeFileParser.NameMustStartWithUpperCaseLetter>(
				() => new TypeFileParser("Test", "abc", new[] { "number a" }));
		}

		[Test]
		public void TypeFileParserContentCannotBeEmpty()
		{
			Assert.Throws<TypeFileParser.ContentCannotBeEmpty>(
				() => new TypeFileParser("Test", "Abc", new string[0]));
		}

		[Test]
		public void TypeCanBeCreatedWithNamespaceNameAndContent()
		{
			var parser = new TypeFileParser("Test", "Abc", new[] { "number value" });
			Assert.AreEqual("Test", parser.Namespace);
			Assert.AreEqual("Abc", parser.Name);
			Assert.AreEqual(1, parser.Members.Count);
			//TODO, after namespace scope merge: Assert.AreEqual("value", parser.Members[0].Name);
		}
	}
}