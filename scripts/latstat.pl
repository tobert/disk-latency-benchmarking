#!/usr/bin/env perl
$| = 1;

our @fields = qw(
    device
    read_complete read_merge read_sectors read_ms
    write_complete write_merge write_sectors write_ms
    io_in_progress io_ms io_ms_weighted
);

# 8  0 sda 298890 2980 5498843 92328 10123211 2314394 134218078 10756944 0 419132 10866136
# 8  5 sda5 5540 826 44511 1528 15558 55975 572334 68312 0 2932 69848
# 8 32 sdc 913492 273 183151490 8217340 2047310 0 37711114 1259728 0 1267508 9476068
# 8 16 sdb 2640 380 18329 2860 1751748 13461886 121702720 249041290 78 2654720 249048720
# 8 1  sda1 35383589 4096190 515794290 173085956 58990656 100542811 1276270912 205189188 0 135658516 378268412
# from Documentation/iostats.txt:
# Field  1 -- # of reads completed
# Field  2 -- # of reads merged, field 6 -- # of writes merged
# Field  3 -- # of sectors read
# Field  4 -- # of milliseconds spent reading
# Field  5 -- # of writes completed
# Field  6 -- # of writes merged
# Field  7 -- # of sectors written
# Field  8 -- # of milliseconds spent writing
# Field  9 -- # of I/Os currently in progress
# Field 10 -- # of milliseconds spent doing I/Os
# Field 11 -- weighted # of milliseconds spent doing I/Os
#
our $re_parse = qr/
    ^\s+   # leading whitespace
    \d+\s+ # major
    \d+\s+ # minor
    (?<device>\w+)\s+         # device
    (?<read_complete>\d+)\s+  # 1 [3]
    (?<read_merge>\d+)\s+     # 2
    (?<read_sectors>\d+)\s+   # 3 [5]
    (?<read_ms>\d+)\s+        # 4
    (?<write_complete>\d+)\s+ # 5 [7]
    (?<write_merge>\d+)\s+    # 6
    (?<write_sectors>\d+)\s+  # 7 [9]
    (?<write_ms>\d+)\s+       # 8
    (?<io_in_progress>\d+)\s+ # 9
    (?<io_ms>\d+)\s+          # 10
    (?<io_ms_weighted>\d+)    # 11
/x;

# if there are no command line args, print all devices
# otherwise match any devices for inclusion
our %include_devs = ();
open(my $fh, "< /proc/diskstats") || die "cannot open /proc/diskstats";
while (my $line = <$fh>) {
    if ($line =~ /$re_parse/) {
        my $dev = $+{device};
        if (@ARGV == 0) {
            $include_devs{$dev} = 1;
        }
        else {
            foreach my $want (@ARGV) {
                if ($dev =~ /$want/) {
                    $include_devs{$dev} = 1;
                }
            }
        }
    }
}
close $fh;

our $count = 0;
our %start = ();
our %prev  = ();

while (1) {
    my $ts = time;

    open(my $fh, "< /proc/diskstats") || die "cannot open /proc/diskstats";
    my %now = ();

    while (my $line = <$fh>) {
        if ($line =~ /$re_parse/ ) {
            # ignore lines for devices we don't care about
            next unless exists $include_devs{$+{device}};

            # copy the captured fields into a hash
            my %data = map { $_ => $+{$_} } @fields;
            $data{ts} = $ts;
            $data{count} = $count;

            # store the current data in %now by device name
            $now{$+{device}} = \%data;
        }
    }
    close $fh;

    # prime %start and %prev on the first iteration
    if ($count == 0) {
        %start = %now;
        %prev  = %now;
    }

    foreach my $dev (sort keys %include_devs) {
        foreach my $f (@fields) {
            my $delta = $now{$dev}->{$f} - $prev{$dev}->{$f};

            printf "% 6s % 16s => % 10d\n",
                $dev, $f, $delta;
        }
    }

    %prev = %now;
    $count++;

    sleep 2;
    print "----------------------------------\n";
}

# vim: et ts=4 sw=4 ai smarttab
