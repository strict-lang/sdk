using System.Globalization;

namespace Strict.Tests.Parsing
{
	public class SampleTestClass
	{
		//ncrunch: no coverage start
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
}