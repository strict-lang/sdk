namespace Strict.Statements
{
	/// <summary>
	/// Represents a boolean value, can only be true or false.
	/// </summary>
	public class Bool : Value
	{
		public Bool(bool value) : base(Base.Bool, value) {}
		public override string ToString() => CurrentValue.ToString().ToLower();
	}
}