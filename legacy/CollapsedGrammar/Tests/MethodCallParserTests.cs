using NUnit.Framework;
using Strict.Statements;

namespace Strict.CollapsedGrammar.Tests
{
	public class MethodCallParserTests : TestWithParser
	{
		[Test]
		public void ParseMethodCall()
		{
			var nodes = parser.Parse("testMethod()");
			Assert.That(nodes, Has.Count.EqualTo(1));
			Assert.That(nodes[0], Is.InstanceOf<MethodCall>());
			var methodCall = nodes[0] as MethodCall;
			Assert.That(methodCall.Method.Name, Is.EqualTo("testMethod"));
			Assert.That(methodCall.ToString(), Is.EqualTo("testMethod()"));
		}

		/*TODO
		[Test]
		public void ParseMethodCallWithContext()
		{
			var nodes = parser.Parse("System.Console.WriteLine()");
			Assert.That(nodes, Has.Count.EqualTo(1));
			var methodCall = nodes[0] as MethodCall;
			Assert.That(methodCall.Context.ToString(), Is.EqualTo("System.Console"));
			Assert.That(methodCall.Identifier.Name, Is.EqualTo("WriteLine"));
			Assert.That(methodCall.ToString(), Is.EqualTo("System.Console.WriteLine()"));
		}

		[Test]
		public void ParseMethodWithArguments()
		{
			var nodes = parser.Parse("Add(1, 2)");
			Assert.That(nodes, Has.Count.EqualTo(1));
			var methodCall = nodes[0] as MethodCall;
			Assert.That(methodCall.Identifier.Name, Is.EqualTo("Add"));
			Assert.That(methodCall.Arguments, Has.Count.EqualTo(2));
			Assert.That((methodCall.Arguments[0] as Number).Value, Is.EqualTo(1));
			Assert.That((methodCall.Arguments[1] as Number).Value, Is.EqualTo(2));
		}

		[Test]
		public void ParseMethodWithContextAndArguments()
		{
			var nodes = parser.Parse("System.Console.WriteLine(\"Hi There\")");
			Assert.That(nodes, Has.Count.EqualTo(1));
			var methodCall = nodes[0] as MethodCall;
			Assert.That(methodCall.Identifier.Name, Is.EqualTo("WriteLine"));
			Assert.That(methodCall.ToString(), Is.EqualTo("System.Console.WriteLine(\"Hi There\")"));
			Assert.That(methodCall.Arguments, Has.Count.EqualTo(1));
			Assert.That((methodCall.Arguments[0] as StringText).Text, Is.EqualTo("Hi There"));
		}

		[Test]
		public void ParseMethodCall()
		{
			var parseTree = Parse("WriteHiThere()");
			Assert.That(parseTree.Root.ChildNodes, Has.Count.EqualTo(1));
			ParseTreeNode expressionNode = parseTree.Root.ChildNodes[0];
			Assert.That(expressionNode.ChildNodes, Has.Count.EqualTo(2));
			Assert.That(expressionNode.ChildNodes[0].Token.Value, Is.EqualTo("WriteHiThere"));
			var arguments = expressionNode.ChildNodes[1];
			Assert.That(arguments.Term.Name, Is.EqualTo("ArgumentList"));
			Assert.That(arguments.ChildNodes, Has.Count.EqualTo(0));
		}
		 */
	}
}