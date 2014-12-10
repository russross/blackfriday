" Vim syntax file
" Language:	mmark
" Maintainer:	Miek Gieben <miek@miek.nl>
" Last Change:	Sun, 28 Oct 2001 21:22:24 +0100
" Filenames:	*.md (markdown)

" TODO(finish)

" Quit when a syntax file was already loaded
if exists("b:current_syntax")
    finish
endif

runtime syntax/pandoc.vim

syn match   mmarkMatter      /{frontmatter}/ contained
syn match   mmarkMatter      /{mainmatter}/ contained
syn match   mmarkMatter      /{backmatter}/ contained

hi link mmarkMatter Special

" vim: ts=8
