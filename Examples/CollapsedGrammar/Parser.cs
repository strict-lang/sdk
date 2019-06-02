using System;
using System.Collections.Generic;
using Irony;
using Irony.Parsing;
using Strict.Statements;

namespace Strict.CollapsedGrammar
{
	/// <summary>
	/// Abstracts away the Irony Grammar and Parser usage to be used from the outside with our Strict
	/// Statements. Uses all the StatementParser classes defined in this namespace for the parsing.
	/// </summary>
	public class Parser
	{
		public Parser(Context currentContext)
		{
			if (currentContext.IsBase)
				throw new CannotParseFilesInBaseContext();
			CurrentContext = currentContext;
			ironyParser = new Irony.Parsing.Parser(grammar);
			statementParsers.Add(new NumberParser(grammar, this));
			statementParsers.Add(new TextParser(grammar, this));
			statementParsers.Add(new BoolParser(grammar, this));
			statementParsers.Add(new MemberParser(grammar, this));
			statementParsers.Add(new MethodCallParser(grammar, this));
			statementParsers.Add(new BinaryOperationParser(grammar, this));
			statementParsers.Add(new NegationParser(grammar, this));
			statementParsers.Add(new MethodParser(grammar, this));
			statementParsers.Add(new ParameterParser(grammar, this));
		}

		public class CannotParseFilesInBaseContext : Exception {}

		internal Context CurrentContext { get; private set; }
		private readonly StrictGrammar grammar = new StrictGrammar();
		private readonly Irony.Parsing.Parser ironyParser;
		private readonly List<StatementParser> statementParsers = new List<StatementParser>();

		public IReadOnlyList<Statement> Parse(string input)
		{
			ParseTree tree = ironyParser.Parse(input);
			if (tree.HasErrors())
				throw new ParsingFailed(input, tree.ParserMessages);
			var result = new List<Statement>();
			foreach (var element in tree.Root.ChildNodes)
				result.Add(ParseNode(element));
			return result;
		}

		public class ParsingFailed : Exception
		{
			public ParsingFailed(string input, LogMessageList messages)
				: base("\n" + input + "\n" + ShowSourcePointer(messages[0].Location) + "\n" +
					ShowParserMessages(messages)) {}

			private static string ShowSourcePointer(SourceLocation location)
				=> new string(' ', location.Column) + "^";

			private static string ShowParserMessages(LogMessageList parserMessages)
			{
				string result = "";
				foreach (var message in parserMessages)
					result += message.Location + " " + message.Message + (result != "" ? "\n" : "");
				return result;
			}
		}

		internal Statement ParseNode(ParseTreeNode element)
		{
			while (element.Term == grammar.AnyStatement || element.Term == grammar.Expression)
				element = element.ChildNodes[0];
			foreach (var statementParser in statementParsers)
				if (statementParser.CanParse(element))
					return statementParser.Parse(element);
			throw new NodeNotSupported(element); //ncrunch: no coverage start, only happens internally
		}

		private class NodeNotSupported : Exception
		{
			public NodeNotSupported(ParseTreeNode element) : base(element.ToString()) {}
		}

		public Method ParseMethodStatementsInItsScope(Method newMethod, ParseTreeNode methodStatements)
		{
			var previousScope = CurrentContext;
			CurrentContext = newMethod.Scope;
			foreach (var statement in methodStatements.ChildNodes)
				newMethod.Add(ParseNode(statement));
			CurrentContext = previousScope;
			return newMethod;
		}
	}
}