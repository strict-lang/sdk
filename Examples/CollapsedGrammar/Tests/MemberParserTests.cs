using NUnit.Framework;
using Strict.Statements;

namespace Strict.CollapsedGrammar.Tests
{
	public class MemberParserTests : TestWithParser
	{
		[Test]
		public void ParseInvalidAssignment()
		{
			Assert.Throws<Parser.ParsingFailed>(() => parser.Parse("a = a = 5"));
		}

		[Test]
		public void ParseInvalidMemberDeclaration()
		{
			Assert.Throws<Parser.ParsingFailed>(() => parser.Parse("number string abc = \"jo\""));
		}

		[Test]
		public void ParseMemberOutsideMethodIsNotAllowed()
		{
			var simpleParser = new Parser(TestContext);
			Assert.Throws<Member.MembersMustBeDefinedInTypesOrMethods>(
				() => simpleParser.Parse("number = 2"));
		}

		[Test]
		public void ParseIdentifierFromSingleWord()
		{
			var nodes = parser.Parse("abc = 0");
			Assert.That(nodes, Has.Count.EqualTo(1));
			CheckMember(nodes[0], "abc", 0);
		}

		protected static void CheckMember(Statement statement, string name, int value)
		{
			Assert.That(statement, Is.InstanceOf<Member>());
			var member = statement as Member;
			Assert.That(member.Name, Is.EqualTo(name));
			Assert.That(member.ToString(), Is.EqualTo(name + " = " + value));
			Assert.That(member.Value, Is.EqualTo(new Number(value)));
			Assert.That(member.UsedBy, Has.Count.EqualTo(0));
		}

		[Test]
		public void ParseIdentifierFromSingleWordWithWhitespaces()
		{
			var nodes = parser.Parse("   abc   =  3");
			Assert.That(nodes, Has.Count.EqualTo(1));
			CheckMember(nodes[0], "abc", 3);
		}

		[Test]
		public void AssignmentToString()
		{
			var nodes = parser.Parse("number = 5");
			Assert.That(nodes[0].ToString(), Is.EqualTo("number = 5"));
		}

		[Test]
		public void ParseTwoIdentifiersFromDifferentWords()
		{
			var nodes = parser.Parse("abc = 0\n" + "def = 1");
			Assert.That(nodes, Has.Count.EqualTo(2));
			CheckMember(nodes[0], "abc", 0);
			CheckMember(nodes[1], "def", 1);
		}
	}
}