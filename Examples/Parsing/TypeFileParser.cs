using System;
using System.Collections.Generic;
using System.IO;
using Strict.Statements;
using Type = Strict.Statements.Type;

namespace Strict.Parsing
{
	/// <summary>
	/// Loads a type with all its members from a file on disk.
	/// </summary>
	public class TypeFileParser : Type
	{
		public TypeFileParser(Context context, string filePath)
			: base(context, GetNameFromFilePathAndCheckValidity(context, filePath))
		{
			Parse(File.ReadAllLines(filePath));
		} //ncrunch: no coverage, only an issue with NCrunch 2015

		internal static string GetNameFromFilePathAndCheckValidity(Context context, string filePath,
			string expectedExtension = ".type")
		{
			if (!filePath.EndsWith(expectedExtension))
				throw new FilenameNeedsToBeValidAndEndWith(expectedExtension);
			var fullFilePath = new FileInfo(filePath);
			var currentDirectory = Directory.GetCurrentDirectory();
			if (Path.GetDirectoryName(fullFilePath.FullName) == currentDirectory)
				throw new FileMustBeInFolderForItsContext(filePath);
			if (!fullFilePath.FullName.StartsWith(currentDirectory))
				throw new FilePathMustBeBelowStrictExecutableFolder(filePath);
			if (context.FullName != fullFilePath.DirectoryName.Substring(currentDirectory.Length + 1))
				throw new ContextMustMatchPath(context, filePath);
			return Path.GetFileNameWithoutExtension(filePath);
		}

		public class FilenameNeedsToBeValidAndEndWith : Exception
		{
			public FilenameNeedsToBeValidAndEndWith(string expectedExtension)
				: base(expectedExtension) {}
		}

		public class FileMustBeInFolderForItsContext : Exception
		{
			public FileMustBeInFolderForItsContext(string filePath)
				: base(filePath) {}
		}

		public class FilePathMustBeBelowStrictExecutableFolder : Exception
		{
			public FilePathMustBeBelowStrictExecutableFolder(string filePath)
				: base(filePath) {}
		}

		public class ContextMustMatchPath : Exception
		{
			public ContextMustMatchPath(Context context, string filePath)
				: base("Content=" + context + ", filePath=" + filePath) {}
		}
		
		private void Parse(IReadOnlyCollection<string> memberLines)
		{
			if (memberLines.Count == 0)
				throw new ContentCannotBeEmpty(Context.Name);
			foreach (var line in memberLines)
				Add(ParseMember(line));
		}

		public class ContentCannotBeEmpty : Exception
		{
			public ContentCannotBeEmpty(string name)
				: base(name) {}
		}
		
		private Member ParseMember(string line)
		{
			if (string.IsNullOrWhiteSpace(line))
				throw new EmptyLinesAreNotAllowedInTypeFile(Name);
			var words = line.Split(' ');
			switch (words.Length)
			{
			case 1:
				return CreateMember(words[0], words[0]);
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
		
		public class InvalidTypeMemberSyntax : Exception
		{
			public InvalidTypeMemberSyntax(string name, string line)
				: base(name + "\n" + line) {}
		}

		public TypeFileParser(Context context, string name, string[] memberLines)
			: base(context, name)
		{
			Parse(memberLines);
		}
	}
}