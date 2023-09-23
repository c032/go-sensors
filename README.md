# go-sensors

## Tests

All tests are hidden behind opt-in tags.

To run all tests, use the `alltests` tag.

    go test -tags alltests

To run individual tests, check each file to see the name of the required
tag.

Tag names are based on the file name without any of the suffixes and
without the file extension, so you should be able to guess which tag is
required for each file just by looking at the file name.
