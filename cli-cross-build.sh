#!/usr/bin/env bash

package=$1
package_out=$2
version=$3

if [[ -z "$package" || -z "$package_out" ]]; then
  echo "usage: $0 <package> <package out> <version>"
  exit 1
fi

platforms=("windows/amd64" "windows/386"
           "darwin/amd64"
           "linux/amd64" "linux/386"
           "freebsd/amd64" "freebsd/386"
           "openbsd/amd64" "openbsd/386")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_base_path='build/'$GOOS'-'$GOARCH
    output_path=$output_base_path'/'$package_out
    output_name=$output_path'/'$package_out
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w -X main.version=$version" -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    # Copy skeleton
    cp -r 'skeleton' $output_path

    # Add README
    if [ $GOOS = "windows" ]; then
        readme_name='resources/cli/readme/windows/README.txt'
    else
        readme_name='resources/cli/readme/unix/README.txt'
    fi
    cp $readme_name $output_path

    # Add LICENSE
    cp 'LICENSE' $output_path'/LICENSE.txt'

    # Create archives for distributing
    package_name=$package_out'-'$GOOS'-'$GOARCH'-'$version
    cd $output_base_path
    if [ $GOOS = "windows" ]; then
        zip -r -9 $package_name'.zip' $package_out
    else
        tar -cvzf $package_name'.tar.gz' $package_out
    fi
    cd -
done
