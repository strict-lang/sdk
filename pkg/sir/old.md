# Strict Intermediate Representation Specification

- [Introduction](#introduction)
- [Structure](#structure)
    - [Overview](#structure_overview)
    - [Constant Poool](#structure_constant_pool)
    - [Modules](#structure_modules)
    - [MetadataTables](#structure_metadata_table)
    - [Mehtod Declarations](#structure_method_declaration)
      - [ParameterList](#structure_params_list)
      - [Local Variable Table](#structure_method_vars)
- [Encoding](#encoding)
- [Metadata](#metadata)
- [Operations](#operations)
    - [Description Format](#operations_format)
    - [Parameterized Operations](#operations_parameterized)
    - [Operations](#operations_list)
      - [`add<>`](#op_add)
      - [`sub<>`](#op_sub)
      - [`mul<>`](#op_mul)
      - [`div<>`](#op_div)
      - [`mod<>`](#op_mod)
      - [`cmp<>`](#op_cmp)
      - [`cmp_eq<>`](#op_cmp_eq)
      - [`cmp_ne<>`](#op_cmp_ne)
      - [`cmp_lt<>`](#op_cmp_lt)
      - [`cmp_le<>`](#op_cmp_le)
      - [`cmp_gt<>`](#op_cmp_gt)
      - [`cmp_ge<>`](#op_cmp_ge)
      - [`call<>`](#op_call)
      
<a name="introduction"></a>
## 1 Introduction
The Strict intermediate representation (*SIR*) is a statically typed, Single Static Assignment (*SSA*) based IR, 
that represents compiled Strict code. It is low level compared to Strict, but still aware of high level concepts.

## 2 Structure
### 2.1 Constant Pool
### 2.2 Modules
Named values are represented as a string of characters with their prefix. For example, %foo, @DivisionByZero, %a.really.long.identifier. The actual regular expression used is ‘[%@][-a-zA-Z$._][-a-zA-Z$._0-9]*’. Identifiers that require other characters in their names can be surrounded with quotes. Special characters may be escaped using "\xx" where xx is the ASCII code for the character in hexadecimal. In this way, any character can be used in a name value, even quotes themselves. The "\01" prefix can be used on global values to suppress mangling.
Unnamed values are represented as an unsigned numeric value with their prefix. For example, %12, @2, %44.
Constants, which are described in the section Constants below.
LLVM requires that values start with a prefix for two reasons: Compilers don’t need to worry about name clashes with reserved words, and the set of reserved words may be expanded in the future without penalty. Additionally, unnamed identifiers allow a compiler to quickly come up with a temporary variable without having to avoid symbol table conflicts.

Reserved words in LLVM are very similar to reserved words in other languages. There are keywords for different opcodes (‘add’, ‘bitcast’, ‘ret’, etc…), for primitive type names (‘void’, ‘i32’, etc…), and others. These reserved words cannot conflict with variable names, because none of them start with a prefix character ('%' or '@').

Here is an example of LLVM code to multiply the integer variable ‘%X’ by 8:

The easy way:

%result = mul i32 %X, 8
After strength reduction:

%result = shl i32 %X, 3
And the hard way:

%0 = add i32 %X, %X           ; yields i32:%0
%1 = add i32 %0, %0           ; yields i32:%1
%result = add i32 %1, %1
This last way of multiplying %X by 8 illustrates several important lexical features of LLVM:

Comments are delimited with a ‘;’ and go until the end of line.
Unnamed temporaries are created when the result of a computation is not assigned to a named value.
Unnamed temporaries are numbered sequentially (using a per-function incrementing counter, starting with 0). Note that basic blocks and unnamed function parameters are included in this numbering. For example, if the entry basic block is not given a label name and all function parameters are named, then it will get number 0.
It also shows a convention that we follow in this document. When demonstrating instructions, we will follow an instruction with a comment that defines the type and name of value produced.

High Level Structure
Module Structure
LLVM programs are composed of Module’s, each of which is a translation unit of the input programs. Each module consists of functions, global variables, and symbol table entries. Modules may be combined together with the LLVM linker, which merges function (and global variable) definitions, resolves forward declarations, and merges symbol table entries. Here is an example of the “hello world” module:

; Declare the string constant as a global constant.
@.str = private unnamed_addr constant [13 x i8] c"hello world\0A\00"

; External declaration of the puts function
declare i32 @puts(i8* nocapture) nounwind

; Definition of main function
define i32 @main() {   ; i32()*
  ; Convert [13 x i8]* to i8*...
  %cast210 = getelementptr [13 x i8], [13 x i8]* @.str, i64 0, i64 0

  ; Call puts function to write out the string to stdout.
  call i32 @puts(i8* %cast210)
  ret i32 0
}

; Named metadata
!0 = !{i32 42, null, !"string"}
!foo = !{!0}
This example is made up of a global variable named “.str”, an external declaration of the “puts” function, a function definition for “main” and named metadata “foo”.

In general, a module is made up of a list of global values (where both functions and global variables are global values). Global values are represented by a pointer to a memory location (in this case, a pointer to an array of char, and a pointer to a function), and have one of the following linkage types.


# The Strict Intermediate Language
The **strict-intermediate-language** (from now on called *SIR*) is a low level but platform independent intermediate representation,
generated by the StrictCompiler. It can be compiled to numerous other IR's such as the llvm-ir, MSIL and
jvm-bytecode or into high-level languages like JavaScript and C++. *SIR* is type-aware but does not impose
any TypeSystem, since most target languages have their own well defined TypeSystems that often differ. The
types are rather hints that allow for linting, minor optimizations and the like. *SIR*'s structure is similar to
that of *java class files*, but instead of being stack-oriented, it uses the *SSA* form.

## Format

Structure of an encoded *SIR* file:
```
0xBADEAFFE
SYMBOL_TABLE_BEGIN
    SYMBOL_BEGIN
        $index 
        $value 
    SYMBOL_END
SYMBOL_TABLE_END
MODULE_BEGIN
    METADATA_TABLE_BEGIN
    
    METADATA_TABLE_END
    DECLARATION_SET_BEGIN
    
    DECLARATION_SET_END
MODULE_END
```

Each translation-unit is called a *module*. Modules have a SymbolTable, containing a set of symbols that are referenced
in the IR, a MetadataTable and DeclarationSet.

| Name | C-Type | Bytes |
|------|--------|------|
| SymbolReference | uint16_t | 2 |
| Identifier | string | 2 + Length |

#### Notation

Structures are represented using a pseudo-language that looks as follows:

```
Example {
    1: TypeName FieldName0
    2: TypeName FieldName1 (Optional Comment)
    N: TypeName FieldNameN
}
```

The first number is the label, it just tells the order in which the fields have to be read.

#### Keywords

Keywords in encoded *SIR* are different than keywords in *Strict* or every other high-level language,
because they are a byte rather than a character-sequence. When the reader encounters on of the keyword
bytes, it knows what structure to decode.

| Mask | Name | Description |
|------|-------|-----------|
| 0x1 | MODULE_BEGIN | Starts a module declaration | 
| 0x2 | MODULE_END | Ends a module declaration |
| 0x3 | SYMBOL_TABLE_BEGIN | Starts the symbol table |
| 0x4 | SYMBOL_TABLE_END | Ends the symbol table |
| 0x5 | SYMBOL_BEGIN  | Stats a symbol table entry |
| 0x6 | SYMBOL_END | Ends a symbol table entry |
| 0x7 | METADATA_TABLE_BEGIN | Begins the ModuleMetadataTable |
| 0x8 | METADATA_TABLE_END | Ends the ModuleMetadataTable |

### Method Declaration
#### Method Argument
```cpp
MethodArgument {
    1: SymbolReference Name
    2: SymbolReference Type (-1 = Unknown)
}
```
#### Modifiers
*SIR* contains some modifiers, which annotate so called *Targets* to give them some
extra semantics. Possible targets are `Field`, `Type`, `Method` and `All`, where `All` 
includes all prior targets.

| Mask | Name | Target | Description |
|------|------|--------|-------------|
| 0x1 | VOLATILE | Field | Volatile memory access |
| 0x2 | SYNTHETIC | All | Generated by Compiler |
| 0x3 | UNSAFE | Method | Contains Unsafe Operations |
| 0x4 | NATIVE | Method | Contains direct native calls |
| 0x5 | CLOSED | Method, Type | Can not be overridden or inherited |
| 0x6 | FINAL | Field | Unmodifiable after initialization | 
| 0x7 | VIRTUAL | Method | Method is virtual (or abstract) |