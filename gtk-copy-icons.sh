#!/usr/bin/env bash

input_path=$1
output_path=$2

if [[ -z "$input_path" || -z "$output_path" ]]; then
  echo "usage: $0 <input_path> <output_path>"
  exit 1
fi

icons=('document-save-symbolic.*' 'pan-down-symbolic.*'
       'edit-clear-symbolic.*' 'pan-up-symbolic.*'
       'edit-delete-symbolic.*' 'view-refresh-symbolic.*'
       'media-playback-start-symbolic.*' 'open-menu-symbolic.*'
       'edit-clear-all-symbolic.*' 'list-add-symbolic.*'
       'list-remove-symbolic.*' 'go-up-symbolic.*'
       'go-down-symbolic.*')

icons_dirs=('16x16' '24x24' '32x32' '512x512' '8x8' 'scalable'
            '22x22' '256x256' '48x48' '64x64' '96x96')

# Copy icons
theme_output_path=$output_path'/share/icons/Adwaita'
theme_path=$input_path'/share/icons/Adwaita'
for icon_dir in "${icons_dirs[@]}"
do
    mkdir -p $theme_output_path'/'$icon_dir'/actions'

    for icon in "${icons[@]}"
    do
        cp $theme_path'/'$icon_dir'/actions/'$icon $theme_output_path'/'$icon_dir'/actions'
    done
done
