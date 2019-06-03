using NUnit.Framework;
using Strict.Statements;

namespace Strict.CollapsedGrammar.Tests
{
	public class MethodParserTests : TestWithParser
	{
		//TODO: we only need to parse statements, the method is already handled by MethodFileParserTests!
		[Test]
		public void ParseMethod()
		{
			var nodes = parser.Parse("void emptyMethod()");
			var method = nodes[0] as Method;
			Assert.That(method.ReturnType, Is.EqualTo(Base.Void));
			Assert.That(method.Name, Is.EqualTo("emptyMethod"));
			Assert.That(method.ToString(), Is.EqualTo("void Test.testMethod.emptyMethod()"));
		}

		[Test]
		public void ParseMethodWithParameters()
		{
			var nodes = parser.Parse("void simpleMethod(number, text)");
			var method = nodes[0] as Method;
			Assert.That(method.ReturnType, Is.EqualTo(Base.Void));
			Assert.That(method.Name, Is.EqualTo("simpleMethod"));
			Assert.That(method.Parameters.Count, Is.EqualTo(2));
			Assert.That(method.Parameters[0], Is.EqualTo(new Parameter(Base.Number, "number")));
			Assert.That(method.Parameters[1], Is.EqualTo(new Parameter(Base.Text, "text")));
			Assert.That(method.ToString(), Is.EqualTo("void Test.testMethod.simpleMethod(number, text)"));
		}

		[Test]
		public void ParseMethodWithNamedParameter()
		{
			var nodes = parser.Parse("void simpleMethod(number value)");
			var method = nodes[0] as Method;
			Assert.That(method.Parameters.Count, Is.EqualTo(1));
			Assert.That(method.Parameters[0], Is.EqualTo(new Parameter(Base.Number, "value")));
			Assert.That(method.ToString(), Is.EqualTo("void Test.testMethod.simpleMethod(value)"));
		}
	}
}