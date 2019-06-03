using NUnit.Framework;
using Strict.Extensions;

namespace Strict.Tests.Extensions
{
	public class StringExtensionsTests
	{
		[Test]
		public void ConvertFloatToInvariantString()
		{
			Assert.That(2.5f.ToInvariantString(), Is.EqualTo("2.5"));
			Assert.That(1.5f.ToInvariantString(), Is.Not.EqualTo("3.5"));
			Assert.That(1.0f.ToInvariantString("00"), Is.EqualTo("01"));
			Assert.That(1.2345f.ToInvariantString("0.00"), Is.EqualTo("1.23"));
			Assert.That(1.2345f.ToInvariantString("0.00"), Is.Not.EqualTo("1.2345"));
			Assert.That(StringExtensions.ToInvariantString(1), Is.EqualTo("1"));
		}

		[Test]
		public void MaxStringLength()
		{
			Assert.That(((string)null).MaxStringLength(4), Is.EqualTo(null));
			Assert.That("".MaxStringLength(4), Is.EqualTo(""));
			Assert.That("abc".MaxStringLength(1), Is.EqualTo(".."));
			Assert.That("abcd".MaxStringLength(4), Is.EqualTo("abcd"));
			Assert.That("abcde".MaxStringLength(4), Is.EqualTo("ab.."));
			Assert.That("abcde".MaxStringLength(2), Is.EqualTo(".."));
		}

		[Test]
		public static void SplitAndTrimByChar()
		{
			string[] components = "abc, 123, def".SplitAndTrim(',');
			Assert.That(3, Is.EqualTo(components.Length));
			Assert.That("abc", Is.EqualTo(components[0]));
			Assert.That("123", Is.EqualTo(components[1]));
			Assert.That("def", Is.EqualTo(components[2]));
		}

		[Test]
		public static void SplitAndTrimByString()
		{
			string[] components = "3 plus 5 is 8".SplitAndTrim("plus", "is");
			Assert.That(3, Is.EqualTo(components.Length));
			Assert.That("3", Is.EqualTo(components[0]));
			Assert.That("5", Is.EqualTo(components[1]));
			Assert.That("8", Is.EqualTo(components[2]));
		}

		[Test]
		public void Compare()
		{
			Assert.That("AbC1".Compare("aBc1"), Is.True);
			Assert.That("1.23".Compare("1.23"), Is.True);
			Assert.That("Hello".Compare("World"), Is.False);
		}

		[Test]
		public void ContainsCaseInsensitive()
		{
			Assert.That("hallo".ContainsCaseInsensitive("ha"), Is.True);
			Assert.That("1.23".ContainsCaseInsensitive("1.2"), Is.True);
			Assert.That("Hello".ContainsCaseInsensitive("hel"), Is.True);
			Assert.That("Banana".ContainsCaseInsensitive("Apple"), Is.False);
			Assert.That(((string)null).ContainsCaseInsensitive("abc"), Is.False);
		}
		
		[Test]
		public void SplitTextIntoWords()
		{
			var stringWords = "TestForSplittingAStringIntoWords".SplitWords(true);
			Assert.That(stringWords[4], Is.EqualTo(' '));
			Assert.That(stringWords, Is.EqualTo("Test for splitting a string into words"));
		}

		[Test]
		public void FromByteArray()
		{
			const string TestString = "TestString";
			Assert.That(StringExtensions.FromByteArray(StringExtensions.ToByteArray(TestString)),
				Is.EqualTo(TestString));
		}

		[Test]
		public void StartsWith()
		{
			Assert.That(StringExtensions.StartsWith("Hi there, what's up?", "Hi"), Is.True);
			Assert.That(StringExtensions.StartsWith("Hi there, what's up?", "what"), Is.False);
			Assert.That(StringExtensions.StartsWith("bcdeuf", "bc"), Is.True);
			Assert.That(StringExtensions.StartsWith("bcdeuf", "abc"), Is.False);
			Assert.That("Hi there, what's up?".StartsWith("Hi", "there", "what"), Is.True);
			Assert.That("Hi there, what's up?".StartsWith("she", "there", "what"), Is.False);
			Assert.That(StringExtensions.StartsWith("ATI Radeon 9500+", "ati"), Is.True);
			Assert.That(StringExtensions.StartsWith("A-t-i da gaga", "ati"), Is.False);
			Assert.That(StringExtensions.StartsWith("NVidia Geforce3", "nvidia"), Is.True);
		}

		[Test]
		public void SplitWordsWithEmptyWord()
		{
			Assert.That("".SplitWords(true), Is.EqualTo(""));
		}
		
		[Test]
		public void TextToInt()
		{
			Assert.That("13".TryParse(0), Is.EqualTo(13));
			Assert.That("-100".TryParse(0), Is.EqualTo(-100));
			Assert.That("abc".TryParse(5), Is.EqualTo(5));
		}

		[Test]
		public void TextToFloat()
		{
			Assert.That("1.3".TryParse(0.0f), Is.EqualTo(1.3f));
			Assert.That("-1.0".TryParse(0.0f), Is.EqualTo(-1.0f));
			Assert.That("abc".TryParse(2.5f), Is.EqualTo(2.5f));
		}
	}
}