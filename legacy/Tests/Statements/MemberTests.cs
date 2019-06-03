using System.Linq;
using NUnit.Framework;
using Strict.Statements;
using Text = Strict.Statements.Text;

namespace Strict.Tests.Statements
{
	public class MemberTests: TestWithContext
	{
		[Test]
		public void CreateNumber()
		{
			var member = new Member("value", new Number(5));
			Assert.That(member.ReturnType, Is.EqualTo(Base.Number));
			Assert.That(member.Name, Is.EqualTo("value"));
			Assert.That((member.Value as Number).CurrentValue, Is.EqualTo(5));
			Assert.That(member.ToString(), Is.EqualTo("value = 5"));
		}

		[Test]
		public void CreateString()
		{
			var member = new Member("text", new Text("Hey"));
			Assert.That(member.ReturnType, Is.EqualTo(Base.Text));
			Assert.That(member.Name, Is.EqualTo("text"));
			Assert.That((member.Value as Text).CurrentValue, Is.EqualTo("Hey"));
			Assert.That(member.ToString(), Is.EqualTo("text = \"Hey\""));
		}

		[Test]
		public void ParameterToString()
		{
			Assert.That(new Parameter(Base.Number, "value").ToString(), Is.EqualTo("value"));
		}

		[Test]
		public void AddingTheSameMemberTwiceIsNotAllowed()
		{
			var testType = new Type(TestContext, "test");
			var value = new Member("value", new Number(0));
			testType.Add(value);
			Assert.Throws<Statement.CannotAddChildThatAlreadyHasAParent>(() => testType.Add(value));
		}

		[Test]
		public void EveryStatementNeedsAValidType()
		{
			Assert.Throws<Statement.EveryStatementNeedsAValidReturnType>(() => new Value(null, null));
		}
		
		[Test]
		public void ParseValueFromString()
		{
			Assert.That(new Value("\"\"").ReturnType, Is.EqualTo(Base.Text));
			Assert.That(new Value("true").ReturnType, Is.EqualTo(Base.Bool));
			Assert.That(new Value("0.1").ReturnType, Is.EqualTo(Base.Number));
		}
		
		[Test]
		public void MembersOfTypeAreAccessibleInMethod()
		{
			var testType = new Type(TestContext, "test");
			UsingTheTypeInAMemberMakesItAvailableInTheTypeScope(testType);
			AddingANewMemberInAMethodMakesItAvailable(testType);
		}

		private static void UsingTheTypeInAMemberMakesItAvailableInTheTypeScope(Type testType)
		{
			Assert.That(testType.Context.AllAccessibleMembersRecursively, Is.Empty);
			testType.Context.AddOrReplace(new Member("value", new Number(0)));
			Assert.That(testType.Context.AllAccessibleMembersRecursively, Is.Not.Empty);
		}
		

		private static void AddingANewMemberInAMethodMakesItAvailable(Type testType)
		{
			var method = new Method(testType.Context, testType, "testMethod");
			Assert.That(method.Scope.AllAccessibleMembersRecursively.Count(), Is.EqualTo(1));
			method.Add(new Member("value2", new Number(1)));
			Assert.That(method.Scope.AllAccessibleMembersRecursively.Count(), Is.EqualTo(2));
		}
		
		[Test]
		public void MembersOfTypeAreAccessibleInNestedContextToo()
		{
			var testType = new Type(TestContext, "testType");
			var member = new Member("value", new Value(testType));
			TestContext.AddOrReplace(member);
			using (var nestedContext = new Context(TestContext, "NestedContext"))
				Assert.That(nestedContext.AllAccessibleMembersRecursively, Contains.Item(member));
		}
		
		[Test]
		public void MemberEquals()
		{
			var member = new Member("value", new Number(5));
			Assert.That(member, Is.EqualTo(member));
			Assert.That(member, Is.EqualTo(new Member(member.Name, member.Value)));
			Assert.That(member, Is.EqualTo(new Member(member.Name, new Number(5))));
			Assert.That(member, Is.Not.EqualTo(member.Value));
			Assert.That(member.GetHashCode(),
				Is.EqualTo(member.Children.GetHashCode() * 397 ^ member.ReturnType.GetHashCode()));
		}

		[Test]
		public void ReplacingAMemberMakesTheOldOneUnusable()
		{
			var oldMember = new Member("value", new Number(1));
			TestMethod.Add(oldMember);
			Assert.That(oldMember.Value, Is.EqualTo(new Number(1)));
			var newMember = new Member("value", new Number(2));
			TestMethod.Add(newMember);
			Assert.That(newMember.Value, Is.EqualTo(new Number(2)));
			Assert.Throws<Member.MemberWasOverwrittenAndIsClearedAndIsNotLongerAccessible>(
				() => Assert.That(oldMember.Equals(newMember), Is.True));
			Assert.Throws<Member.MemberWasOverwrittenAndIsClearedAndIsNotLongerAccessible>(
				() => Assert.That(newMember.Equals(oldMember), Is.True));
			Assert.Throws<Member.MemberWasOverwrittenAndIsClearedAndIsNotLongerAccessible>(
				() => Assert.That(oldMember.Value, Is.EqualTo(1)));
		}
	}
}