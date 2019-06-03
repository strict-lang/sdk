using System.Globalization;

namespace Strict.Statements
{
	/// <summary>
	/// Represents a constant number in code with a specific value to be assigned to a Member or
	/// passed along in a MethodCall. About numbers: https://strict.fogbugz.com/default.asp?W11
	/// </summary>
	public class Number : Value
	{
		public Number(double value) : base(Base.Number, value) {}

		public override string ToString()
			=> ((double)CurrentValue).ToString(CultureInfo.InvariantCulture);
	}
}