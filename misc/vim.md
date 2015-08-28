# Syntax Highlighting

For Vim's syntax highlighting you can use just the vim-pandoc-syntax.
I personally use the following extra highlight groups for mmark. Note
that block-tables are not highlighted corrected.

    " Mmark stuff
    au Filetype pandoc syn match mmMatter /^{mainmatter}/
    au Filetype pandoc syn match mmMatter /^{frontmatter}/
    au Filetype pandoc syn match mmMatter /^{backmatter}/
    au Filetype pandoc syn match mmPartHeader /^\s*-#.*\n/ contains=pandocEmphasis,pandocStrong,pandocNoFormatted,@Spell
    au Filetype pandoc syn match mmSpecialHeader /^\s*\.#.*\n/ contains=pandocEmphasis,pandocStrong,pandocNoFormatted,@Spell
    au Filetype pandoc syn region mmAside       start="^\s*A>.*" end="$" contains=@Spell
    au Filetype pandoc syn region mmNote        start="^\s*N>.*" end="$" contains=@Spell
    au Filetype pandoc syn region mmFigure      start="^\s*F>.*" end="$" contains=@Spell
    au Filetype pandoc syn match mmXref       /(#[[:graph:]äëïöüáéíóúàèìòùłßÄËÏÖÜÁÉÍÓÚÀÈÌÒÙŁß].*)/ contains=@NoSpell display
    au Filetype pandoc call TextEnableCodeSnip('toml', '^% *', '$', 'PreProc')
    au Filetype pandoc hi link pandocHTMLComment Comment
    au Filetype pandoc hi def link mmPartHeader SpecialKey
    au Filetype pandoc hi def link mmSpecialHeader PreProc
    au Filetype pandoc hi def link mmMatter Special
    au Filetype pandoc hi def link mmAside Operator
    au Filetype pandoc hi def link mmNote Operator
    au Filetype pandoc hi def link mmFigure Identifier
    au Filetype pandoc hi def link mmXref Identifier
    au Filetype pandoc hi def link pandocCodeBlock String

    function! TextEnableCodeSnip(filetype,start,end,textSnipHl) abort
      let ft=toupper(a:filetype)
      let group='textGroup'.ft
      if exists('b:current_syntax')
        let s:current_syntax=b:current_syntax
        " Remove current syntax definition, as some syntax files (e.g. cpp.vim)
        " do nothing if b:current_syntax is defined.
        unlet b:current_syntax
      endif
      execute 'syntax include @'.group.' syntax/'.a:filetype.'.vim'
      try
        execute 'syntax include @'.group.' after/syntax/'.a:filetype.'.vim'
      catch
      endtry
      if exists('s:current_syntax')
        let b:current_syntax=s:current_syntax
      else
        unlet b:current_syntax
      endif
      execute 'syntax region textSnip'.ft.'
      \ matchgroup='.a:textSnipHl.'
      \ start="'.a:start.'" end="'.a:end.'"
      \ contains=@'.group
    endfunction
