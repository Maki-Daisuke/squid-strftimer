#!/usr/bin/env perl

use DateTime;
use DateTime::Format::ISO8601::Format;

my $format = DateTime::Format::ISO8601::Format->new;
sub format_time {
    my $dt = DateTime->from_epoch(epoch => shift);
    local $_ = $format->format_datetime($dt);
    my ($msec) = m/\.(\d+)Z$/;
    $msec .= "0"  while length $msec < 6;
    s/\.(\d+)Z$/.${msec}Z/;
    "[$_]";
}

while (<>) {
    s{^(\d+\.\d+)}{ format_time($1) }e;
    print;
}
