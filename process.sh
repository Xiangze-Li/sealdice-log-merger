#!/usr/bin/env bash

fnList=$(ls | grep .zip)

for fn in ${fnList[@]}; do
    if [[ $fn =~ ^QQ-Group[0-9]+_(.*)\.[0-9]+\.zip$ ]]; then
        grpName=${BASH_REMATCH[1]}
        echo "unzipping $grpName.zip"
        unzip -q -o -O=gb2312 $fn
        mv 文本log.txt $grpName.txt
        rm 海豹标准log-粘贴到染色器可格式化.txt
    fi
done
