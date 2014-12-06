% title = "Signaling Type of Denial via Delegation Signer Records"
% abbrev = "DS Denial Type Signalling"
% date = "20131201T00:00:00-00:00"
% 
% [[author]]
% initials = "R."
% surname = "Gieben"
% fullnanm = "R. (Miek) Gieben"
% organization = "Google"
% address.email = "miek@google.com"
% area = "Network"
% keyword = ["DNSSEC"]

AB> This document defines a transition mechanism for using new hash algorithms 
AB> when providing hashed authenticated denial of existence in a zone. The transition mechanism
AB> defines a new digest type for Delegation Signer (DS) Resource
AB> Records that points to extra data embedded in the digest to 
AB> include the type of authenticated denial used in the zone.

{mainmatter}

# Introduction

The DS Resource Record [@RFC3658,i]
is published in parent zones to distribute a cryptographic digest of one key in a child's
DNSKEY RRset. With the DS published, a zone sets expectations for a validator. In
particular a digest of the DNSKEY, the algorithm used for signature of the
DNSKEY and the type of authenticated denial of existence used.

When NSEC3 [@RFC5155,n] was ....

> This transition method is best described as a hack.

In this document, the key words "MUST", "MUST NOT", "REQUIRED",
"SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY",
and "OPTIONAL" are to be interpreted as described in [@RFC2119,n].

# DS Record Field Values

Indicating the type of denial of existence in use at the child zone is done by
prefixing the digest in the DS record with two octets defining
the Denial Type:

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

This document does *not* change the presentation format of DS records.

##  DS Record with Denial Type Wire Format

The resulting on-the-wire format for the resulting DS record will be as follows:

{#fig:wire}
F> # The on-the-wire format for the DS. The length of the digest is specified in the respective RFCs defining the digest type.
F>
F>                          1 1 1 1 1 1 1 1 1 1 2 2 2 2 2 2 2 2 2 2 3 3
F>      0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
F>     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
F>     |           Key Tag             |  Algorithm    |  DigestType   |
F>     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
F>     |          Denial Type          |                               /
F>     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+          Digest               /
F>     /                                                               /
F>     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|

The Denial Type is a 16 bit unsigned integer value stored in network order.

##  Example DS Record Using SHA-256 and Denial Type TBD

The following is an example DNSKEY, and a matching DS record that 
includes denial type TBD that refences that NSEC3/SHA1 is in use in 
the child zone. This
DNSKEY record comes from the example DNSKEY/DS records found in
section 5.4 of [@RFC4034].

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

The Digest Types to be used for supporting Denial Type information within
DS records has been assigned by IANA.

At the time of this writing, the current digest types assigned for
use in DS records are as follows:

VALUE  |  Digest Type     |   Status
-------|------------------|-------------
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

T> # A table caption.
T>
T>  VALUE   |  Denial Type
T>  --------|-------------------
T>     0    |  Reserved
T>     1    |  NSEC
T>     2    |  NSEC3 w/ SHA-1
T>     3    |  NSEC3 w/ SHA-256
T>     4    |  NSEC3 w/ SHA-384
T>  5-65535 |  Unassigned

# Acknowledgements

...

{backmatter}

# Other Options

## Querying for NSEC3PARAM

A different approach could have been to use the NSEC3PARAM record at
the zone apex to determine the hash algorithm, as well as other
parameters. This would however result in mandatory extra queries
during validation, and it would be a change in the resolution
process, as the NSEC3PARAM is currently not used by validating
resolvers.

## Algorithm Aliasing

This is a good, or maybe the best way to deal with this transition, but 
because the algorithm namespace is only 8 bits and each aliases need to
alias all previous aliases...

# Changelog

## 00

* Initial release.
