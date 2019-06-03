using System;
using Irony.Parsing;
using Strict.Statements;

namespace Strict.CollapsedGrammar
{
	internal class NumberParser : StatementParser
	{
		public NumberParser(StrictGrammar grammar, Parser parser) : base(grammar, parser) {}
		public override bool CanParse(ParseTreeNode element) => element.Term == grammar.Number;

		public override Statement Parse(ParseTreeNode element)
		{
			var details = element.Token.Details as CompoundTerminalBase.CompoundTokenDetails;
			var numberType = details?.TypeCodes[0] ?? TypeCode.Int32;
			var number = element.Token.Value;
			return new Number(numberType == TypeCode.Int32 ? (int)number : (double)number);
		}
	}
}