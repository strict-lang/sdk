using NUnit.Framework;
using Strict.Statements;
using Text = Strict.Statements.Text;

namespace Strict.Tests.Statements
{
	public class NegationTests
	{
		[Test]
		public void NegateIsNotAllowedForNonNumberOrBool()
		{
			Assert.Throws<Type.OperatorIsNotSupportedForThisType>(
				() => new Negation(new Text("Hi")));
		}

		[Test]
		public void NegateNumber()
		{
			var unary = new Negation(new Number(3));
			Assert.That(unary.ReturnType, Is.EqualTo(Base.Number));
			Assert.That(unary.Operator, Is.EqualTo("-"));
			Assert.That((unary.Argument as Number).CurrentValue, Is.EqualTo(3));
			Assert.That(unary.ToString(), Is.EqualTo("-3"));
		}

		[Test]
		public void NegateBool()
		{
			var unary = new Negation(new Bool(true));
			Assert.That(unary.ReturnType, Is.EqualTo(Base.Bool));
			Assert.That(unary.Operator, Is.EqualTo("not"));
			Assert.That((unary.Argument as Bool).CurrentValue, Is.True);
			Assert.That(unary.ToString(), Is.EqualTo("not true"));
		}
	}
}