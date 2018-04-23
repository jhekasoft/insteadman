#!/usr/bin/env bash

# Please read before compilation: https://github.com/gotk3/gotk3/wiki/Cross-Compiling

package=$1
package_out=$2
version=$3

if [[ -z "$package" || -z "$package_out" ]]; then
  echo "usage: $0 <package> <package out> <version>"
  exit 1
fi

platform="windows/386"
platform_split=(${platform//\// })
GOOS=${platform_split[0]}
GOARCH=${platform_split[1]}
output_base_path='build/'$GOOS'-'$GOARCH
output_path=$output_base_path'/'$package_out
output_name=$output_path'/'$package_out'.exe'

# TODO: copy GTK+ libs + resources

# Generate .syso resouces (executable icon, manifest, info)
go generate $package

# Building .exe
env CGO_ENABLED=1 \
	CC=i686-w64-mingw32-cc \
	GOOS=$GOOS \
	GOARCH=$GOARCH \
	go build -ldflags "-H=windowsgui -s -w" -o $output_name $package

# Copy skeleton
cp -r 'skeleton' $output_path

# Copy resources
resources_path=$output_path'/resources'
images_path=$resources_path'/images'
mkdir $resources_path
cp -r 'resources/gtk' $resources_path
cp 'resources/images/{logo.png}' $images_path

# Add LICENSE
cp 'LICENSE' $output_path'/LICENSE.txt'
