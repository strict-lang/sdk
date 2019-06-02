using System;
using System.Collections.Generic;
using System.Linq;

namespace Strict.Statements
{
	/// <summary>
	/// Types are loaded directly from .type files via <see cref="Parsing.TypeFileParser" />, which
	/// do not contain any statements, just members. Types are always written in lower case letters,
	/// like most things in Strict. No OOP is used, see https://strict.fogbugz.com/default.asp?W11
	/// Type names are unique, generics are not supported in the long run, but currently for
	/// supporting DeltaEngine they are via &lt;&gt; after the name, <see cref="GenericType" />.
	/// </summary>
	public class Type : Statement, IDisposable
	{
		public Type(Context inContext, string name) : base(null)
		{
			if (inContext.ContainsType(name))
				throw new CannotCreateAlreadyExistingType(name);
			if (inContext.IsBase && inContext.NumberOfTypes >= Base.NumberOfAllowedBaseTypes)
				throw new CannotCreateTypeInBaseContext(name);
			if (inContext.ContainsChildContext(char.ToUpper(name[0]) + name.Substring(1)))
				throw new Context.TypeWithThisContextNameAlreadyExists(name);
			Context = inContext;
			Name = name;
			Context.Add(this);
		}

		public class CannotCreateAlreadyExistingType : Exception
		{
			public CannotCreateAlreadyExistingType(string name) : base(name) {}
		}

		public class CannotCreateTypeInBaseContext : Exception
		{
			public CannotCreateTypeInBaseContext(string name) : base(name) {}
		}

		public class TypesMustStartWithLowercaseLetter : Exception
		{
			public TypesMustStartWithLowercaseLetter(string name) : base(name) {}
		}

		public Context Context { get; }
		public string Name { get; }
		/// <summary>
		/// Support for C# Base Types, not used in Strict, but needed for importing and exporting C#.
		/// </summary>
		public string BaseTypeNames { get; set; }
		public IReadOnlyList<Member> Members => members;
		private readonly List<Member> members = new List<Member>();
		public IReadOnlyList<Statement> UsedBy => usedBy;
		internal readonly List<Statement> usedBy = new List<Statement>();
		
		public void Add(Member member)
		{
			members.Add(member);
			base.Add(member);
		}
		
		public Method GetBinaryMethod(BinaryOperator binaryOperator)
		{
			string operatorMethodName = GetBinaryMethodName(binaryOperator);
			Method existingMethod = methods.FirstOrDefault(m => m.Name == operatorMethodName);
			if (existingMethod != null)
				return existingMethod;
			const string LeftArgumentName = "left";
			const string RightArgumentName = "right";
			return new Method(Context, this, operatorMethodName, new Parameter(this, LeftArgumentName),
				new Parameter(this, RightArgumentName));
		}

		private string GetBinaryMethodName(BinaryOperator binaryOperator)
		{
			switch (binaryOperator)
			{
			case BinaryOperator.Add:
				return "add";
			case BinaryOperator.Subtract:
				return "subtract";
			case BinaryOperator.Multiply:
				return "multiply";
			case BinaryOperator.Divide:
				return "divide";
			case BinaryOperator.Modulate:
				return "modulate";
			case BinaryOperator.And:
				return "and";
			case BinaryOperator.Or:
				return "or";
			case BinaryOperator.Is:
				return "is";
			case BinaryOperator.IsNot:
				return "isnot";
			case BinaryOperator.Smaller:
				return "smaller";
			case BinaryOperator.Bigger:
				return "bigger";
			default:
				throw new OperatorIsNotSupportedForThisType(binaryOperator.ToString(), this);
			}
		}
		
		public class OperatorIsNotSupportedForThisType : Exception
		{
			public OperatorIsNotSupportedForThisType(string operatorName, Type type)
				: base(operatorName + " " + type) {}
		}

		public Method GetNegateMethod()
		{
			const string NegateOperatorName = "negate";
			if (this != Base.Number && this != Base.Bool)
				throw new OperatorIsNotSupportedForThisType(NegateOperatorName, this);
			return methods.FirstOrDefault(m => m.Name == NegateOperatorName) ??
				new Method(Context, this, NegateOperatorName, new Parameter(this, "argument"));
		}

		public override string ToString() => Context.IsBase ? Name : Context + "." + Name;

		public void Dispose()
		{
			if (Context.IsBase)
				throw new BaseContextTypesShouldNotBeDisposed(this);
			if (methods.Count > 0)
			{
				foreach (var method in methods.ToList())
					method.Dispose();
				methods.Clear();
			}
			Context.Remove(this);
		}

		public class BaseContextTypesShouldNotBeDisposed : Exception
		{
			public BaseContextTypesShouldNotBeDisposed(Type type) : base(type.ToString()) {}
		}

		public bool IsDisposed => !Context.ContainsType(Name);

		public override bool Equals(Statement other)
		{
			var otherType = other as Type;
			return otherType != null && Context == otherType.Context && Name == otherType.Name;
		}
		
		/// <summary>
		/// Adds a method that has this Type as the return type. Methods can be defined anywhere, types
		/// have no methods and the context does not have to match.
		/// </summary>
		internal void Add(Method method) => methods.Add(method);
		internal void Remove(Method method) => methods.Remove(method);
		private readonly List<Method> methods = new List<Method>();

		internal bool TryGetMethod(string methodName, Statement[] methodCallArguments,
			out Method foundMethod)
		{
			foundMethod =
				methods.FirstOrDefault(m => m.Name == methodName && m.DoArgumentsMatch(methodCallArguments));
			return foundMethod != null;
		}

		public void RemoveMethodsFromContext(Context context)
		{
			if (methods.Count == 0)
				return;
			var methodsToRemove = new List<Method>();
			foreach (var method in methods)
				if (method.Context == context)
					methodsToRemove.Add(method);
			foreach (var method in methodsToRemove)
				methods.Remove(method);
		}

		public IEnumerable<Method> RegisteredMethods => methods;
	}
}