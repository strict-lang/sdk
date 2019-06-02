using System;
using System.Globalization;

namespace Strict.Statements
{
	/// <summary>
	/// Value that can be put into a member or passed into methods, basically just a type + value.
	/// Also base class for Number, Bool, Text and other statements just holding a value.
	/// </summary>
	public class Value : Statement
	{
		public Value(Type valueType, object currentValue = null) : base(valueType)
		{
			if (currentValue != null)
				CurrentValue = currentValue;
			else
				CurrentValue = ReturnType == Base.Text
					? "" : ReturnType == Base.Number ? 0.0 : ReturnType == Base.Bool ? (object)false : null;
		}

		public object CurrentValue { get; }

		public Value(string valueToParse) : base(ParseType(valueToParse))
		{
			CurrentValue = ParseValue(valueToParse);
		}

		private static Type ParseType(string valueToParse)
		{
			if (valueToParse[0] == '\"')
				return Base.Text;
			if (valueToParse == "true" || valueToParse == "false")
				return Base.Bool;
			double result;
			if (double.TryParse(valueToParse, out result))
				return Base.Number;
			throw new TypeParsingIsNotYetSupported(valueToParse);
		}

		public class TypeParsingIsNotYetSupported : Exception
		{
			public TypeParsingIsNotYetSupported(string valueToParse) : base(valueToParse) {}
		}
		
		private object ParseValue(string valueToParse)
		{
			return ReturnType == Base.Text
				? valueToParse.Substring(1, valueToParse.Length - 2)
				: ReturnType == Base.Number
					? double.Parse(valueToParse, CultureInfo.InvariantCulture)
					: ReturnType == Base.Bool ? (object)bool.Parse(valueToParse) : null;
		}

		public override bool Equals(Statement other)
		{
			var otherValue = other as Value;
			return otherValue != null && ReturnType == otherValue.ReturnType &&
				Equals(CurrentValue, otherValue.CurrentValue);
		}
	}
}