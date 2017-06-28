#!/bin/bash -e

# This script builds the Upspin commands and pushes them to release@upspin.io.
# It is executed by the release Docker container.
# The Docker container provides the Upspin repo in the /workspace directory and
# sets the environment COMMIT_HASH variable to the current Git commit hash of
# that repo.

# The commands to build and distribute.
# Command "upspin" must be one of these,
# as it is used to copy the files to release@upspin.io.
cmds="upspin upspinfs cacheserver"

# The operating systems and processor architectures to build for.
oses="darwin linux windows"
arches="amd64"

# Thet tree that contains the release binaries.
user="release@upspin.io"

echo 1>&2 "Repo has base path $1"
mkdir -p /gopath/src
cp -R /workspace/ /gopath/src/$1

export GOOS
export GOARCH
for GOOS in $oses; do
	for GOARCH in $arches; do
		for cmd in $cmds; do
			if [[ $GOOS == "windows" && $cmd == "upspinfs" ]]; then
				# upspinfs doesn't run on Windows.
				continue
			fi
			echo 1>&2 "Building $cmd for ${GOOS}_${GOARCH}"
			/usr/local/go/bin/go get -v upspin.io/cmd/$cmd
		done
	done
done

function upspin() {
	/gopath/bin/upspin -config=/config "$@"
}

for GOOS in $oses; do
	for GOARCH in $arches; do
		osarch="${GOOS}_${GOARCH}"
		srcdir="/gopath/bin"
		if [[ $osarch != "linux_amd64" ]]; then
			srcdir="$srcdir/$osarch"
		fi
		destdir="$user/all/$osarch/$COMMIT_SHA"
		for cmd in $cmds; do
			if [[ $GOOS == "windows" && $cmd == "upspinfs" ]]; then
				# upspinfs doesn't run on Windows.
				continue
			fi
			if [[ $GOOS == "windows" ]]; then
				# Windows commands have a ".exe" suffix.
				cmd="${cmd}.exe"
			fi
			src="$srcdir/$cmd"
			dest="$destdir/$cmd"
			link="$linkdir/$cmd"
			echo 1>&2 "Copying $src to $dest"
			upspin mkdir $destdir || echo 1>&2 "mkdir can fail if the directory exists"
			upspin cp $src $dest
		done
		link="$user/latest/$osarch"
		echo 1>&2 "Linking $link to $destdir"
		upspin rm $link || echo 1>&2 "rm can fail if the link does not already exist"
		upspin link $destdir $link
	done
done
