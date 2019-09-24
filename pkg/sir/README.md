# Strict Intermediate Representation Specification

- [Introduction](#Introduction)
- [Structure](#Structure)
  - [Overview](#Overview)
  - [Constant Pool](#Constant Pool)
  - [Modules](#Modules)
  - [Metadata Tables](#Metadata Tables)
  - [Method Declarations](#Method Declarations)
    - [Parameter List](#Parameter List)
    - [Local Variable Table](#Local Variable Table)
- [Encoding](#Encoding) 
- [Metadata](#Metadata)
- [Operations](#operations)
  - [Description Format](#Description Format)
  - [Parameterized Operations](#Parameterized Operations)
  - [Operations](#Operations)
    - [`add<>`](#add)
    - [`sub<>`](#sub)
    - [`mul<>`](#mul)
    - [`div<>`](#div)
    - [`mod<>`](#mod)
    - [`cmp<>`](#cmp)
    - [`cmp_eq<>`](#cmp_eq)
    - [`cmp_ne<>`](#cmp_ne)
    - [`cmp_lt<>`](#cmp_lt)
    - [`cmp_le<>`](#cmp_le)
    - [`cmp_gt<>`](#cmp_gt)
    - [`cmp_ge<>`](#cmp_ge)
    - [`call<>`](#call)
    - [`cast<>`](#cast)
    - [`len`](#len)
    - [`throw`](#throw)
    
##Introduction
The Strict intermediate representation (*SIR*) is a statically typed, Single Static Assignment (*SSA*) based IR, 
that represents compiled Strict code. It is low level compared to Strict, but still aware of high level concepts.

##Structure
###Overview
###Constant Pool
###Modules
###Metadata Tables
###Method Declarations
####Parameter List
####Local Variable Table

# Encoding 

# Metadata

# Operations 

## Description Format

## Parameterized Operations

## Operation List
___
###add<>
Adds two numbers.
#### Format
`$2 = add<int>($0, $1)`
#### Parameter
Any number type
#### Description
Adds two values of the same type.
___

###sub<>
###mul<>
###mod<>
###div<>
###cmp<>
###cmp_eq<>
###cmp_ne<>
###cmp_lt<>
###cmp_le<>
###cmp_gt<>
###cmp_ge<>
###call<>
###cast<>
###len
###throw

