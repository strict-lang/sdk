using System;
using System.Collections.Generic;
using Strict.Extensions;

namespace Strict.Statements
{
	public class Method : NamedStatement, IDisposable
	{
		public Method(Context inContext, Type returnType, string name, params Parameter[] parameters)
			: base(returnType, name)
		{
			Parameters = parameters;
			Context = inContext;
			if (Context.IsBase && !IsBaseTypeOperator)
				throw new CannotCreateMethodInBaseContext(name);
			Scope = new Context(inContext, name, true);
			returnType.Add(this);
		}
		
		public IReadOnlyList<Parameter> Parameters { get; }
		public Context Context { get; }
		public bool IsBaseTypeOperator
			=> Parameters.Count == 2 && Parameters[0].Name == "left" && Parameters[1].Name == "right" ||
				Parameters.Count == 1 && Parameters[0].Name == "argument";

		public class CannotCreateMethodInBaseContext : Exception
		{
			public CannotCreateMethodInBaseContext(string name) : base(name) {}
		}

		public Context Scope { get; }

		public void Add(Member member)
		{
			base.Add(member);
			Scope.AddOrReplace(member);
		}

		public void Dispose()
		{
			ReturnType.Remove(this);
			Scope.Dispose();
		}

		public bool DoArgumentsMatch(Statement[] methodCallArguments)
		{
			if (methodCallArguments != null && Parameters.Count != methodCallArguments.Length)
				return false;
			for (int i = 0; i < Parameters.Count; i++)
				if (Parameters[i].ReturnType != methodCallArguments[i].ReturnType)
					return false;
			return true;
		}

		public override string ToString() => ReturnType + " " + Scope + "(" + Parameters.ToText() + ")";
	}
}