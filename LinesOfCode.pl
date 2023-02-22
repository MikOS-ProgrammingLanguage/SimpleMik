use strict;
use warnings;
use diagnostics;
use feature "say";
use v5.30.3;

sub parse_directories {
    state $lines_in_directories = 0;
    my ($directory, @suffixes) = @_;

    my @files = glob $directory;
    foreach (@files) {
        if (-d $_) {
            $_ = substr $_, 2;
            parse_directories("./".$_."/*", @suffixes);
        } else {
            my $ref = $_;
            my $applies = 0;
            foreach (@suffixes) {
                if ($_ eq substr $ref, -length($_)) {
                    $applies = 1;
                    last;
                }
            }
            open my $fh, "<", $_ if $applies
                or next;
            while (<$fh>) {$lines_in_directories++}
        }
    }
    return $lines_in_directories
}

say parse_directories "./*", ("go", "mik", "py", "pl", "Makefile");