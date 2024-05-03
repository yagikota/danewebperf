#!/bin/bash

for location in tokyo frankfurt; do
    for num in `seq -f '%02g' 1 10`; do
        directory="${location}-v2-${num}"
        echo "Copying ${directory}..."
        mkdir -p "$directory"
        cd "$directory"
        aws s3 cp "s3://pageloadtime-results/${directory}" . --recursive --exclude "*" --include "pageload*" --profile default
        cd ..
    done
done
