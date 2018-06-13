#!/usr/bin/env bash

# Please read before compilation: https://github.com/gotk3/gotk3/wiki/Cross-Compiling

package=$1
package_out=$2
version=$3
os=$4
arch=$5

if [[ -z "$package" || -z "$package_out" ]]; then
  echo "usage: $0 <package> <package out> <version> <arch>"
  exit 1
fi

instead_app_path='/Applications/Instead.app'

platform='darwin/amd64'
platform_split=(${platform//\// })
GOOS=${platform_split[0]}
GOARCH=${platform_split[1]}
output_base_path='build/'$GOOS'-'$GOARCH
output_path=$output_base_path'/'$package_out
tmp_path=$output_path'/tmp'
tmp_icons_path=$tmp_path'/icons'
tmp_contents_path=$tmp_path'/Contents'
tmp_macos_path=$tmp_contents_path'/MacOS'
tmp_resources_path=$tmp_contents_path'/Resources'

mkdir -p $tmp_path
mkdir -p $tmp_contents_path
mkdir -p $tmp_macos_path
mkdir -p $tmp_resources_path

# Copy icons
./gtk-copy-icons.sh $JHBUILD_PREFIX $tmp_icons_path

# Copy INSTEAD
cp -R $instead_app_path'/Contents/Frameworks' $tmp_contents_path'/' # Copy with symlinks
cp -r $instead_app_path'/Contents/Resources/__private__' $tmp_resources_path'/'
cp $instead_app_path'/Contents/MacOS/sdl-instead' $tmp_macos_path'/sdl-instead'

# Update version in plist
plist_path="./resources/darwin/bundle-gtk/Info-insteadman.plist"
sed -e "s/{{version}}/$version/g" "$plist_path.in" > $plist_path
