using System;

namespace Strict.Statements
{
	public class BinaryOperation : MethodCall
	{
		public BinaryOperation(Statement left, BinaryOperator binaryOperator, Statement right)
			: base(left.ReturnType.GetBinaryMethod(binaryOperator), left, right)
		{
			if (left.ReturnType != right.ReturnType)
				throw new TypesMustMatchForOperator(left.ReturnType, binaryOperator, right.ReturnType);
		}

		public class TypesMustMatchForOperator : Exception
		{
			public TypesMustMatchForOperator(Type leftType, BinaryOperator binary, Type rightType)
				: base(leftType + " " + binary + " " + rightType) {}
		}

		public override string ToString() => Left + " " + Operator + " " + Right;
		public Statement Left => Arguments[0];

		public string Operator
		{
			get
			{
				var operatorName = Method.Name;
				switch (operatorName)
				{
				case "add":
					return "+";
				case "subtract":
					return "-";
				case "multiply":
					return "*";
				case "divide":
					return "/";
				case "modulate":
					return "%";
				case "smaller":
					return "<";
				case "bigger":
					return ">";
				default:
					return operatorName;
				}
			}
		}
		public Statement Right => Arguments[1];
	}
}