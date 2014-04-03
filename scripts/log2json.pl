#!/usr/bin/env perl
# quickly convert the fio clat log to JSON

our $infile = shift(@ARGV) || die "2 arguments required but 0 given: infile outfile";
our $outfile = shift(@ARGV) || die "2 arguments required but 1 given: infile outfile";

open(my $in, "< $infile") || die "Could not open $infile for read: $!";
open(my $out, "> $outfile") || die "Could not open $outfile for write: $!";

print $out "[\n";

my $do_comma = undef;
while (my $line = <$in>) {
    my($time_usec, $lat_usec, $a, $b) = split /\s*,\s*/, $line, 4;
    if ($do_comma) {
        print $out ",\n";
    }
    print $out "{\"x\":$time_usec,\"y\":$lat_usec}";
    $do_comma = 1;
}

print $out "\n]\n";

close $in;
close $out;


# vim: et ts=4 sw=4 ai smarttab
