#!/usr/bin/perl

use Switch;

use XML::Parser;
use Data::Dumper;

sub ext { return 'perl-unhandled-extern-ref' }

sub start {
    shift;
    $e = shift;
    switch ($e) {
        case "title" {
            for ($i = 0; $i < scalar @_ - 1; $i++) {
                print $_[$i] . " = " . '"' . $_[$i+1] . '"' . "\n";
            }
            print "title = \"";
        }
        case "rfc" {
            for ($i = 0; $i < scalar @_ - 1; $i++) {
                print $_[$i] . " = " . '"' . $_[$i+1] . '"' . "\n";
                $i++;
            }
        }
        case "author" {
            print "\n[[author]]\n";
            for ($i = 0; $i < scalar @_ - 1; $i++) {
                print $_[$i] . " = " . '"' . $_[$i+1] . '"' . "\n";
                $i++;
            }
        }
        case "address"      { print "[author.address]\n" }
        case "postal"       { print "[author.address.postal]\n" }
        # todo street code country
        # todo date
        case "keyword"      { print "keyword = [" }
        case "area"         { print "area = \"" }
        case "workgroup"    { print "workgroup = \"" }
        case "organization" { print "organization = \"" }
        case "email"        { print "email = \"" }
        case "uri"          { print "uri = \"" }
        case "phone"        { print "phone = \"" }
    }
}

sub end {
    shift;
    $e = shift;
    switch ($e) {
        # rfc ipr="trust200902" category="exp" docName="draft-gieben-nsec4-02">
        case "rfc"          { }
        case "address"      { }
        case "keyword"      { print "]\n"; }
        case "author"       { print "\n"; }
        case "email"        { print "\"\n"; }
        case "uri"          { print "\"\n"; }
        case "phone"        { print "\"\n"; }
        case "title"        { print "\"\n"; }
        case "area"         { print "\"\n"; }
        case "organization" { print "\"\n"; }
        case "workgroup"    { print "\"\n"; }
        case "front"        { exit; }
    }
}

sub char {
    $p = shift;
    switch ($p->current_element) {
        case "keyword"        { print join ' ', map qq("$_"), @_; }
        case "title"          { print $_[0] ; }
        case "area"           { print $_[0] ; }
        case "workgroup"      { print $_[0] ; }
        case "organization"   { print $_[0] ; }
        case "email"          { print $_[0] ; }
        case "uri"            { print $_[0] ; }
        case "phone"          { print $_[0] ; }
    }
}

my $xmlfile = shift @ARGV;
my $parser = XML::Parser->new( Style => 'Tree', ErrorContext => 2, NoExpand => 1 );
$parser->setHandlers(ExternEnt => \&ext, Start => \&start, End => \&end, Char => \&char);
eval { $dom = $parser->parsefile( $xmlfile ); };
