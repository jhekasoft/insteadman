#!/usr/bin/env bash

package="insteadman-cli.go"
package_name="insteadman-cli"

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64" "linux/386" "freebsd/amd64" "freebsd/386" "openbsd/amd64" "openbsd/386")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name='build/'$GOOS'-'$GOARCH'/insteadman-cli/'$package_name
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w"  -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
