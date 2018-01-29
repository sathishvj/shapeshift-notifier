#!/bin/bash

package_name="shapeshift-notifier"

platforms=("windows/amd64" "linux/386" "darwin/amd64")

echo "Deleting existing releases folder and recreating."
find ./releases -mindepth 1 -delete
mkdir -p releases

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name='releases/'$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi  

	echo "Building $output_name"
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
    if [ $? -ne 0 ]; then
        echo "An error has occurred! Aborting the script execution..."
        exit 1
    fi
done

