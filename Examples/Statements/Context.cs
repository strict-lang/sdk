using System;
using System.Collections.Generic;
using System.Linq;

namespace Strict.Statements
{
	/// <summary>
	/// As described on https://strict.fogbugz.com/default.asp?W10#toc_3 contexts are scopes, they
	/// contain all types with its members and methods defined in it, which are saved as files in the
	/// context folder. You can compare contexts to namespaces in other languages, however a context
	/// can be in a method too, be nested and there will be a lot more than namespaces. Members that
	/// are not longer used at the end of a method scope are automatically cleaned up.
	/// </summary>
	public sealed class Context : IDisposable
	{
		public Context(Context parent, string name, bool isMethodScope = false)
		{
			if (name.Contains("."))
				throw new NameShouldNotContainDots(name);
			if (parent != null && string.IsNullOrEmpty(name))
				throw new ContextNameCannotBeEmpty();
			if (parent != null && !isMethodScope && !char.IsUpper(name[0]))
				throw new NameMustStartWithUpperCaseLetter(name);
			if (parent != null && !isMethodScope && parent.ContainsChildContext(name))
				throw new ContextAlreadyExists((parent.IsBase ? "" : parent + ".") + name);
			if (parent != null && parent.types.ContainsKey(char.ToLower(name[0]) + name.Substring(1)))
				throw new TypeWithThisContextNameAlreadyExists(name);
			parent?.children.Add(this);
			Parent = parent;
			Name = name;
			FullName = (Parent != null && !Parent.IsBase ? Parent + "." : "") + Name;
		}

		public class NameShouldNotContainDots : Exception
		{
			public NameShouldNotContainDots(string name) : base(name) {}
		}

		public class ContextNameCannotBeEmpty : Exception {}

		public class NameMustStartWithUpperCaseLetter : Exception
		{
			public NameMustStartWithUpperCaseLetter(string name) : base(name) {}
		}

		public class ContextAlreadyExists : Exception
		{
			public ContextAlreadyExists(string fullName) : base(fullName) {}
		}

		public class TypeWithThisContextNameAlreadyExists : Exception
		{
			public TypeWithThisContextNameAlreadyExists(string name) : base(name) {}
		}

		private readonly List<Context> children = new List<Context>();
		public Context Parent { get; }
		public string Name { get; }
		public string FullName { get; }
		public bool IsBase => string.IsNullOrEmpty(Name);
		public IEnumerable<Context> ChildContexts => children;

		public void Dispose()
		{
			Parent?.children.Remove(this);
			while (children.Count > 0)
			{
				var child = children[children.Count - 1];
				children.Remove(child);
				child.Dispose();
			}
			children.Clear();
			members.Clear();
			var checkContext = this;
			do
			{
				foreach (var type in checkContext.types.Values)
					type.RemoveMethodsFromContext(checkContext);
				checkContext = checkContext.Parent;
			} while (checkContext != null);
			foreach (var type in types.Values.ToList())
				type.Dispose();
			types.Clear();
		}

		public void AddOrReplace(Member member)
		{
			if (this == Base.Context)
				//ncrunch: no coverage start, can only happen internally
				throw new NotAllowedToAddDirectlyToBaseContextCreateChildContextInstead();
			//ncrunch: no coverage end
			if (members.ContainsKey(member.Name))
			{
				members[member.Name].ClearAndMakeUnusable();
				members[member.Name] = member;
			}
			else
				members.Add(member.Name, member);
		}

		public class NotAllowedToAddDirectlyToBaseContextCreateChildContextInstead : Exception {}

		private readonly Dictionary<string, Member> members = new Dictionary<string, Member>();
		/// <summary>
		/// Returns all members accessible in this scope and all parent scopes.
		/// </summary>
		public IEnumerable<Member> AllAccessibleMembersRecursively
		{
			get
			{
				foreach (var pair in members)
					yield return pair.Value;
				if (Parent != null)
					foreach (var method in Parent.AllAccessibleMembersRecursively)
						yield return method;
			}
		}

		internal bool ContainsChildContext(string name) => children.Any(child => child.Name == name);

		public void Add(Type type)
		{
			if (types.ContainsKey(type.Name))
				throw new TypeWasAlreadyAdded(type); //ncrunch: no coverage, can only happen internally
			types.Add(type.Name, type);
		}

		public class TypeWasAlreadyAdded : Exception
		{
			public TypeWasAlreadyAdded(Type type) : base(type.ToString()) {} //ncrunch: no coverage
		}

		private readonly Dictionary<string, Type> types = new Dictionary<string, Type>();
		public IEnumerable<Type> GetAllTypesRecursively()
		{
			foreach (var pair in types)
				yield return pair.Value;
			if (Parent != null)
				foreach (var pair in Parent.types)
					yield return pair.Value;
		}
		public int NumberOfTypes => types.Count;
		
		public bool ContainsType(string name)
		{
			if (string.IsNullOrEmpty(name) || char.IsUpper(name[0]))
				throw new Type.TypesMustStartWithLowercaseLetter(name);
			return types.ContainsKey(name) || Parent != null && Parent.ContainsType(name);
		}

		public Type GetType(string name)
		{
			Type type;
			if (types.TryGetValue(name, out type))
				return type;
			if (Parent != null)
				return Parent.GetType(name);
			throw new TypeNotFound(name);
		}
		
		public class TypeNotFound : Exception
		{
			public TypeNotFound(string typeName) : base(typeName) {}
		}

		internal void Remove(Type type) => types.Remove(type.Name);
		public override string ToString() => FullName;
		public bool IsInMethod => !string.IsNullOrEmpty(Name) && char.IsLower(Name[0]);

		internal Context GetChildContext(string[] nameParts)
		{
			foreach (var child in children)
				if (child != null && child.Name == nameParts[0])
					return nameParts.Length > 1 ? child.GetChildContext(nameParts.Skip(1).ToArray()) : child;
			throw new ChildContextNotFound(nameParts);
		}

		public class ChildContextNotFound : Exception
		{
			public ChildContextNotFound(string[] nameParts) : base(string.Join(".", nameParts)) {}
		}

		/// <summary>
		/// Finds the method in the given context (or parent contexts) that matches the given name
		/// and arguments. We don't know the return type yet, which is where this method lives.
		/// </summary>
		public Method GetMethod(string methodName, params Statement[] arguments)
		{
			Method method;
			var guessedReturnType = arguments != null && arguments.Length > 0
				? arguments[0].ReturnType : Base.Void;
			if (guessedReturnType.TryGetMethod(methodName, arguments, out method))
				return method;
			var context = this;
			do
			{
				if (context.FindMethodWithGivenArguments(methodName, arguments, out method))
					return method;
				context = context.Parent;
			} while (context != null);
			throw new MethodNotFound(ToString(), methodName);
		}

		private bool FindMethodWithGivenArguments(string methodName, Statement[] arguments,
			out Method method)
		{
			foreach (var type in types.Values)
				if (type.TryGetMethod(methodName, arguments, out method))
					return true;
			method = null;
			return false;
		}

		public class MethodNotFound : Exception
		{
			public MethodNotFound(string contentName, string methodName)
				: base(methodName + " not found in context: " + contentName) {}
		}

		public Context GetSubContextFromNamespace(string namespaceName)
		{
			if (!namespaceName.Contains('.'))
				return ChildContexts.FirstOrDefault(c => c.Name == namespaceName) ??
					new Context(this, namespaceName);
			var nameParts = namespaceName.Split(new[] { '.' }, 2);
			var subContext = ChildContexts.FirstOrDefault(c => c.Name == nameParts[0]) ??
				new Context(this, nameParts[0]);
			return subContext.GetSubContextFromNamespace(nameParts[1]);
		}

		public Type GetTypeInAnyChildContext(string name)
		{
			Type type;
			if (types.TryGetValue(name, out type))
				return type;
			foreach (var child in ChildContexts)
			{
				var childType = child.GetTypeInAnyChildContext(name);
				if (childType != null)
					return childType;
			}
			return null;
		}
	}
}