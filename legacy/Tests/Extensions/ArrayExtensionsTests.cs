using System.Collections.Generic;
using NUnit.Framework;
using Strict.Extensions;

namespace Strict.Tests.Extensions
{
	public class ArrayExtensionsTests
	{
		[SetUp]
		public void SetUp()
		{
			dictionary = new Dictionary<string, object> { { "int", 1 }, { "string", "string" } };
		}

		private Dictionary<string, object> dictionary;

		[Test]
		public void Compare()
		{
			var numbers1 = new[] { 1, 2, 5 };
			var numbers2 = new[] { 1, 2, 5 };
			var numbers3 = new[] { 1, 2, 5, 7 };
			Assert.That(numbers1.Compare(numbers2), Is.True);
			Assert.That(numbers1.Compare(null), Is.False);
			Assert.That(numbers1.Compare(numbers3), Is.False);
			Assert.That(numbers3.Compare(numbers1), Is.False);
			byte[] optionalData = null;
			// ReSharper disable once ExpressionIsAlwaysNull
			Assert.That(optionalData.Compare(null), Is.True);
		}

		[Test]
		public void Combine()
		{
			int[] numbers1 = { 1, 2, 3 };
			int[] numbers2 = { 4, 5, 6 };
			CollectionAssert.AreEqual(numbers1.Combine(numbers2), new[] { 1, 2, 3, 4, 5, 6 });
			CollectionAssert.AreEqual(numbers1.Combine(numbers1), numbers1);
			CollectionAssert.AreEqual(numbers1.Combine(null), numbers1);
			CollectionAssert.AreEqual(ArrayExtensions.Combine(null, numbers1), numbers1);
		}

		[Test]
		public void ToText()
		{
			var texts = new List<string> { "Hi", "there", "whats", "up?" };
			Assert.That(texts.ToText(), Is.EqualTo("Hi, there, whats, up?"));
			Assert.That(texts.ToText(", ", 2), Is.EqualTo("Hi, there"));
			Assert.That((null as List<string>).ToText(), Is.EqualTo(""));
		}

		[Test]
		public void GetWithDefaultReturnsDefaultIfNotInDictionary()
		{
			int result = ArrayExtensions.GetWithDefault<string, int>(dictionary, "Missing");
			Assert.That(result, Is.EqualTo(0));
		}

		[Test]
		public void GetWithDefaultReturnsValueIfFound()
		{
			int result = ArrayExtensions.GetWithDefault<string, int>(dictionary, "int");
			Assert.That(result, Is.EqualTo(1));
		}

		[Test]
		public void GetWithDefaultReturnsDefaultIfValueIsWrongType()
		{
			int result = ArrayExtensions.GetWithDefault<string, int>(dictionary, "string");
			Assert.That(result, Is.EqualTo(0));
		}
		
		[Test]
		public void Insert()
		{
			int[] source = { 1, 2, 4, 5, 6 };
			source = source.Insert(3, 2);
			Assert.That(source.Compare(new[] { 1, 2, 3, 4, 5, 6 }), Is.True);
			source = source.Insert(0, 0);
			Assert.That(source.Compare(new[] { 0, 1, 2, 3, 4, 5, 6 }), Is.True);
			source = source.Insert(7, source.Length);
			Assert.That(source.Compare(new[] { 0, 1, 2, 3, 4, 5, 6, 7 }), Is.True);
		}
	}
}