using Irony.Parsing;
using Strict.Statements;

namespace Strict.CollapsedGrammar
{
	internal class MethodParser : StatementParser
	{
		public MethodParser(StrictGrammar grammar, Parser parser) : base(grammar, parser) {}
		public override bool CanParse(ParseTreeNode element) => element.Term == grammar.Method;

		public override Statement Parse(ParseTreeNode element)
		{
			var parameterElements = element.ChildNodes[3].ChildNodes;
			var parameters = new Parameter[parameterElements.Count];
			for (int i = 0; i < parameters.Length; i++)
				parameters[i] = (Parameter)parser.ParseNode(parameterElements[i]);
			var newMethod = new Method(parser.CurrentContext, ParseType(element.ChildNodes[0]),
				ParseName(element.ChildNodes[1]), parameters);
			return parser.ParseMethodStatementsInItsScope(newMethod, element.ChildNodes[2]);
		}
	}
}