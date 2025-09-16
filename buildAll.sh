#!/usr/bin/env bash

# this script builds all go packages in the current directory

packages=$(go list ./...)
for package in $packages; do
    echo "building $package"
    if ! go build -o /dev/null "$package"; then
        read -p "Build failed for $package. Do you want to continue? (Y/n): " choice
        choice=${choice:-y}
        if [[ "$choice" != "y" && "$choice" != "Y" ]]; then
            echo "Exiting build process."
            exit 1
        fi
    fi
done
