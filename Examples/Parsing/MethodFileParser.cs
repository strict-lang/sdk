using System;
using System.Collections.Generic;
using System.IO;
using Strict.Statements;
using Type = Strict.Statements.Type;

namespace Strict.Parsing
{
	/// <summary>
	/// Loads a method with all its statements from a file on disk.
	/// </summary>
	public class MethodFileParser : Method
	{
		public MethodFileParser(Context context, string filePath)
			: this(context, GetNameFromFilePathAndCheckValidity(context, filePath),
				File.ReadAllLines(filePath)) {}

		private static string GetNameFromFilePathAndCheckValidity(Context context, string filePath)
		{
			string filename = Path.GetFileNameWithoutExtension(filePath);
			if (!filename.Contains("(") && !filename.Contains(")"))
				throw new MethodFilenameMustContainBrackets();
			//TODO: must contain (), optionally parameters and return type if not void!
			//e.g. customAdd(number,number)number.method
			return TypeFileParser.GetNameFromFilePathAndCheckValidity(context, filePath, ".method");
		}

		public class MethodFilenameMustContainBrackets : Exception {}

		public MethodFileParser(Context context, string name, string[] fileLines)
			: base(context, ParseReturnType(context.Name + "." + name, fileLines), name,
					ParseParameters(context.Name+"."+name, fileLines))
		{
			ParseStatements(fileLines);
		}

		private static Type ParseReturnType(string fullMethodName, IReadOnlyCollection<string> fileLines)
		{
			if (fileLines.Count == 0)
				throw new TypeFileParser.ContentCannotBeEmpty(fullMethodName);
			return Base.Void;//TODO
		}

		private static Parameter[] ParseParameters(string fullMethodName, string[] fileLines)
		{
			return new Parameter[] { };
		}

		private static void ParseStatements(string[] fileLines) {}

		/*TODO
		private Member ParseStatement(string line)
		{
			if (string.IsNullOrWhiteSpace(line))
				throw new EmptyLinesAreNotAllowedInTypeFile(Name);
			var words = line.Split(' ');
			switch (words.Length)
			{
			case 1:
				return CreateMember(CapitalizeFirstLetter(words[0]), words[0]);
			case 2:
				return CreateMember(words[0], words[1]);
			case 3:
				if (words[1] == "=")
					return new Member(words[0], new Value(words[2]));
				break;
			}
			throw new InvalidTypeMemberSyntax(Name, line);
		}

		private Member CreateMember(string typeName, string name)
		{
			return new Member(name, new Value(Context.GetType(typeName)));
		}

		public class EmptyLinesAreNotAllowedInTypeFile : Exception
		{
			public EmptyLinesAreNotAllowedInTypeFile(string name)
				: base(name) {}
		}

		private static string CapitalizeFirstLetter(string name)
		{
			return char.ToUpper(name[0]) + name.Substring(1);
		}

		public class InvalidTypeMemberSyntax : Exception
		{
			public InvalidTypeMemberSyntax(string name, string line)
				: base(name + "\n" + line) {}
		}
	*/
	}
}