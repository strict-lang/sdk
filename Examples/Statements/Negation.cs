namespace Strict.Statements
{
	public class Negation : MethodCall
	{
		public Negation(Statement argument) : base(argument.ReturnType.GetNegateMethod(), argument) {}
		public override string ToString() => (ReturnType == Base.Number ? "-" : "not ") + Argument;
		public string Operator => ReturnType == Base.Number ? "-" : "not";
		public Statement Argument => Arguments[0];
	}
}