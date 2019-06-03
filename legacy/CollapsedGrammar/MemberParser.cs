using Irony.Parsing;
using Strict.Statements;

namespace Strict.CollapsedGrammar
{
	internal class MemberParser : StatementParser
	{
		public MemberParser(StrictGrammar grammar, Parser parser) : base(grammar, parser) {}
		public override bool CanParse(ParseTreeNode element) => element.Term == grammar.Assignment;

		public override Statement Parse(ParseTreeNode element)
		{
			if (!parser.CurrentContext.IsInMethod)
				throw new Member.MembersMustBeDefinedInTypesOrMethods();
			var member = new Member(ParseName(element.ChildNodes[0]),
				parser.ParseNode(element.ChildNodes[2]));
			parser.CurrentContext.AddOrReplace(member);
			return member;
		}
	}
}