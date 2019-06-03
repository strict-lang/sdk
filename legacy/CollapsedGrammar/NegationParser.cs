using Irony.Parsing;
using Strict.Statements;

namespace Strict.CollapsedGrammar
{
	internal class NegationParser : StatementParser
	{
		public NegationParser(StrictGrammar grammar, Parser parser) : base(grammar, parser) {}
		public override bool CanParse(ParseTreeNode element) => element.Term == grammar.Negation;

		public override Statement Parse(ParseTreeNode element)
			=> new Negation(parser.ParseNode(element.ChildNodes[1]));
	}
}