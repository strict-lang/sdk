using Irony.Parsing;
using Strict.Statements;

namespace Strict.CollapsedGrammar
{
	internal class BoolParser : StatementParser
	{
		public BoolParser(StrictGrammar grammar, Parser parser) : base(grammar, parser) {}
		public override bool CanParse(ParseTreeNode element) => element.Term == grammar.Bool;

		public override Statement Parse(ParseTreeNode element)
			=> new Bool(element.ChildNodes[0].Token.Value.ToString() == "true");
	}
}