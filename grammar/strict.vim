" Vim syntax file
" Language:	Strict
" Maintainer:	Merlin Osayimwen <merlinosayimwen@gmail.com>
" Last Change:	2019 Aug 28
if exists("b:current_syntax")
  finish
endif

let s:cpo_save = &cpo
set cpo&vim

if exists("strict_no_doctest_highlight")
  let strict_no_doctest_code_highlight = 1
endif

if exists("strict_highlight_all")
  if exists("strict_no_builtin_highlight")
    unlet strict_no_builtin_highlight
  endif
  if exists("strict_no_doctest_code_highlight")
    unlet strict_no_doctest_code_highlight
  endif
  if exists("strict_no_doctest_highlight")
    unlet strict_no_doctest_highlight
  endif
  if exists("strict_no_exception_highlight")
    unlet strict_no_exception_highlight
  endif
  if exists("strict_no_number_highlight")
    unlet strict_no_number_highlight
  endif
  let strict_space_error_highlight = 1
endif

syn keyword strictStatement	false true empty
syn keyword strictStatement	assert break continue
syn keyword strictStatement	lambda test  return  yield
syn keyword strictStatement	method nextgroup=strictFunction skipwhite
syn keyword strictConditional	else if do
syn keyword strictRepeat for
syn keyword strictOperator	and or is isnt
syn keyword strictException	throw try catch
syn keyword strictInclude	import as

syn match   strictFunction	"\h\w*" display contained

syn match   strictComment	"//.*$" contains=strictTodo,@Spell
syn keyword strictTodo		FIXME NOTE NOTES TODO XXX contained
syn region  strictString matchgroup=strictQuotes
      \ start=+[uU]\=\z(['"]\)+ end="\z1" skip="\\\\\|\\\z1"
      \ contains=strictEscape,@Spell
syn region  strictString matchgroup=strictTripleQuotes
      \ start=+[uU]\=\z('''\|"""\)+ end="\z1" keepend
      \ contains=strictEscape,strictSpaceError,strictDoctest,@Spell
syn region  strictRawString matchgroup=strictQuotes
      \ start=+[uU]\=[rR]\z(['"]\)+ end="\z1" skip="\\\\\|\\\z1"
      \ contains=@Spell
syn region  strictRawString matchgroup=strictTripleQuotes
      \ start=+[uU]\=[rR]\z('''\|"""\)+ end="\z1" keepend
      \ contains=strictSpaceError,strictDoctest,@Spell

syn match   strictEscape	+\\[abfnrtv'"\\]+ contained
syn match   strictEscape	"\\\o\{1,3}" contained
syn match   strictEscape	"\\x\x\{2}" contained
syn match   strictEscape	"\%(\\u\x\{4}\|\\U\x\{8}\)" contained
syn match   strictEscape	"\\N{\a\+\%(\s\a\+\)*}" contained
syn match   strictEscape	"\\$"

if !exists("strict_no_number_highlight")
  syn match  strictNumber	"\<0[oO]\=\o\+[Ll]\=\>"
  syn match  strictNumber	"\<0[xX]\x\+[Ll]\=\>"
  syn match  strictNumber	"\<0[bB][01]\+[Ll]\=\>"
  syn match  strictNumber	"\<\%([1-9]\d*\|0\)[Ll]\=\>"
  syn match  strictNumber	"\<\d\+[jJ]\>"
  syn match  strictNumber	"\<\d\+[eE][+-]\=\d\+[jJ]\=\>"
  syn match  strictNumber
	\ "\<\d\+\.\%([eE][+-]\=\d\+\)\=[jJ]\=\%(\W\|$\)\@="
  syn match  strictNumber
	\ "\%(^\|\W\)\zs\d*\.\d\+\%([eE][+-]\=\d\+\)\=[jJ]\=\>"
endif

if !exists("strict_no_builtin_highlight")
  syn keyword strictBuiltin	false true empty
  syn keyword strictBuiltin	Log LogFormatted
  syn keyword strictBuiltin	String List Sequence Enumerator
  syn keyword strictBuiltin	int float bool
endif

if exists("strict_space_error_highlight")
  syn match   strictSpaceError	display excludenl "\s\+$"
  syn match   strictSpaceError	display " \+\t"
  syn match   strictSpaceError	display "\t\+ "
endif

if !exists("strict_no_doctest_highlight")
  if !exists("strict_no_doctest_code_highlight")
    syn region strictDoctest
	  \ start="^\s*>>>\s" end="^\s*$"
	  \ contained contains=ALLBUT,strictDoctest,strictFunction,@Spell
    syn region strictDoctestValue
	  \ start=+^\s*\%(>>>\s\|\.\.\.\s\|"""\|'''\)\@!\S\++ end="$"
	  \ contained
  else
    syn region strictDoctest
	  \ start="^\s*>>>" end="^\s*$"
	  \ contained contains=@NoSpell
  endif
endif

syn sync match strictSync grouphere NONE "^\%(method\)\s\+\h\w*\s*("

hi def link strictStatement		Statement
hi def link strictConditional		Conditional
hi def link strictRepeat		Repeat
hi def link strictOperator		Operator
hi def link strictInclude		Include
hi def link strictAsync			Statement
hi def link strictFunction		Function
hi def link strictComment		Comment
hi def link strictTodo			Todo
hi def link strictString		String
hi def link strictRawString		String
hi def link strictQuotes		String
hi def link strictTripleQuotes		strictQuotes
hi def link strictEscape		Special
if !exists("strict_no_number_highlight")
  hi def link strictNumber		Number
endif
if !exists("strict_no_builtin_highlight")
  hi def link strictBuiltin		Function
endif
if !exists("strict_no_exception_highlight")
  hi def link strictExceptions		Structure
endif
if exists("strict_space_error_highlight")
  hi def link strictSpaceError		Error
endif
if !exists("strict_no_doctest_highlight")
  hi def link strictDoctest		Special
  hi def link strictDoctestValue	Define
endif

let b:current_syntax = "strict"

let &cpo = s:cpo_save
unlet s:cpo_save

