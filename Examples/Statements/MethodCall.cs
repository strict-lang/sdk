using System.Collections.Generic;
using System.Linq;
using Strict.Extensions;

namespace Strict.Statements
{
	public class MethodCall : Statement
	{
		public MethodCall(Method method, params Statement[] arguments)
			: base(method.ReturnType, arguments)
		{
			Method = method;
			Arguments = arguments;
		}

		public Method Method { get; }
		public IReadOnlyList<Statement> Arguments { get; }
		public override string ToString() => Method.Name + "(" + Arguments.ToText() + ")";

		public override bool Equals(Statement other)
		{
			var otherMethodCall = other as MethodCall;
			return otherMethodCall != null && Method.Equals(otherMethodCall.Method) &&
				Arguments.SequenceEqual(otherMethodCall.Arguments);
		}
	}
}