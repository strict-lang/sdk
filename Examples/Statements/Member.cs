using System;
using System.Collections.Generic;
using System.Linq;

namespace Strict.Statements
{
	/// <summary>
	/// A member can be defined in a type or method. Global or namespace members are not allowed.
	/// Each member is assigned with a Statement value, which does not have to be constant, it could
	/// be a method or any other statement. Assignment is not allowed in Strict, you have to define
	/// a new member if anything changes and this new member is only available on other threads when
	/// re-evaluated. Try to be as functional as possible to avoid side-effects.
	/// </summary>
	public class Member : NamedStatement
	{
		public Member(string name, Statement value) : base(value.ReturnType, name, value) {}
		public IReadOnlyList<Statement> UsedBy => usedBy;
		private readonly List<Statement> usedBy = new List<Statement>();
		public Statement Value
		{
			get
			{
				if (!Children.Any())
					throw new MemberWasOverwrittenAndIsClearedAndIsNotLongerAccessible(this);
				return Children[0];
			}
		}

		public override string ToString() => Name + " = " + Value;

		public override bool Equals(Statement other)
		{
			if (Children.Count == 0)
				throw new MemberWasOverwrittenAndIsClearedAndIsNotLongerAccessible(this);
			var otherMember = other as Member;
			if (otherMember != null && otherMember.Children.Count == 0)
				throw new MemberWasOverwrittenAndIsClearedAndIsNotLongerAccessible(otherMember);
			return base.Equals(other) && otherMember != null && Value.Equals(otherMember.Value);
		}

		public class MemberWasOverwrittenAndIsClearedAndIsNotLongerAccessible : Exception
		{
			public MemberWasOverwrittenAndIsClearedAndIsNotLongerAccessible(Member member)
				: base(member.Name + " (" + member.ReturnType + ")") {}
		}

		/// <summary>
		/// Normally we do not need to worry about members, they always exists for Types and in Methods
		/// they are defined and used until we go out of scope and everything is cleaned up. However
		/// when a member is reassigned to a new value or statement, the old one with the same name
		/// becomes inactive and unusable. Kill all UsedBy links here and make it unusable from now on.
		/// </summary>
		public void ClearAndMakeUnusable()
		{
			usedBy.Clear();
			ReturnType?.usedBy.Remove(this);
			while (Children.Any())
				Remove(Children[0]);
		}

		public class MembersMustBeDefinedInTypesOrMethods : Exception {}
	}
}