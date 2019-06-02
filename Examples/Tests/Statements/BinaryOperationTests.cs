using NUnit.Framework;
using Strict.Statements;
using Text = Strict.Statements.Text;

namespace Strict.Tests.Statements
{
	public class BinaryOperationTests
	{
		[Test]
		public void BinaryExpressionWithDifferentArgumentTypesIsNotAllowed()
		{
			Assert.Throws<BinaryOperation.TypesMustMatchForOperator>(() =>
				new BinaryOperation(new Number(1), BinaryOperator.Add, new Text("Yo")));
		}

		[Test]
		public void CreateBinaryExpression()
		{
			var binary = new BinaryOperation(new Number(3), BinaryOperator.Add, new Number(5));
			Assert.That(binary.ReturnType, Is.EqualTo(Base.Number));
			Assert.That(binary.Operator, Is.EqualTo("+"));
			Assert.That((binary.Left as Number).CurrentValue, Is.EqualTo(3));
			Assert.That((binary.Right as Number).CurrentValue, Is.EqualTo(5));
			Assert.That(binary.ToString(), Is.EqualTo("3 + 5"));
		}

		[Test]
		public void GetInvalidOperatorMethod()
		{
			Assert.Throws<Type.OperatorIsNotSupportedForThisType>(
				() => Base.Number.GetBinaryMethod((BinaryOperator)3945));
		}

		[TestCase(BinaryOperator.Add)]
		[TestCase(BinaryOperator.Subtract)]
		[TestCase(BinaryOperator.Multiply)]
		[TestCase(BinaryOperator.Divide)]
		[TestCase(BinaryOperator.Modulate)]
		[TestCase(BinaryOperator.And)]
		[TestCase(BinaryOperator.Or)]
		[TestCase(BinaryOperator.Is)]
		[TestCase(BinaryOperator.IsNot)]
		[TestCase(BinaryOperator.Smaller)]
		[TestCase(BinaryOperator.Bigger)]
		public void GetBinaryOperatorMethodName(BinaryOperator binaryOperator)
		{
			Assert.That(Base.Number.GetBinaryMethod(binaryOperator).Name,
				Is.EqualTo(binaryOperator.ToString().ToLower()));
		}

		[TestCase(BinaryOperator.Add, "+")]
		[TestCase(BinaryOperator.Subtract, "-")]
		[TestCase(BinaryOperator.Multiply, "*")]
		[TestCase(BinaryOperator.Divide, "/")]
		[TestCase(BinaryOperator.Modulate, "%")]
		[TestCase(BinaryOperator.And, "and")]
		[TestCase(BinaryOperator.Or, "or")]
		[TestCase(BinaryOperator.Is, "is")]
		[TestCase(BinaryOperator.IsNot, "isnot")]
		[TestCase(BinaryOperator.Smaller, "<")]
		[TestCase(BinaryOperator.Bigger, ">")]
		public void GetBinaryOperatorSymbol(BinaryOperator binaryOperator, string operatorSymbol)
		{
			var binary = new BinaryOperation(new Number(3), binaryOperator, new Number(5));
			Assert.That(binary.Operator, Is.EqualTo(operatorSymbol));
		}
	}
}