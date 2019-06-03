using System;
using NUnit.Framework;
using Strict.Parsing;

namespace Strict.Tests.Parsing
{
	public class MethodOrClassGrabberTests
	{
		[Test]
		public void GetMethodLinesSimple()
		{
			var code = CreateFullSourceCode(SimpleMethodCode).Split(new[] { Environment.NewLine },
				StringSplitOptions.None);
			Assert.That(new MethodOrClassGrabber(code).GetCode("OnePlusOneMethod"),
				Is.EqualTo(SimpleMethodCode));
		}

		private const string SimpleMethodCode = @"		private static int OnePlusOneMethod()
		{
			return 1 + 1;
		}
";
		private static string CreateFullSourceCode(string methodCode)
		{
			return @"using System;
namespace Strict
{
	public class Test
	{
" + methodCode + @"
	}
}";
		}

		[Test]
		public void GetMethodLinesAdvanced()
		{
			var code = AdvancedFullSourceCode.Split(new[] { Environment.NewLine },
				StringSplitOptions.None);
			Assert.That(new MethodOrClassGrabber(code).GetCode("MethodWithArguments"),
				Is.EqualTo(MethodWithArgumentsCode));
		}

		private const string MethodWithArgumentsCode =
			@"		private static int MethodWithArguments(string name, int num)
		{
			int returnNumber = num;
			foreach (var c in name)
				returnNumber++;
			return returnNumber;
		}
";
		private const string AdvancedFullSourceCode = @"using System;
namespace Strict
{
	public class Test
	{
" + SimpleMethodCode + @"
" + MethodWithArgumentsCode + @"
	}
}";

		[Test]
		public void GetMethodLinesWithNCrunchComments()
		{
			var code =
				CreateFullSourceCode(SimpleMethodCodeWithNCrunchComments).Split(
					new[] { Environment.NewLine }, StringSplitOptions.None);
			Assert.That(new MethodOrClassGrabber(code).GetCode("OnePlusOneMethod"),
				Is.EqualTo(SimpleMethodCode));
		}

		private const string SimpleMethodCodeWithNCrunchComments =
@"		private static int OnePlusOneMethod()
		{
			//ncrunch: no coverage start
			return 1 + 1;
		} //ncrunch: no coverage end
";
	}
}