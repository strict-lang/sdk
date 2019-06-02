using System;
using Irony.Parsing;
using Strict.Statements;

namespace Strict.CollapsedGrammar
{
	internal class BinaryOperationParser : StatementParser
	{
		public BinaryOperationParser(StrictGrammar grammar, Parser parser) : base(grammar, parser) {}
		public override bool CanParse(ParseTreeNode element) => element.Term == grammar.BinaryOperation;

		public override Statement Parse(ParseTreeNode element)
		{
			var binaryOperator = GetBinaryOperator(ParseName(element.ChildNodes[1].ChildNodes[0]));
			return new BinaryOperation(parser.ParseNode(element.ChildNodes[0]), binaryOperator,
				parser.ParseNode(element.ChildNodes[2]));
		}

		private static BinaryOperator GetBinaryOperator(string operatorText)
		{
			switch (operatorText)
			{
			case "+":
				return BinaryOperator.Add;
			case "-":
				return BinaryOperator.Subtract;
			case "*":
				return BinaryOperator.Multiply;
			case "/":
				return BinaryOperator.Divide;
			case "%":
				return BinaryOperator.Modulate;
			case "and":
				return BinaryOperator.And;
			case "or":
				return BinaryOperator.Or;
			case "is":
				return BinaryOperator.Is;
			case "isnot":
				return BinaryOperator.IsNot;
			case "<":
				return BinaryOperator.Smaller;
			case ">":
				return BinaryOperator.Bigger;
			}
			//ncrunch: no coverage start
			throw new BinaryOperatorNotSupported(operatorText);
		}

		private class BinaryOperatorNotSupported : Exception
		{
			public BinaryOperatorNotSupported(string operatorText) : base(operatorText) {}
		}
	}
}