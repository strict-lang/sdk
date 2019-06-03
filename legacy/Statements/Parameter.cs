namespace Strict.Statements
{
	/// <summary>
	/// Parameter for method definitions. Arguments must match the types defined here.
	/// </summary>
	public class Parameter : NamedStatement
	{
		public Parameter(Type parameterType, string name) : base(parameterType, name) {}
		public override string ToString() => Name;

		public override bool Equals(Statement other)
		{
			var otherParameter = other as Parameter;
			return base.Equals(other) && otherParameter != null && ReturnType == otherParameter.ReturnType;
		}
	}
}