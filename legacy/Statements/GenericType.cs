namespace Strict.Statements
{
	/// <summary>
	/// Currently we have no good way to support generic C# types, so this is our placeholder.
	/// </summary>
	public class GenericType : Type
	{
		public GenericType(Context inContext, string name, string constraintTypeName)
			: base(inContext, name + "<>")
		{
			this.constraintTypeName = constraintTypeName;
		}

		private readonly string constraintTypeName;
		public bool IsStructConstraint => constraintTypeName == "struct";
		public Type ConstraintType
			=> IsStructConstraint ? null : Base.Context.GetTypeInAnyChildContext(constraintTypeName);
	}
}