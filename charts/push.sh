#!/bin/bash

REPO=${1:-"oci://localhost:5000/charts"}

set -e

for folder in */; do
	folder=${folder%/}  # Remove trailing slash
	echo "Processing folder: $folder"
	pushd "$folder" > /dev/null
	helm package .
	helm push *.tgz "$REPO"
	popd > /dev/null
done
