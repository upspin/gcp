#!/bin/bash -e

# Usage:
#   $ makerelease.sh <repo-base>
# where repo-base is one of "upspin.io" or "augie.upspin.io".
#
# This script builds Upspin commands for multiple platforms and pushes them to
# the release@upspin.io Upspin tree.
# It is executed by the release Docker container.
#
# The Docker container is run by Google Container Builder, which provides the
# Git repository (nominated as repo-base) in the /workspace directory and sets
# the environment variable COMMIT_SHA to the current Git commit hash of that
# repo.
#
# The Docker container is built atop xgo (https://github.com/karalabe/xgo)
# which is a framework for cross-compiling cgo-enabled binaries. Its magic
# environment variables are:
#  EXT_GOPATH,	the location of the Go workspace, and
#  TARGETS,	a space-separated list of os/arch combinations.

# The operating systems and processor architectures to build for.
oses="darwin linux windows"
arches="amd64"

# The tree that contains the release binaries.
user="release@upspin.io"

# The repository to build for.
repo="$1"

case "$repo" in
upspin.io|augie.upspin.io)
	;;
*)
	echo >&2 "must supply upspin.io or augie.upspin.io as first argument"
	exit 1
esac

echo 1>&2 "Repo has base path ${repo}"
export EXT_GOPATH="/gopath"
mkdir -p $EXT_GOPATH/src
cp -R /workspace/ $EXT_GOPATH/src/$repo
mkdir /build

# A short name for this repo; used in the's release path name.
shortrepo="upspin"

# The commands to build and distribute for repo "upspin.io".
cmds="upspin upspinfs cacheserver"

# For other repos, set cmds appropriately,
# and perform code generation.
case "$repo" in
augie.upspin.io)
	shortrepo="augie"
	cmds="upspin-ui"
	GOPATH="$EXT_GOPATH" go generate augie.upspin.io/cmd/upspin-ui/static
esac

# Generate the version strings for the commands.
if [[ "$repo" != "upspin.io" ]]; then
	GOPATH="$EXT_GOPATH" go get -d upspin.io/cmd/upspin
fi
GOPATH="$EXT_GOPATH" go generate -run make_version upspin.io/version

# Build the upspin tool, used to copy the files to release@upspin.io.
if [[ "$repo" != "upspin.io" ]]; then
	TARGETS="linux/amd64" $BUILD upspin.io/cmd/upspin
fi

# Build cmds for oses and arches.
for cmd in $cmds; do
	TARGETS=""
	for GOOS in $oses; do
		for GOARCH in $arches; do
			if [[ $GOOS == "windows" && $cmd == "upspinfs" ]]; then
				# upspinfs doesn't run on Windows.
				continue
			fi
			TARGETS="$TARGETS ${GOOS}/${GOARCH}"
		done
	done
	echo 1>&2 "Building $cmd for $TARGETS"
	export TARGETS
	$BUILD ${repo}/cmd/$cmd
done

# Upload the built commands to the release@upspin.io tree.
function upspin() {
	/build/upspin-linux-amd64 -config=/config "$@"
}
for GOOS in $oses; do
	for GOARCH in $arches; do
		osarch="${GOOS}_${GOARCH}"
		destdir="${user}/${shortrepo}/commits/${osarch}/$COMMIT_SHA"

		for cmd in $cmds; do
			if [[ $GOOS == "windows" && $cmd == "upspinfs" ]]; then
				# upspinfs doesn't run on Windows.
				continue
			fi
			# Use wildcard between os and arch to match OS version
			# numbers in the binaries produced by xgo's build script.
			src="/build/${cmd}-${GOOS}-*${GOARCH}"
			if [[ $GOOS == "windows" ]]; then
				# Windows commands have a ".exe" suffix.
				src="${src}.exe"
				cmd="${cmd}.exe"
			fi
			dest="${destdir}/${cmd}"
			link="${linkdir}/${cmd}"
			echo 1>&2 "Copying $src to $dest"
			upspin mkdir -p "$destdir" || echo 1>&2 "mkdir can fail if the directory exists"
			upspin cp "$src" "$dest"
		done

		# Re-link the 'latest' link.
		link="${user}/${shortrepo}/latest/${osarch}"
		echo 1>&2 "Linking $link to $destdir"
		upspin rm "$link" || echo 1>&2 "rm can fail if the link does not already exist"
		upspin link "$destdir" "$link"
	done
done
