using System;
using System.Linq;

namespace Strict.Parsing
{
	/// <summary>
	/// Simple implementation to grab lines from source code file without parsing. Used to build
	/// source code files from methods in SourceCodeGrabber. Only used from services tests projects.
	/// </summary>
	public class MethodOrClassGrabber
	{
		public MethodOrClassGrabber(string[] lines)
		{
			this.lines = lines;
		}

		private readonly string[] lines;

		public string GetCode(string methodOrClass)
		{
			methodOrClass = methodOrClass.Replace("`1", "");
			code = "";
			currentAttributes = "";
			foreach (var line in lines)
				if (ProcessLine(methodOrClass, line))
					break;
			return code;
		}
		
		private string code;
		private string currentAttributes;

		private bool ProcessLine(string methodOrClass, string line)
		{
			if (!inBlock && IsAttribute(line))
				currentAttributes += line + Environment.NewLine;
			if (line.Contains("@\""))
				inMultilineString = true;
			if (inMultilineString && line.Replace("\"\"", "").Replace("@\"", "").Contains("\""))
				inMultilineString = false;
			if (!inMultilineString && (IsMethod(methodOrClass, line) || IsClass(methodOrClass, line)))
				inBlock = true;
			if (inBlock)
				code += currentAttributes + GetLine(line);
			if (inBlock || line.Trim().Length == 0 || line.Trim().StartsWith("{"))
				currentAttributes = "";
			return LastLineOfMethod;
		}

		private int bracketCounter;
		private bool inBlock;
		private bool inMultilineString;

		private static bool IsAttribute(string line)
		{
			return line.TrimStart().StartsWith("[") || line.TrimEnd().StartsWith("]");
		}

		private static bool IsMethod(string methodName, string line)
		{
			return line.Contains(" " + methodName + "(") && !line.TrimEnd().EndsWith(";");
		}

		private static bool IsClass(string className, string line)
		{
			return line.Contains(className) &&
				(line.Contains(" class ") || line.Contains(" struct ") || line.Contains("\tclass ") ||
				line.Contains("\tstruct "));
		}

		private string GetLine(string line)
		{
			bracketCounter += line.Sum(c => c == '{' ? 1 : 0) - line.Sum(c => c == '}' ? 1 : 0);
			if (line.Trim().StartsWith("//ncrunch: no coverage"))
				return "";
			return line.Replace("//ncrunch: no coverage end", "").TrimEnd() + Environment.NewLine;
		}

		private bool LastLineOfMethod
			=>
				inBlock && bracketCounter <= 0 &&
				(code.TrimEnd().EndsWith("}") || code.TrimStart().StartsWith("}"));
	}
}