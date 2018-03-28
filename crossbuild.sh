#!/usr/bin/env bash

package=$1
package_out=$2

if [[ -z "$package" || -z "$package_out" ]]; then
  echo "usage: $0 <package> <package out>"
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
    output_name='build/'$GOOS'-'$GOARCH'/'$package_out'/'$package_out
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w" -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
