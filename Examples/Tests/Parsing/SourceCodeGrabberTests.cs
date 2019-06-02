using System;
using NUnit.Framework;
using Strict.Parsing;

namespace Strict.Tests.Parsing
{
	public class SourceCodeGrabberTests
	{
		[Test]
		public void GetSourceCodeOfMethod()
		{
			var method = new Func<int>(OnePlusOneMethod);
			var sourceCode = new SourceCodeGrabber(method.Method);
			Assert.That(sourceCode.Code, Is.EqualTo(@"using System;
using Strict.Parsing;

namespace Strict.Tests.Parsing
{
	public class Test
	{
		private static int OnePlusOneMethod()
		{
			return 1 + 1;
		}
	}
}"));
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
			Assert.That(sourceCode.Code, Is.EqualTo(@"using System;
using Strict.Parsing;

namespace Strict.Tests.Parsing
{
	public class Test
	{
		[Test, Category(""Used to test multiple attributes"")]
		public static void MethodWithAttributes()
		{
			Assert.That(OnePlusOneMethod(), Is.EqualTo(2));
		}
	}
}"));
		}

		//ncrunch: no coverage start
		[Test, Category("Used to test multiple attributes")]
		public static void MethodWithAttributes()
		{
			Assert.That(OnePlusOneMethod(), Is.EqualTo(2));
		} //ncrunch: no coverage end

		[Test]
		public void GetSourceCodeOfSingleClass()
		{
			var sourceCode = new SourceCodeGrabber(typeof(SampleTestClass));
			Assert.That(sourceCode.Code, Is.EqualTo(@"using System.Globalization;

namespace Strict.Tests.Parsing
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
}"));
		}

		[Test]
		public void GetSourceCodeOfNestedClass()
		{
			var sourceCode = new SourceCodeGrabber(typeof(NestedClass));
			Assert.That(sourceCode.Code, Is.EqualTo(@"using System;
using Strict.Parsing;

namespace Strict.Tests.Parsing
{
		public class NestedClass : Object
		{
			private int number = 7;

			public int GetNumberAndIncreaseIt()
			{
				return number++;
			}
		}
}"));
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