using NUnit.Framework;
using Strict.Statements;

namespace Strict.CollapsedGrammar.Tests
{
	/// <summary>
	/// Parse examples are in the Strict Wiki: https://strict.fogbugz.com/default.asp?W1
	/// Tests are split into many files found in this namespace, they all test the Parser class 
	/// </summary>
	public class ParserTests : TestWithParser
	{
		[Test]
		public void CreatingParserInBaseContextIsNotAllowed()
		{
			Assert.Throws<Parser.CannotParseFilesInBaseContext>(() => new Parser(Base.Context));
		}

		[Test]
		public void ParsingNothingResultsInNoStatements()
		{
			var statements = parser.Parse("");
			Assert.That(statements, Is.Empty);
		}

		[Test]
		public void ParseSyntaxError()
		{
			Assert.Throws<Parser.ParsingFailed>(() => parser.Parse("$%@%(#"));
		}
	}
}