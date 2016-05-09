% Title = "Scenic Routing for IPv6"
% abbrev = "Scenic Routing for IPv6"
% category = "info"
% docName = "rfc-7511"
% ipr= "trust200902"
% area = "Internet"
% workgroup = "Network Working Group"
%
% date = 2015-04-01T00:00:00Z
%
% [[author]]
% initials="M."
% surname="Wilhelm"
% fullname="Maximilian Wilhelm"
%  [author.address]
%  email = "max@rfc2324.org"
%  phone = "+49 176 62 05 94 27"
%   [author.address.postal]
%   city = "Paderborn, NRW"
%   country = "Germany"

.# Abstract

This document specifies a new routing scheme for the current version
of the Internet Protocol version 6 (IPv6) in the spirit of "Green
IT", whereby packets will be routed to get as much fresh-air time as
possible.

{mainmatter}

#  Introduction

In times of Green IT, a lot of effort is put into reducing the energy
consumption of routers, switches, servers, hosts, etc., to preserve
our environment.  This document looks at Green IT from a different
angle and focuses on network packets being routed and switched around
the world.

Most likely, no one ever thought about the millions of packets being
disassembled into bits every second and forced through copper wires
or being shot through dark fiber lines by powerful lasers at
continuously increasing speeds.  Although RFC 5841 [@!RFC5841] provided
some thoughts about Packet Moods and began to represent them as a TCP
option, this doesn't help the packets escape their torturous routine.

This document defines another way to deal with Green IT for traffic
and network engineers and will hopefully aid the wellbeing of a
myriad of network packets around the world.  It proposes Scenic
Routing, which incorporates the green-ness of a network path into the
routing decision.  A routing engine implementing Scenic Routing
should therefore choose paths based on Avian IP Carriers [@?RFC1149]
and/or wireless technologies so the packets will get out of the
miles/kilometers of dark fibers that are in the ground and get as
much fresh-air time and sunlight as possible.

As of the widely known acceptance of the current version of the
Internet Protocol (IPv6), this document only focuses on version 6 and
ignores communication still based on Vintage IP [@?RFC0791].

##  Conventions and Terminology

The key words "**MUST**", "**MUST NOT**", "**REQUIRED**", "**SHALL**", "**SHALL NOT**",
"**SHOULD**", "**SHOULD NOT**", "**RECOMMENDED**", "**MAY**", and "**OPTIONAL**" in this
document are to be interpreted as described in RFC 2119 [@!RFC2119].

Additionally, the key words "**MIGHT**", "**COULD**", "**MAY WISH TO**", "**WOULD
PROBABLY**", "**SHOULD CONSIDER**", and "**MUST (BUT WE KNOW YOU WON'T)**" in
this document are to interpreted as described in RFC 6919 [@!RFC6919].

#  Scenic Routing

Scenic Routing can be enabled with a new option for IPv6 datagrams.

##  Scenic Routing Option (SRO)

The Scenic Routing Option (SRO) is placed in the IPv6 Hop-by-Hop
Options Header that must be examined by every node along a packet's
delivery path [@!RFC2460].

The SRO can be included in any IPv6 datagram, but multiple SROs **MUST
NOT** be present in the same IPv6 datagram.  The SRO has no alignment
requirement.

If the SRO is set for a packet, every node en route from the packet
source to the packet's final destination **MUST** preserve the option.

The following Hop-by-Hop Option is proposed according to the
specification in Section 4.2 of RFC 2460 [@!RFC2460].

{#fig-scenic-routing-option-layout}
~~~
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
                                +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
                                |  Option Type  | Option Length |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|   SRO Param   |                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
~~~
Figure: Scenic Routing Option Layout

Option Type
: <br/>8-bit identifier of the type of option.  The option identifier
    0x0A (On Air) is proposed for Scenic Routing.

    {#fig-option-type}
    ~~~~~
    HEX         act  chg  rest
    ---         ---  ---  -----
    0A           00   0   01010     Scenic Routing
    ~~~~~
    Figure: Scenic Routing Option Type

    The highest-order two bits are set to 00 so any node not
    implementing Scenic Routing will skip over this option and
    continue processing the header.  The third-highest-order bit
    indicates that the SRO does not change en route to the packet's
    final destination.

Option Length
: <br/>8-bit unsigned integer.  The length of the option in octets
    (excluding the Option Type and Option Length fields).  The value
    **MUST** be greater than 0.

SRO Param
: <br/>8-bit identifier indicating Scenic Routing parameters encoded as a bit string.

    {#fig-bit-string-layout}
    ~~~~~
    +-+-+-+-+-+-+-+-+
    | SR A W AA X Y |
    +-+-+-+-+-+-+-+-+
    ~~~~~
    Figure: SRO Param Bit String Layout

    The highest-order two bits (SR) define the urgency of Scenic
    Routing:

    {style="empty"}
    - 00 - Scenic Routing **MUST NOT** be used for this packet.
    - 01 - Scenic Routing **MIGHT** be used for this packet.
    - 10 - Scenic Routing **SHOULD** be used for this packet.
    - 11 - Scenic Routing **MUST** be used for this packet.

    The following BIT (A) defines if Avian IP Carriers should be used:

    {style="empty"}
    - 0 - Don't use Avian IP Carrier links (maybe the packet is
      afraid of pigeons).
    - 1 - Avian IP Carrier links may be used.

    The following BIT (W) defines if wireless links should be used:

    {style="empty"}
    - 0 - Don't use wireless links (maybe the packet is afraid of
      radiation).
    - 1 - Wireless links may be used.

    The following two bits (AA) define the affinity for link types:

    {style="empty"}
    - 00 - No affinity.
    - 01 - Avian IP Carriers **SHOULD** be preferred.
    - 10 - Wireless links **SHOULD** be preferred.
    - 11 - RESERVED

    The lowest-order two bits (XY) are currently unused and reserved
    for future use.

# Implications

## Routing Implications

If Scenic Routing is requested for a packet, the path with the known
longest Avian IP Carrier and/or wireless portion **MUST** be used.

Backbone operators who desire to be fully compliant with Scenic
Routing **MAY WISH TO** -- well, they **SHOULD** -- have separate MPLS paths
ready that provide the most fresh-air time for a given path and are
to be used when Scenic Routing is requested by a packet.  If such a
path exists, the path MUST be used in favor of any other path, even
if another path is considered cheaper according to the path costs
used regularly, without taking Scenic Routing into account.

## Implications for Hosts

Host systems implementing this option of receiving packets with
Scenic Routing requested **MUST** honor this request and **MUST** activate
Scenic Routing for any packets sent back to the originating host for
the current connection.

If Scenic Routing is requested for connections of local origin, the
host MUST obey the request and route the packet(s) over a wireless
link or use Avian IP Carriers (if available and as requested within
the SRO Params).

System administrators **MIGHT** want to configure sensible default
parameters for Scenic Routing, when Scenic Routing has been widely
adopted by operating systems.  System administrators **SHOULD** deploy
Scenic Routing information where applicable.

##  Proxy Servers

If a host is running a proxy server or any other packet-relaying
application, an application implementing Scenic Routing **MUST** set the
same SRO Params on the outgoing packet as seen on the incoming
packet.

Developers **SHOULD CONSIDER** Scenic Routing when designing and
implementing any network service.

#  Security Considerations

The security considerations of RFC 6214 [@!RFC6214] apply for links
provided by Avian IP Carriers.

General security considerations of wireless communication apply for
links using wireless technologies.

As the user is able to influence where flows and packets are being
routed within the network, this **MIGHT** influence traffic-engineering
considerations and network operators **MAY WISH TO** take this into
account before enabling Scenic Routing on their devices.

#  IANA Considerations

This document defines a new IPv6 Hop-by-Hop Option, the Scenic
Routing Option, described in (#scenic-routing-option-sro).
If this work is standardized, IANA is requested to assign a value from the "Destination Options and
Hop-by-Hop Options" registry for the purpose of Scenic Routing.

There are no IANA actions requested at this time.

#  Related Work

As Scenic Routing is heavily dependent on network paths and routing
information, it might be worth looking at designing extensions for
popular routing protocols like BGP or OSPF to leverage the full
potential of Scenic Routing in large networks built upon lots of
wireless links and/or Avian IP Carriers.  When incorporating
information about links compatible with Scenic Routing, the routing
algorithms could easily calculate the optimal paths providing the
most fresh-air time for a packet for any given destination.

This would even allow preference for wireless paths going alongside
popular or culturally important places.  This way, the packets don't
only avoid the dark fibers, but they get to see the world outside of
the Internet and are exposed to different cultures around the globe,
which may help build an understanding of cultural differences and
promote acceptance of these differences.

{backmatter}

# Acknowledgements

The author wishes to thank all those poor friends who were kindly
forced to read this document and that provided some nifty comments.
