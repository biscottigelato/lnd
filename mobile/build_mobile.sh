#!/usr/bin/env bash
set -e -u -o pipefail # Fail on error

dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
echo "current dir is $dir"
cd $dir

arg=${1:-}

package="github.com/lightningnetwork/lnd/mobile"

mobile_dir="$dir"
build_dir="$mobile_dir/build"
mkdir -p $build_dir

# Generate API bindings.
api_filename="api_generated.go"
echo "Generating mobile bindings ($api_filename)."
go generate
echo "Done."

ios_dir="$build_dir/ios"
mkdir -p $ios_dir
ios_dest="$ios_dir/Lndmobile.framework"
echo "Building for iOS ($ios_dest)..."
"$GOPATH/bin/gomobile" bind -target=ios -tags="ios" -v -o "$ios_dest" "$package"

android_dir="$build_dir/android"
mkdir -p $android_dir
android_dest="$android_dir/Lndmobile.aar"
echo "Building for Android ($android_dest)..."
"$GOPATH/bin/gomobile" bind -target=android -tags="android" -v -o "$android_dest" "$package"

echo "Cleaning up."
rm $api_filename
