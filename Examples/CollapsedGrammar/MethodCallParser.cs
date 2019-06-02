using Irony.Parsing;
using Strict.Statements;

namespace Strict.CollapsedGrammar
{
	internal class MethodCallParser : StatementParser
	{
		public MethodCallParser(StrictGrammar grammar, Parser parser) : base(grammar, parser) {}
		public override bool CanParse(ParseTreeNode element) => element.Term == grammar.MethodCall;

		public override Statement Parse(ParseTreeNode element)
		{
			var firstName = ParseName(element.ChildNodes[0]);
			bool usesContext = ParseName(element.ChildNodes[1]) == ".";
			var methodName = usesContext ? ParseName(element.ChildNodes[2]) : firstName;
			var arguments =
				element.ChildNodes[element.ChildNodes.Count - 2].ChildNodes.ConvertAll(parser.ParseNode).
					ToArray();
			var context = usesContext ? Base.GetContext(firstName) : parser.CurrentContext;
			return new MethodCall(context.GetMethod(methodName, arguments), arguments);
		}
	}
}