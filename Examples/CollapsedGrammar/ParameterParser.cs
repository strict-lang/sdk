using Irony.Parsing;
using Strict.Statements;

namespace Strict.CollapsedGrammar
{
	internal class ParameterParser : StatementParser
	{
		public ParameterParser(StrictGrammar grammar, Parser parser) : base(grammar, parser) {}
		public override bool CanParse(ParseTreeNode element) => element.Term == grammar.Parameter;

		public override Statement Parse(ParseTreeNode element)
		{
			var typeName = ParseName(element.ChildNodes[0]);
			var identifier = typeName;
			if (element.ChildNodes.Count > 1)
				identifier = ParseName(element.ChildNodes[1]);
			return new Parameter(parser.CurrentContext.GetType(typeName), identifier);
		}
	}
}