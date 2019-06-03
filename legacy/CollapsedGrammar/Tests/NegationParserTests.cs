using NUnit.Framework;
using Strict.Statements;

namespace Strict.CollapsedGrammar.Tests
{
	public class NegationParserTests : TestWithParser
	{
		[Test]
		public void ParseNumberNegation()
		{
			var nodes = parser.Parse("-5");
			var unary = nodes[0] as Negation;
			Assert.That(unary.ReturnType, Is.EqualTo(Base.Number));
			Assert.That((unary.Argument as Number).CurrentValue, Is.EqualTo(5));
		}

		[Test]
		public void ParseBoolNegation()
		{
			var nodes = parser.Parse("not false");
			var negation = nodes[0] as Negation;
			Assert.That(negation.ReturnType, Is.EqualTo(Base.Bool));
			Assert.That(negation.Operator, Is.EqualTo("not"));
			Assert.That((negation.Argument as Bool).CurrentValue, Is.False);
		}
	}
}