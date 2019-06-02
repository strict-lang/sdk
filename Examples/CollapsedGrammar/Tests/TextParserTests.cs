namespace Strict.CollapsedGrammar.Tests
{
	public class TextParserTests : TestWithParser
	{
		/*TODO
		[Test]
		public void ParseEmptyString()
		{
			ParseStringAndCheckIt("");
		}

		private void ParseStringAndCheckIt(params string[] texts)
		{
			var input = "";
			foreach (var text in texts)
				input += "string text = \"" + text + "\"\n";
			var nodes = parser.Parse(input);
			Assert.That(nodes.Count, Is.EqualTo(texts.Length));
			var declaration = nodes[0] as MemberDeclaration;
			Assert.IsInstanceOf<StringText>(declaration.Value);
			for (int index = 0; index < texts.Length; index++)
				Assert.That((nodes[index] as MemberDeclaration).Value.ToString(),
					Is.EqualTo("\"" + texts[index] + "\""));
		}

		[Test]
		public void ParseHelloString()
		{
			ParseStringAndCheckIt("Hello");
		}

		[Test]
		public void ParseMultipleStrings()
		{
			ParseStringAndCheckIt("abc", "def", "");
		}

		[Test]
		public void StringTextToString()
		{
			var nodes = parser.Parse("string text = \"abc\"");
			var declaration = nodes[0] as MemberDeclaration;
			Assert.That(declaration.Value.ToString(), Is.EqualTo("\"abc\""));
		}
	 */
	}
}