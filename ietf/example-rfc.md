% title = "Signaling Type of Denial via Delegation Signer Records"
% abbrev = "DS Denial Type Signalling"
% docName = "example-00"
% ipr = "trust200902"
% category = "info"
%
% date = 2014-12-01T00:00:00Z
% area = "Network"
% keyword = ["DNSSEC"]
%
% [[author]]
% initials = "R."
% surname = "Gieben"
% fullname = "R. (Miek) Gieben"
% organization = "Google"
%   [author.address]
%   email = "miek@google.com"

A> This will become a paragraph in the abstract.

{mainmatter}

# Introduction

The DS Resource Record [@RFC3658,i]...
particular a digest of the **DNSKEY**, *the* algorithm used for signature of the

When NSEC3 [@RFC5155,n 5.5] was ....

> This transition method is best described as a hack.

In this document, the key words "MUST", "MUST NOT", "REQUIRED",
"SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY",
and "OPTIONAL" are to be interpreted as described in [@RFC2119,n].

# DS Record Field Values

AS> When typesetting something in an aide
AS> you get an aside.

<!-- Miek: are you sure you want to include this stuff? -->

Indicating the type of denial of existence in use at the child zone is done by

Denial Type:
:   An extra 16 bit integer value (see [](#iana-considerations)) encoded in the DS' digest
    that indicates the denial of existence in use in the (child) zone.

Digest:
:   The digest value is calculated by using the following
    formula ("|" denotes concatenation, HASH denotes that
    hash algorithm in use).

            digest = Denial Type | HASH(DNSKEY owner name | DNSKEY RDATA)

    where DNSKEY RDATA is defined by [@RFC4034,n] as:

            DNSKEY RDATA = Flags | Protocol | Algorithm | Public Key

    The Key Tag field and Algorithm fields remain unchanged by this
    document and are specified in the [@RFC4034] specification.

Denial Type:
:   An extra 16 bit integer value (see [](#iana-considerations)) encoded in the DS' digest
    that indicates the denial of existence in use in the (child) zone.

This document does *not* change the presentation format of DS records.

##  DS Record with Denial Type Wire Format

The resulting on-the-wire format for the resulting DS record will be as follows:

{#fig:wire}
Figure: The on-the-wire format for the DS. The length of the digest is specified in the respective RFCs defining the digest type.
                         1 1 1 1 1 1 1 1 1 1 2 2 2 2 2 2 2 2 2 2 3 3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |           Key Tag             |  Algorithm    |  DigestType   |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |          Denial Type          |                               /
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+          Digest               /
    /                                                               /
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|

The Denial Type is a 16 bit unsigned integer value stored in network order.

##  Example DS Record Using SHA-256 and Denial Type TBD

Fenced code block
Figure: This is a fenced code block
``` go
println(hallo)
```

DNSKEY record comes from the example DNSKEY/DS records found in section 5.4 of [@RFC4034].

As you can use over at this URL <http://www.miek.nl>.

The DNSKEY record:

    dskey.example.com. 86400 IN DNSKEY 256 3 5 ( AQOeiiR0GOMYkDshWoSKz9Xz
                                              fwJr1AYtsmx3TGkJaNXVbfi/
                                              2pHm822aJ5iI9BMzNXxeYCmZ
                                              DRD99WYwYqUSdjMmmAphXdvx
                                              egXd/M5+X7OrzKBaMbCVdFLU
                                              Uh6DhweJBjEVv5f2wwjM9Xzc
                                              nOf+EPbtG9DMBmADjFDc2w/r
                                              ljwvFw==) ;  key id = 60485

# IANA Considerations

The following action for IANA are required by this document:

At the time of this writing, the current digest types assigned for
use in DS records are as follows:

VALUE  |  Digest Type     |   Status
------:|:-----------------|-------------
 0     | Reserved         |      -
 1     | SHA-1            |   MANDATORY
 2     | SHA-256          |   MANDATORY
 3     | GOST R 34.11-94  |   OPTIONAL
 4     | SHA-384          |   OPTIONAL
 TBD   | DT-SHA-256       |   OPTIONAL
TDB-255| Unassigned       |      -

All future assigned Digest Types MUST assume that there is a Denial Type incorporated in the Digest.

This document creates a new IANA registry for Denial Types.  This
registry is named "DNSSEC DENIAL TYPES".  The initial contents of this
registry are:

Table: As shown here.
VALUE   |  Denial Type
--------|-------------------
   0    |  Reserved
   1    |  NSEC
   2    |  NSEC3 w/ SHA-1
   3    |  NSEC3 w/ SHA-256
   4    |  NSEC3 w/ SHA-384
5-65535 |  Unassigned

<!--  Miek Gieben: This is a comment -->

# Acknowledgements

The people in the following list:

* ...
* And ...
    1. another list
    2. another list
* And ...

This is an citation that shows up in the references, but not in the document: [-@RFC1033].
This needs to happen *BEFORE* the `{backmatter`} otherwise the references are outputted, but
you use the references.

{backmatter}

# Other Options

## Images

{type="ascii-art"}
![alt text](https://github.com/adam-p/markdown-here/raw/master/src/common/images/icon48.png "Title")

{type="ascii-art"}
![alt text](https://github.com/adam-p/markdown-here/raw/master/src/common/images/icon48.png "Title2")

## Algorithm Aliasing

Now we can use the reference to [@RFC1033] here.

This is a good, or maybe the best way to deal with this transition, but
because the algorithm namespace is only 8 bits and each aliases need to
alias all previous aliases...

# Changelog

## 00

* Initial release.
