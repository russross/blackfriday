% Title = "Using mmark to create I-Ds and RFCs"
% abbrev = "mmark2rfc"
% category = "ifno"
% docName = "draft-gieben-mmark2rfc-00"
% ipr= "trust200902"
% date = 2014-12-10T00:00:00Z
% area = "Internet"
% workgroup = ""
% keyword = ["markdown", "xml", "mmark"]
%
% [[author]]
% initials="R."
% surname="Gieben"
% fullname="R. (Miek) Gieben"
% organization = "Google"
%   [author.address]
%   email = "miek@google.com"

A> This document describes an markdown variant called mmark [@!mmark] that can
A> be used to create RFC documents. It's aim is to make using mmark is natural
A> as possible, while providing a lot of power on how to structure and layout
A> the document.

{mainmatter}

# Introduction


{backmatter}


# Raw references?

R!> <reference anchor='mmark' target="http://github.com/miekg/mmark">
R!>     <front>
R!>         <title abbrev='mmark'>Mmark git repository</title>
R!>         <author initials='R.' surname='Gieben' fullname='R. (Miek) Gieben'>
R!>             <address>
R!>                 <email>miek@miek.nl</email></address></author>
R!>         <date year='2014' month='December' />
R!>     </front>
R!> </reference>

R?> <reference anchor='mmark' target="http://github.com/miekg/mmark">
R?>     <front>
R?>         <title abbrev='mmark'>Mmark git repository</title>
R?>         <author initials='R.' surname='Gieben' fullname='R. (Miek) Gieben'>
R?>             <address>
R?>                 <email>miek@miek.nl</email></address></author>
R?>         <date year='2014' month='December' />
R?>     </front>
R?> </reference>
