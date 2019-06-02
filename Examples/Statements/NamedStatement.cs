using System;

namespace Strict.Statements
{
	public class NamedStatement : Statement
	{
		public NamedStatement(Type returnType, string name, params Statement[] childStatement)
			: base(returnType, childStatement)
		{
			if (string.IsNullOrEmpty(name) || char.IsUpper(name[0]))
				throw new MemberOrMethodMustStartWithLowerCaseLetter(name);
			Name = name;
		}

		public class MemberOrMethodMustStartWithLowerCaseLetter : Exception
		{
			public MemberOrMethodMustStartWithLowerCaseLetter(string name) : base(name) {}
		}

		public string Name { get; }

		public override bool Equals(Statement other)
			=> other is NamedStatement && Name == ((NamedStatement)other).Name;
	}
}