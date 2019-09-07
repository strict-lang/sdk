# Type Inference in Strict

Example code:

```strict
method String concat(String[] strings)
	let builder = StringBuilder()
	for string in strings do
		builder.Append(string)
	return builder.ToString()
```

So now there are two fields, who's values have to be deduced (builder and string).
They are both fairly simple to infer, since the assignment of `builder` has a
StringBuilder constructor as its RHS. string on the other hand has the value
of the String List `strings`. Lists are builtin types but still implement the
`Sequence<T>` interface. Thus the code could be written as

```strict
method String concat(Sequence<String> strings)
	let builder = StringBuilder()
	for string in strings do
		builder.Append(string)
	return builder.ToString()
```

Now instead of using the type parameter of the builtin List type, we use that
of the `Sequence` Type given. The Foreach statement can only have a `Sequence`
to iterate thus implementing inference for this field is still simple.

What, if the expression is slightly more advanced:

```strict
let sum = a + b
```

In this case, the type of `a` and `b` have to be known first. Then, the `add`
method is looked up for their type. Even primitives will have this method. 
To be more precise, the method is looked up on the `typeof(a)` (LHS) of a 
unary expression. If it exists, it has to have a single argument of 
type `typeof(b)` or a super-type of `typeof(b)`. 


