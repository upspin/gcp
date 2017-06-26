#!/bin/bash -e

# This script builds the Upspin commands and pushes them to release@upspin.io.
# It is executed by the release Docker container.

# The commands to build and distribute.
# Command "upspin" must be one of these,
# as it is used to copy the files to release@upspin.io.
cmds="upspin upspinfs cacheserver"

echo "Repo has base path $1"
mkdir -p /gopath/src
cp -R /workspace/ /gopath/src/$1

for cmd in $cmds; do
	echo "Building $cmd"
	/usr/local/go/bin/go get -v upspin.io/cmd/$cmd
done

up="/gopath/bin/upspin -config=/config"
user="release@upspin.io"
osarch="linux_amd64"
destdir="$user/all/$osarch/$COMMIT_SHA"
linkdir="$user/latest/$osarch"
for cmd in $cmds; do
	echo "Copying $cmd to release@upspin.io"
	dest="$destdir/$cmd"
	link="$linkdir/$cmd"
	$up mkdir $destdir || echo 1>&2 "mkdir can fail if the directory exists"
	$up cp /gopath/bin/$cmd $dest
	$up rm $link || echo 1>&2 "rm can fail if the link does not already exist"
	$up link $dest $link
done
