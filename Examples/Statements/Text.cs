namespace Strict.Statements
{
	/// <summary>
	/// Represents a constant string text, often combined with other things to build texts.
	/// </summary>
	public class Text : Value
	{
		public Text(string value) : base(Base.Text, value) {}
		public override string ToString() => "\"" + CurrentValue + "\"";
	}
}