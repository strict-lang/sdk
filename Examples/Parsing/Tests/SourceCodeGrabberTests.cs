using System;
using NUnit.Framework;

namespace Strict.Parsing.Tests
{
	public class SourceCodeGrabberTests
	{
		[Test]
		public void GetSourceCodeOfMethod()
		{
			var method = new Func<int>(OnePlusOneMethod);
			var sourceCode = new SourceCodeGrabber(method.Method);
			Assert.AreEqual(@"using System;

namespace Strict.Parsing.Tests
{
	public class Test
	{
		private static int OnePlusOneMethod()
		{
			return 1 + 1;
		}
	}
}", sourceCode.Code);
		}

		//ncrunch: no coverage start
		private static int OnePlusOneMethod()
		{
			return 1 + 1;
		} //ncrunch: no coverage end

		[Test]
		public void GetSourceCodeOfMethodWithAttributes()
		{
			var method = new Action(MethodWithAttributes);
			var sourceCode = new SourceCodeGrabber(method.Method);
			Assert.AreEqual(@"using System;

namespace Strict.Parsing.Tests
{
	public class Test
	{
		[Test, Ignore]
		public static void MethodWithAttributes()
		{
			Assert.AreEqual(2, OnePlusOneMethod());
		}
	}
}", sourceCode.Code);
		}

		//ncrunch: no coverage start
		[Test, Ignore]
		public static void MethodWithAttributes()
		{
			Assert.AreEqual(2, OnePlusOneMethod());
		} //ncrunch: no coverage end

		[Test]
		public void GetSourceCodeOfSingleClass()
		{
			var sourceCode = new SourceCodeGrabber(typeof(SampleTestClass));
			Assert.AreEqual(@"using System.Globalization;

namespace Strict.Parsing.Tests
{
	public class SampleTestClass
	{
		public SampleTestClass(int number, string text)
		{
			this.number = number;
			this.text = text;
		}

		private readonly int number;
		private readonly string text;

		public string PublicMethod()
		{
			return text + GetNumberText();
		}

		private string GetNumberText()
		{
			return number.ToString(CultureInfo.InvariantCulture);
		}
	}
}", sourceCode.Code);
		}

		[Test]
		public void GetSourceCodeOfNestedClass()
		{
			var sourceCode = new SourceCodeGrabber(typeof(NestedClass));
			Assert.AreEqual(@"using System;

namespace Strict.Parsing.Tests
{
		public class NestedClass : Object
		{
			private int number = 7;

			public int GetNumberAndIncreaseIt()
			{
				return number++;
			}
		}
}", sourceCode.Code);
		}

		//ncrunch: no coverage start
		public class NestedClass : Object
		{
			private int number = 7;

			public int GetNumberAndIncreaseIt()
			{
				return number++;
			}
		} //ncrunch: no coverage end
	}
}