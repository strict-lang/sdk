namespace Strict.CollapsedGrammar.Tests
{
	public class NumberParserTests : TestWithParser
	{
		/*TODO
		[Test]
		public void ParseIntNumber()
		{
			var nodes = parser.Parse("number value = 2");
			var declaration = nodes[0] as MemberDeclaration;
			CheckSingleUseNumber(declaration.Value as Number, 2, new Position(1, 16));
			Assert.That(declaration.Type.Name, Is.EqualTo("number"));
			Assert.That(declaration.Identifier.Name, Is.EqualTo("value"));
		}

		[Test]
		public void ParseFloatingPointNumber()
		{
			var nodes = parser.Parse("number value = 1.5");
			var declaration = nodes[0] as MemberDeclaration;
			CheckSingleUseNumber(declaration.Value as Number, 1.5, new Position(1, 16));
		}

		[Test]
		public void ParseCSharpFloatNumberIsNotAllowed()
		{
			Assert.Throws<Number.FormatNotSupported>(() => parser.Parse("number value = 2.5f"));
		}

		[Test]
		public void ParseUintNumberIsNotSupported()
		{
			Assert.Throws<Number.FormatNotSupported>(() => parser.Parse("number value = 3U"));
		}
		*/
	}
}