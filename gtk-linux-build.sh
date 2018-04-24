#!/usr/bin/env bash

# Please read before compilation: https://github.com/gotk3/gotk3/wiki/Cross-Compiling

package=$1
package_out=$2
version=$3
arch=$4

if [[ -z "$package" || -z "$package_out" ]]; then
  echo "usage: $0 <package> <package out> <version> <arch>"
  exit 1
fi

platform='linux/'$arch
platform_split=(${platform//\// })
GOOS=${platform_split[0]}
GOARCH=${platform_split[1]}
output_base_path='build/'$GOOS'-'$GOARCH
output_path=$output_base_path'/'$package_out
output_name=$output_path'/'$package_out

# Building binary
env GOOS=$GOOS \
	GOARCH=$GOARCH \
	go build -ldflags "-s -w" -o $output_name $package

# Copy skeleton
cp -r 'skeleton' $output_path

# Copy resources
resources_path=$output_path'/resources'
images_path=$resources_path'/images'
mkdir $resources_path
mkdir $images_path
cp -r 'resources/gtk' $resources_path
cp 'resources/images/logo.png' $images_path
cp -r resources/windows/gtk/* $output_path

# Add README
readme_name='resources/gtk/readme/windows/README.txt'
cp $readme_name $output_path

# Add LICENSE
cp 'LICENSE' $output_path'/LICENSE.txt'

 Create archives for distributing
cd $output_base_path
package_name=$package_out'-'$GOOS'-'$GOARCH'-'$version
tar -cvzf $package_name'.tar.gz' $package_out
cd -
