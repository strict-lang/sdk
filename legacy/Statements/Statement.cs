using System;
using System.Collections.Generic;

namespace Strict.Statements
{
	/// <summary>
	/// Each statement is just one code line, they always have a parent except for root statements in
	/// files (.method) or if just evaluating or parsing a simple code block. All statements have a
	/// return type and child statements, which are hocked up automatically via Add to this parent.
	/// </summary>
	public abstract class Statement : IEquatable<Statement>
	{
		/// <summary>
		/// Statements can constantly be created and reused by moving them to new parents. Comparing
		/// statements is also easily supported for testing and pattern matching (parents are not
		/// checked, but all data and usually all children must match). Statements are often GCed.
		/// </summary>
		protected Statement(Type returnType, params Statement[] childStatement)
		{
			if (!(this is Type) && (returnType == null || returnType.IsDisposed))
				throw new EveryStatementNeedsAValidReturnType(returnType);
			ReturnType = returnType;
			returnType?.usedBy.Add(this);
			foreach (var child in childStatement)
				Add(child.Parent != null ? child.Clone() : child);
		}

		public class EveryStatementNeedsAValidReturnType : Exception
		{
			public EveryStatementNeedsAValidReturnType(Type returnType)
				: base(returnType + " is disposed: " + returnType?.IsDisposed) {}
		}

		public Type ReturnType { get; }
		public Statement Parent { get; private set; }
		private readonly List<Statement> children = new List<Statement>();
		public IReadOnlyList<Statement> Children => children;
		public abstract bool Equals(Statement other);

		private Statement Clone()
		{
			var clonedStatement = MemberwiseClone() as Statement;
			clonedStatement.Parent = null;
			return clonedStatement;
		}

		public void Add(Statement child)
		{
			if (child.Parent != null)
				throw new CannotAddChildThatAlreadyHasAParent(child);
			children.Add(child);
			child.Parent = this;
		}

		public void Remove(Statement child)
		{
			if (!Equals(child.Parent, this))
				throw new CannotRemoveChildThatIsNotLinkedToThisParent(child);
			children.Remove(child);
			child.Parent = null;
		}

		public override bool Equals(object other) => Equals(other as Statement);

		public override int GetHashCode()
		{
			unchecked
			{
				return (children?.GetHashCode() * 397 ?? 0) ^ (ReturnType?.GetHashCode() ?? 0);
			}
		}

		public class CannotAddChildThatAlreadyHasAParent : Exception
		{
			public CannotAddChildThatAlreadyHasAParent(Statement child)
				: base("Child: " + child + ", Parent: " + child.Parent) {}
		}

		public class CannotRemoveChildThatIsNotLinkedToThisParent : Exception
		{
			public CannotRemoveChildThatIsNotLinkedToThisParent(Statement child)
				: base("Child: " + child + ", Parent: " + child.Parent) {}
		}
	}
}