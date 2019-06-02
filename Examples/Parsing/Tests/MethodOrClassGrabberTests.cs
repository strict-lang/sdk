using System;
using NUnit.Framework;

namespace Strict.Parsing.Tests
{
	public class MethodOrClassGrabberTests
	{
		[Test]
		public void GetMethodLinesSimple()
		{
			var code = SimpleFullSourceCode.Split(new[] { Environment.NewLine }, StringSplitOptions.None);
			Assert.AreEqual(SimpleMethodCode, new MethodOrClassGrabber(code).GetCode("OnePlusOneMethod"));
		}

		private const string SimpleMethodCode = @"		private static int OnePlusOneMethod()
		{
			return 1 + 1;
		}
";
		private const string SimpleFullSourceCode = @"using System;
namespace Strict
{
	public class Test
	{
" + SimpleMethodCode + @"
	}
}";

		[Test]
		public void GetMethodLinesAdvanced()
		{
			var code = AdvancedFullSourceCode.Split(new[] { Environment.NewLine },
				StringSplitOptions.None);
			Assert.AreEqual(MethodWithArgumentsCode,
				new MethodOrClassGrabber(code).GetCode("MethodWithArguments"));
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
	}
}