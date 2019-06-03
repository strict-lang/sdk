# The Strict Programming Language

## Translation Units

### Names


Every *user defined type* and method is in a separate file. The files name encodes
information of its content. This makes filenames a crucial part of Strict source code. 
Strict files do not have a common filename extension like most programming languages.
Extensions tell whether the content of a file is a method or type.

Following table shows the valid filename extensions and their corresponding content.

| Extension | Content | Filename Example           |
|-----------|---------|----------------------------|
| .method   | Method  | floor(number)number.method |
| .type     | Type    | account.type      		     |

### Content

Content of the `person.type` File

```strict
name
age
```

