using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Reflection;

namespace Strict.Parsing
{
	/// <summary>
	/// Grabs source code from a given method or class type dynamically and builds compilable output,
	/// which can be passed to a compiler or AST parser for further conversion. Used in NCrunch tests.
	/// Usings are just copied from the source file, use separate files for reduced usings. 
	/// </summary>
	public class SourceCodeGrabber
	{
		public SourceCodeGrabber(MethodInfo method)
		{
			Code = GrabCode(method.DeclaringType, method.Name);
		}

		public string Code { get; private set; }
		public List<string> AssemblyReferenceNames { get; private set; }

		public SourceCodeGrabber(Type classType)
		{
			Code = GrabCode(classType, classType.Name);
		}

		private string GrabCode(Type classTypeForFilename, string methodOrClassName)
		{
			AssemblyReferenceNames = new List<string>();
			foreach (var reference in classTypeForFilename.Assembly.GetReferencedAssemblies())
				if (!IsServicesAssembly(reference) && !AssemblyReferenceNames.Contains(reference.Name))
					AssemblyReferenceNames.Add(reference.Name);
			var lines = GetFileAndReadAllLines(classTypeForFilename);
			return BuildClassCodeFromMethod(GetUsings(lines),
				new MethodOrClassGrabber(lines).GetCode(methodOrClassName));
		}

		private static bool IsServicesAssembly(AssemblyName reference)
		{
			return reference.Name == "nunit.framework" || reference.Name == "System.Core" ||
				reference.Name == "nCrunch.TestRuntime" || reference.Name.StartsWith("Strict") ||
				reference.Name == "Mono.Cecil" || reference.Name.StartsWith("ICSharpCode.NRefactory");
		}

		private string[] GetFileAndReadAllLines(Type classTypeForFilename)
		{
			FileName = classTypeForFilename.Name + ".cs";
			var path = Path.GetDirectoryName(classTypeForFilename.Assembly.Location);
			path = Path.Combine(path, "..", "..");
			var names = classTypeForFilename.FullName.Replace("`1", "").Split('+')[0].Split('.');
			var expectedPath = Path.Combine(path, names.Last() + ".cs");
			if (!File.Exists(expectedPath))
				expectedPath = Path.Combine(path,
					string.Join(Path.DirectorySeparatorChar.ToString(), names.Skip(1)) + ".cs");
			return File.ReadAllLines(expectedPath);
		}

		public string FileName { get; private set; }
		public string Namespace { get; private set; }

		private string BuildClassCodeFromMethod(string usings, string extractedCode)
		{
			const string TestClassCode = @"	public class Test
	{
";
			var hasClassInCode = extractedCode.Contains("class ") || extractedCode.Contains("struct ");
			return usings + @"namespace " + Namespace + @"
{
" + (hasClassInCode ? "" : TestClassCode) + extractedCode + (hasClassInCode ? "" : @"	}
") + "}";
		}

		private string GetUsings(string[] lines)
		{
			string usings = "";
			Namespace = "";
			foreach (var line in lines)
			{
				if (line.TrimStart().StartsWith("using ") && line.TrimEnd().EndsWith(";") &&
						!line.Contains("NUnit.Framework;"))
					usings += line + Environment.NewLine;
				if (line.TrimStart().StartsWith("namespace"))
					Namespace = line.Replace("namespace ", "").Replace("{", "").Replace("}", "").Trim();
				if (line.Contains("class ") || line.Contains("struct ") || line.Contains("enum"))
					break;
			}
			return usings != ""? usings + Environment.NewLine : usings;
		}
	}
}