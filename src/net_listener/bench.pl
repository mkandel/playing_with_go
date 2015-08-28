#!/usr/local/bin/perl -w
use strict;
use warnings;

use 5.012;

use Time::HiRes qw( gettimeofday );
my $start = gettimeofday();

use Carp;

use Getopt::Long;
Getopt::Long::Configure ("bundling");

use Data::Dumper;
# Some Data::Dumper settings:
local $Data::Dumper::Useqq  = 1;
local $Data::Dumper::Indent = 3;
local $Data::Dumper::Deparse  = 1;
use Benchmark qw{ cmpthese timethese };

my $debug = 1;
my $iters = 100;

GetOptions(
    "debug|d"        => \$debug,
) or die "Error: $!\n";

my $prog = $0;
$prog =~ s/^.*\///;

cmpthese( $iters, {
        script      => sub {
                `go run client.go -f lorem 2>&1 /dev/null`;
            },
        binary    => sub {
                `./client -f lorem 2>&1 /dev/null`;
            },
    }
);
#cmpthese( -10, {
#        one => sub { create_em() },
#        five => sub { foreach ( 1..5 ) { create_em() }},
#        ten => sub { foreach ( 1..10 ) { create_em() }},
#    }
#);

END{
    if ( $debug ){
        my $run_time = gettimeofday() - $start;
        print "$prog ran for ";
        if ( $run_time < 60 ){
            print "$run_time seconds.\n";
        } else {
            use integer;
            print $run_time / 60 . " minutes " . $run_time % 60
                . " seconds ($run_time seconds).\n";
        }
    }   
} 

__END__
