using Irony.Parsing;
using Strict.Statements;

namespace Strict.CollapsedGrammar
{
	internal abstract class StatementParser
	{
		protected StatementParser(StrictGrammar grammar, Parser parser)
		{
			this.grammar = grammar;
			this.parser = parser;
		}

		protected readonly StrictGrammar grammar;
		protected readonly Parser parser;
		public abstract bool CanParse(ParseTreeNode element);
		public abstract Statement Parse(ParseTreeNode element);
		protected static Type ParseType(ParseTreeNode childNode) => Base.Void;//TODO: use parser.CurrentContext, find type and return it for use!
		protected static string ParseName(ParseTreeNode childNode) => childNode.Token.Text;
	}
}