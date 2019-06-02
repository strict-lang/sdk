using NUnit.Framework;
using Strict.Statements;
using Text = Strict.Statements.Text;

namespace Strict.CollapsedGrammar.Tests
{
	public class BinaryOperationParserTests : TestWithParser
	{
		[TestCase("1+2", 1, 2)]
		[TestCase("2-5", 2, 5)]
		[TestCase("2/4", 2, 4)]
		[TestCase("5*8", 5, 8)]
		[TestCase("7%2", 7, 2)]
		[TestCase("1<3", 1, 3)]
		[TestCase("4>2", 4, 2)]
		[TestCase("1 is 1", 1, 1)]
		[TestCase("3 isnot 5", 3, 5)]
		[TestCase("1 and 2", 1, 2)]
		[TestCase("4 or 5", 4, 5)]
		public void CheckOperator(string statement, int left, int right)
		{
			string operatorName =
				statement.Substring(1, statement.Length - 2).Replace("+", "add").Replace("-", "subtract").
					Replace("*", "multiply").Replace("/", "divide").Replace("%", "modulate").Replace("<",
						"smaller").Replace(">", "bigger").Trim();
			CheckOperator(parser.Parse(statement)[0], operatorName, left, right);
		}

		private static void CheckOperator(Statement statement, string operatorName, int left, int right)
		{
			var operation = statement as BinaryOperation;
			Assert.That(operation.ReturnType, Is.EqualTo(Base.Number));
			Assert.That(operation.Method.Name, Is.EqualTo(operatorName));
			Assert.That(operation.Left, Is.EqualTo(new Number(left)));
			Assert.That(operation.Right, Is.EqualTo(new Number(right)));
		}
		
		[Test]
		public void AddAndMultiply()
		{
			var nodes = parser.Parse("2+3*4");
			Assert.That(nodes, Has.Count.EqualTo(1));
			var add = nodes[0] as BinaryOperation;
			Assert.That(add.Method.Name, Is.EqualTo("add"));
			Assert.That(add.ReturnType, Is.EqualTo(Base.Number));
			Assert.That(add.Left, Is.EqualTo(new Number(2)));
			CheckOperator(add.Right, "multiply", 3, 4);
		}

		[Test]
		public void AddStringTexts()
		{
			var nodes = parser.Parse("\"Hi\" + \" there\"");
			var add = nodes[0] as BinaryOperation;
			Assert.That(add.ReturnType, Is.EqualTo(Base.Text));
			Assert.That((add.Left as Text).CurrentValue, Is.EqualTo("Hi"));
			Assert.That((add.Right as Text).CurrentValue, Is.EqualTo(" there"));
		}

		[Test]
		public void SubtractAndDivide()
		{
			var nodes = parser.Parse("6/4-2");
			Assert.That(nodes, Has.Count.EqualTo(1));
			var subtract = nodes[0] as BinaryOperation;
			Assert.That(subtract.Method.Name, Is.EqualTo("subtract"));
			Assert.That(subtract.ReturnType, Is.EqualTo(Base.Number));
			CheckDivideNumbers(subtract.Left as BinaryOperation);
			Assert.That(subtract.Right, Is.EqualTo(new Number(2)));
		}

		private static void CheckDivideNumbers(BinaryOperation divide)
		{
			Assert.That(divide.ReturnType, Is.EqualTo(Base.Number));
			Assert.That(divide.Left, Is.EqualTo(new Number(6)));
			Assert.That(divide.Right, Is.EqualTo(new Number(4)));
		}

		[Test]
		public void ComplexBinaryExpression()
		{
			var nodes = parser.Parse("number = 2/1 + 4*3");
			Assert.That(nodes, Has.Count.EqualTo(1));
			Assert.That(nodes[0], Is.InstanceOf<Member>());
			var member = nodes[0] as Member;
			Assert.That(member.Name, Is.EqualTo("number"));
			var value = member.Value as BinaryOperation;
			Assert.That(value.ReturnType, Is.EqualTo(Base.Number));
			CheckOperator(value.Left, "divide", 2, 1);
			Assert.That(value.Operator, Is.EqualTo("+"));
			CheckOperator(value.Right, "multiply", 4, 3);
		}
	}
}