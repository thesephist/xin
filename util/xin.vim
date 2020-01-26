" place this in the init path (.vimrc)
au BufNewFile,BufRead *.xin set filetype=xin

" place this in $HOME/.vim/syntax/ink.vim
if exists("b:current_syntax")
    finish
endif

" xin syntax definition for vi/vim
syntax sync fromstart

" lisp-style indentation
set lisp

" booleans
syntax keyword xinBoolean true false
highlight link xinBoolean Boolean

" numbers should be consumed first by identifiers, so comes before
syntax match xinNumber "\v-?\d+[.\d+]?"
highlight link xinNumber Number

" builtin functions
syntax keyword xinKeyword if contained
syntax keyword xinKeyword do contained
syntax keyword xinKeyword co contained
highlight link xinKeyword Keyword

" functions
syntax match xinFunctionForm "\v\(\s*[A-Za-z0-9\-?!+*/:><=%&|]*" contains=xinFunctionName,xinKeyword
syntax match xinFunctionName "\v[A-Za-z0-9\-?!+*/:><=%&|]*" contained
highlight link xinFunctionName Function

syntax match xinDefinitionForm "\v\(\s*\:\s*(\(\s*)?[A-Za-z0-9\-?!+*/:><=%&|]*" contains=xinDefinition
syntax match xinDefinition "\v[A-Za-z0-9\-?!+*/:><=%&|]*" contained
highlight link xinDefinition Type

" strings
syntax region xinString start=/\v'/ skip=/\v(\\.|\r|\n)/ end=/\v'/
highlight link xinString String

" comment
" -- block
" -- line-ending comment
syntax match xinComment "\v;.*" contains=xinTodo
highlight link xinComment Comment
" -- shebang, highlighted as comment
syntax match xinShebangComment "\v^#!.*"
highlight link xinShebangComment Comment
" -- TODO in comments
syntax match xinTodo "\v(TODO\(.*\)|TODO)" contained
syntax keyword xinTodo XXX contained
highlight link xinTodo Todo

syntax region inkForm start="(" end=")" transparent fold
set foldmethod=syntax
set foldlevel=20

let b:current_syntax = "xin"
