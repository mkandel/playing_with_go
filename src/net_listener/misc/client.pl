#!/usr/local/bin/perl -w
# $Id:$
# $HeadURL:$
use strict;
use warnings;

use Carp;
use Getopt::Long;
Getopt::Long::Configure qw/bundling no_ignore_case/;
use Data::Dumper;
# Some Data::Dumper settings:
local $Data::Dumper::Useqq  = 1;
local $Data::Dumper::Indent = 3;
use Pod::Usage;
use IO::Socket::INET;

local $| = 1;

my $debug   = 0;
my $dryrun  = 0;
my $file;

GetOptions(
    "help|h"         => sub { pod2usage( 1 ); },
    "debug|d"        => \$debug,
    "dryrun|n"       => \$dryrun,
    "file|f=s"       => \$file,
);

my $prog = $0;
$prog =~ s/^.*\///;
 
# create a connecting socket
my $socket = new IO::Socket::INET (
    PeerHost => '127.0.0.1',
    #PeerHost => '10.11.9.181',
    PeerPort => '7777',
    Proto => 'tcp',
);
die "cannot connect to the server $!\n" unless $socket;
print "connected to the server\n";

my $req;

if ( $file && -e $file ){
    open my $IN, '<', $file || die "Error opening '$file' for read: $!\n";
    my @data = <$IN>;
    $req = join '\n', @data;
    close $IN || die "Error closing '$file' after read: $!\n";
}

$req //= "hello world (from Perl)";
my $size = $socket->send($req);
print "sent data of length $size\n";
 
# notify server that request has been sent
#shutdown($socket, 1);
 
# receive a response of up to 1024 characters from server
#my $response = "";
#$socket->recv($response, 1024);
#print "received response: $response\n";
 
$socket->close();

__END__

