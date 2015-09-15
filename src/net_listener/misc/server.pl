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
use threads;

local $| = 1;

my $debug   = 0;
my $dryrun  = 0;
my $port = 7777;

GetOptions(
    "help|h"         => sub { pod2usage( 1 ); },
    "debug|d"        => \$debug,
    "dryrun|n"       => \$dryrun,
    "port|p=n"       => \$port,
);

my $prog = $0;
$prog =~ s/^.*\///;
 
my $sep = "=====================================================================\n";
# creating a listening socket
my $socket = new IO::Socket::INET (
    LocalHost => '127.0.0.1',
    #LocalHost => '10.11.9.181',
    LocalPort => $port,
    Proto => 'tcp',
    Listen => 5,
    Reuse => 1
);
autoflush $socket;

croak "cannot create socket $!\n" unless $socket;
print $sep;
print "server waiting for client connection on port '$port'\n";
print $sep;
 
while(1)
{
    ## waiting for a new client connection
    my $client_socket = $socket->accept();
    ## my $thr1 = threads->create(\&sub1, 'Param 1', 'Param 2', $Param3);
    ## handle_client( $client_socket );
    my $thr = threads->create( \&handle_client, $client_socket );
    $thr->detach();
}
 
$socket->close();

sub handle_client {
    my $client_socket = shift || croak "Socket needed!!\n";
    ## get information about a newly connected client
    my $client_address = $client_socket->peerhost();
    my $client_port = $client_socket->peerport();
    print "connection from $client_address:$client_port\n";
 
    ## read up to 1024 characters from the connected client
    my $data = "";
    print "received data: ";
    while( $data = <$client_socket> ){
        print "$data";
    }

    print "\nProcessed by threadID: '" . threads->tid() . "'\n";
    print $sep;
 
    ## write response data to the connected client
    $data = "ok\n";
    $client_socket->send($data);
 
    ## notify client that response has been sent
    shutdown($client_socket, 1);
}

__END__

