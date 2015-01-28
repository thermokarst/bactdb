#!/usr/bin/env bash

# FROM: https://github.com/getlantern/flashlight-build/blob/devel/testandcover.bash

function die() {
    echo $*
    exit 1
}

export GOPATH=`pwd`:$GOPATH

echo "mode: count" > coverage.txt

ERROR=""

for pkg in `cat testpackages.txt`
do
    go test -v -covermode=count -coverprofile=profile_tmp.cov $pkg || ERROR="Error testing $pkg"
    tail -n +2 profile_tmp.cov >> coverage.txt || die "Unable to append coverage for $pkg"
done

if [ ! -z "$ERROR" ]
then
    die "Encountered error, last error was: $ERROR"
fi
