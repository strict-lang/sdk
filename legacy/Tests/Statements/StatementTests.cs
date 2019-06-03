using System.Linq;
using NUnit.Framework;
using Strict.Statements;

namespace Strict.Tests.Statements
{
	public class StatementTests : TestWithContext
	{
		[Test]
		public void CreateNestedStatements()
		{
			var firstEntry = new Number(1);
			var secondEntry = new Number(2);
			var methodCall = new BinaryOperation(firstEntry, BinaryOperator.Add, secondEntry);
			Assert.That(methodCall.Children, Has.Count.EqualTo(2));
			Assert.That(firstEntry.Children, Is.Empty);
			Assert.That(secondEntry.Children, Is.Empty);
			Assert.That(methodCall.Parent, Is.Null);
			Assert.That(firstEntry.Parent, Is.EqualTo(methodCall));
			Assert.That(secondEntry.Parent, Is.EqualTo(methodCall));
		}

		[Test]
		public void AddChildrenAfterToMethodTree()
		{
			using (var method = new Method(TestContext, Base.Number, "test"))
			{
				var thirdEntry = new Number(1);
				method.Add(thirdEntry);
				Assert.That(thirdEntry.Parent, Is.EqualTo(method));
			}
		}

		[Test]
		public void RemoveChildrenAfterTreeCreation()
		{
			using (var method = new Method(TestContext, Base.Number, "test"))
			{
				method.Add(new Number(1));
				method.Add(new Number(2));
				method.Remove(method.Children.Last());
				Assert.That(method.Children, Has.Count.EqualTo(1));
			}
		}

		[Test]
		public void CreateBinaryOperationMethodCall()
		{
			var left = new Number(1);
			var right = new Number(2);
			var binaryExpression = new BinaryOperation(left, BinaryOperator.Add, right);
			using (var method = new Method(TestContext, Base.Number, "test"))
			{
				method.Add(binaryExpression);
				Assert.That(binaryExpression.Parent, Is.EqualTo(method));
				Assert.That(binaryExpression.Left, Is.EqualTo(left));
				Assert.That(binaryExpression.Right, Is.EqualTo(right));
				Assert.That(binaryExpression.Left.Parent, Is.EqualTo(binaryExpression));
				Assert.That(binaryExpression.Right.Parent, Is.EqualTo(binaryExpression));
				Assert.That(binaryExpression.Children, Has.Count.EqualTo(2));
			}
		}

		[Test]
		public void AddingSameEntryTwiceIsNotAllowed()
		{
			using (var method = new Method(TestContext, Base.Number, "test"))
			{
				var number = new Number(1);
				method.Add(number);
				Assert.Throws<Statement.CannotAddChildThatAlreadyHasAParent>(() => method.Add(number));
			}
		}

		[Test]
		public void RemovingNotPreviouslyAddedEntryIsNotAllowed()
		{
			using (var method = new Method(TestContext, Base.Number, "test"))
			{
				var firstEntry = new Number(1);
				Assert.Throws<Statement.CannotRemoveChildThatIsNotLinkedToThisParent>(
					() => method.Remove(firstEntry));
			}
		}
	}
}