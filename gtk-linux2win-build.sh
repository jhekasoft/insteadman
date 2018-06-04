#!/usr/bin/env bash

# Please read before compilation: https://github.com/gotk3/gotk3/wiki/Cross-Compiling

package=$1
package_out=$2
version=$3

if [[ -z "$package" || -z "$package_out" ]]; then
  echo "usage: $0 <package> <package out> <version>"
  exit 1
fi

goversioninfo=~/go/bin/goversioninfo
mingw_path='/usr/i686-w64-mingw32'
bin_path=$mingw_path'/bin'

dynamic_libs=('libatk-1.0-0.dll' 'libgdk_pixbuf-2.0-0.dll' 'libjpeg-8.dll'
          'libbz2-1.dll' 'libgio-2.0-0.dll' 'libpango-1.0-0.dll'
          'libcairo-2.dll' 'libglib-2.0-0.dll' 'libpangocairo-1.0-0.dll'
          'libcairo-gobject-2.dll' 'libgmodule-2.0-0.dll' 'libpangoft2-1.0-0.dll'
          'libepoxy-0.dll' 'libgobject-2.0-0.dll' 'libpangowin32-1.0-0.dll'
          'libexpat-1.dll' 'libgraphite2.dll' 'libpcre-1.dll'
          'libffi-6.dll' 'libgtk-3-0.dll' 'libpixman-1-0.dll'
          'libfontconfig-1.dll' 'libharfbuzz-0.dll' 'libpng16-16.dll'
          'libfreetype-6.dll' 'libiconv-2.dll' 'libstdc++-6.dll'
          'libgcc_s_sjlj-1.dll' 'libintl-8.dll' 'libwinpthread-1.dll'
          'libgdk-3-0.dll' 'libjasper.dll' 'zlib1.dll')

icons=('document-save-symbolic.*' 'pan-down-symbolic.*'
       'edit-clear-symbolic.*' 'pan-up-symbolic.*'
       'edit-delete-symbolic.*' 'view-refresh-symbolic.*'
       'media-playback-start-symbolic.*' 'open-menu-symbolic.*'
       'edit-clear-all-symbolic.*' 'list-add-symbolic.*'
       'list-remove-symbolic.*' 'go-up-symbolic.*'
       'go-down-symbolic.*')

icons_dirs=('16x16' '24x24' '32x32' '512x512' '8x8' 'scalable'
            '22x22' '256x256' '48x48' '64x64' '96x96')

platform='windows/386'
platform_split=(${platform//\// })
GOOS=${platform_split[0]}
GOARCH=${platform_split[1]}
output_base_path='build/'$GOOS'-'$GOARCH
output_path=$output_base_path'/'$package_out
output_name=$output_path'/'$package_out'.exe'
syso_name="resource.syso"
version_split=(${version//./ })
version_major=${version_split[0]}
version_minor=${version_split[1]}
version_patch=${version_split[2]}
version_build=0

iscc=~/.wine/drive_c/Program\ Files\ \(x86\)/Inno\ Setup\ 5/ISCC.exe

# Generate .syso resouces (executable icon, manifest, info)
$goversioninfo \
    -company="JhekaSoft" \
    -copyright="JhekaSoft" \
    -description="INSTEAD Manager" \
    -file-version=$version \
    -icon="resources/images/logo.ico" \
    -manifest="resources/windows/insteadman.manifest" \
    -o=$package"/"$syso_name \
    -product-name="InsteadMan" \
    -product-version=$version \
    -ver-major=$version_major \
    -ver-minor=$version_minor \
    -ver-patch=$version_patch \
    -ver-build=$version_build \
    -product-ver-major=$version_major \
    -product-ver-minor=$version_minor \
    -product-ver-patch=$version_patch \
    -product-ver-build=$version_build

# Building .exe
env CGO_ENABLED=1 \
	CC=i686-w64-mingw32-cc \
	GOOS=$GOOS \
	GOARCH=$GOARCH \
	go build -ldflags "-H=windowsgui -s -w -X main.version=$version" -o $output_name $package

# Remove .syso
rm $package'/'$syso_name

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

# Copy INSTEAD
cp -r 'resources/windows/instead' $output_path

# Add LICENSE
cp 'LICENSE' $output_path'/LICENSE.txt'

# Copy dynamic libs
for lib in "${dynamic_libs[@]}"
do
    cp $bin_path'/'$lib $output_path
done

# Copy icons
./gtk-copy-icons.sh $mingw_path $output_path

# Create archives for distributing
cd $output_base_path
package_name=$package_out'-'$GOOS'-'$GOARCH'-'$version
zip -r -9 '../'$package_name'.zip' $package_out
cd -

# Create installator
setup_script_path='resources/windows/setup.iss'
sed -e "s/{{version}}/$version/g" "$setup_script_path.in" > $setup_script_path
wine "$iscc" $setup_script_path
