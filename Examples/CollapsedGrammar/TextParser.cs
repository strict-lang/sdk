using Irony.Parsing;
using Strict.Statements;

namespace Strict.CollapsedGrammar
{
	internal class TextParser : StatementParser
	{
		public TextParser(StrictGrammar grammar, Parser parser) : base(grammar, parser) {}
		public override bool CanParse(ParseTreeNode element) => element.Term == grammar.Text;
		public override Statement Parse(ParseTreeNode element) => new Text((string)element.Token.Value);
	}
}