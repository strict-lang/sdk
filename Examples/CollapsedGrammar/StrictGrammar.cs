using Irony.Parsing;

namespace Strict.CollapsedGrammar
{
	/// <summary>
	/// Parses collapsed strict language expressions and converts them into the strict statement tree.
	/// Basically very similar to Python with indentation and syntax, but overall more similar to C#.
	/// </summary>
	[Language("Strict", "0.1",
		"Parses collapsed strict language syntax into Strict Statements via the Parser class")]
	internal class StrictGrammar : Grammar
	{
		public StrictGrammar()
		{
			CreateTypes();
			CreateNonTerminals();
			CreateMemberRules();
			CreateMethodRules();
			CreateRootRules();
		}

		private void CreateTypes()
		{
			Number = TerminalFactory.CreateCSharpNumber("Number");
			Text = TerminalFactory.CreateCSharpString("StringText");
			Bool = new NonTerminal("Bool") { Rule = ToTerm("true") | "false" };
			Type = Identifier = TerminalFactory.CreateCSharpIdentifier("Identifier");
			Context = new NonTerminal("Context");
			Context.Rule = MakePlusRule(Context, ToTerm("."), Identifier);
		}

		public NumberLiteral Number { get; private set; }
		public StringLiteral Text { get; private set; }
		public NonTerminal Bool { get; private set; }
		/// <summary>
		/// Type or member or namespace, we don't know at parsing time from Irony. Could be anything.
		/// </summary>
		public IdentifierTerminal Identifier { get; private set; }
		public IdentifierTerminal Type { get; private set; }
		public NonTerminal Context { get; private set; }

		private void CreateNonTerminals()
		{
			Root = new NonTerminal("Root");
			Member = new NonTerminal("Member");
			Assignment = new NonTerminal("Assignment");
			Method = new NonTerminal("Method");
			Parameter = new NonTerminal("Parameter");
			Expression = new NonTerminal("Expression");
			MethodCall = new NonTerminal("MethodCall");
			BinaryOperation = new NonTerminal("BinaryOperation");
			Negation = new NonTerminal("Negation");
			AnyStatement = new NonTerminal("AnyStatement");
		}

		public NonTerminal Member { get; private set; }
		public NonTerminal Assignment { get; private set; }
		public NonTerminal Method { get; private set; }
		public NonTerminal Parameter { get; private set; }
		/// <summary>
		/// For an assignment or argument for a method call, can be any identifier, constant number,
		/// string or bool, but a MethodCall, BinaryOperation, UnaryOperation or Negation.
		/// </summary>
		public NonTerminal Expression { get; private set; }
		public NonTerminal MethodCall { get; private set; }
		public NonTerminal BinaryOperation { get; private set; }
		public NonTerminal Negation { get; private set; }
		public NonTerminal AnyStatement { get; private set; }

		private void CreateMemberRules()
		{
			Member.Rule = Type + Assignment;
			Assignment.Rule = Identifier + "=" + Expression;
			Expression.Rule = Number | Text | Bool | Identifier | MethodCall | BinaryOperation |
				Negation;
		}

		private void CreateMethodRules()
		{
			CreateMethodRule();
			CreateMethodCallRule();
			CreateBinaryOperationRule();
			CreateNegationRule();
		}

		/// <summary>
		/// Must match <see cref="Strict.Statements.BinaryOperator" /> and how
		/// <see cref="Strict.Statements.BinaryOperation" /> converts names back to the operator enum.
		/// </summary>
		private void CreateBinaryOperationRule()
		{
			var binaryOperator = new NonTerminal("binaryOperator");
			SetOperatorPrecedence();
			binaryOperator.Rule = ToTerm("+") | "-" | "*" | "/" | "%" | "and" | "or" | "isnot" | "is" |
				"<" | ">";
			BinaryOperation.Rule = Expression + ImplyPrecedenceHere(5) + binaryOperator + Expression;
		}

		private void SetOperatorPrecedence()
		{
			RegisterOperators(-1, "=");
			RegisterOperators(1, "or");
			RegisterOperators(2, "and");
			RegisterOperators(3, "is", "is not");
			RegisterOperators(4, "<", ">", "<=", ">=");
			RegisterOperators(5, "+", "-");
			RegisterOperators(6, "*", "/", "%");
			RegisterOperators(7, "%");
			RegisterOperators(8, "not");
		}
		
    private void CreateNegationRule()
			=> Negation.Rule = ToTerm("-") + Expression | ToTerm("not") + Expression;

		private void CreateMethodRule()
		{
			Parameter.Rule = Type + Identifier | Identifier;
			var optionalParameters = new NonTerminal("optionalParameters");
			optionalParameters.Rule = Empty | MakePlusRule(optionalParameters, ToTerm(","), Parameter);
			Method.Rule = Type + Identifier + "(" + optionalParameters + ")";
		}

		private void CreateMethodCallRule()
		{
			var optionalArguments = new NonTerminal("optionalArguments");
			optionalArguments.Rule = Empty | MakePlusRule(optionalArguments, ToTerm(","), Expression);
			MethodCall.Rule = Identifier + "(" + optionalArguments + ")" |
				Context + "." + Identifier + "(" + optionalArguments + ")";
		}

		private void CreateRootRules()
		{
			AnyStatement.Rule = Expression | Negation | Member | Method | Assignment | MethodCall;
			Root.Rule = MakeStarRule(Root, AnyStatement);
		}
	}
}