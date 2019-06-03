namespace Strict.Statements
{
	/// <summary>
	/// Allows to access base functionality of .NET and build in types and all GAC namespaces. The
	/// Base.Context is always available, but nothing can be added to it as no ones will to clean up.
	/// Create a child contexts instead. Also see http://en.wikipedia.org/wiki/Data_type
	/// </summary>
	public static class Base
	{
		static Base()
		{
			Context = new Context(null, string.Empty);
			Number = new Type(Context, "number");
			Text = new Type(Context, "text");
			Bool = new Type(Context, "bool");
			List = new Type(Context, "list");
			Map = new Type(Context, "map");
			Anything = new Type(Context, "anything");
			Void = new Type(Context, "void");
		}

		public static Context Context { get; }
		public static Type Number { get; }
		public static Type Text { get; }
		public static Type Bool { get; }
		public static Type List { get; }
		public static Type Map { get; }
		public static Type Anything { get; }
		public static Type Void { get; }
		internal const int NumberOfAllowedBaseTypes = 7;

		public static Context GetContext(string fullName)
		{
			return Context.GetChildContext(fullName.Split('.'));
		}
	}
}