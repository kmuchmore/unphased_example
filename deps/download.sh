#!/bin/bash

# Get the latest release version from GitHub API
latest_version=$(curl --silent "https://api.github.com/repos/protocolbuffers/protobuf/releases/latest" | jq -r .tag_name)

# Remove the 'v' from the beginning of the latest_version string
clean_version=${latest_version#v}

# Construct the download URL
download_url="https://github.com/protocolbuffers/protobuf/releases/download/${latest_version}/protoc-${clean_version}-linux-x86_64.zip"

# Download the latest version
curl -L -o protoc.zip $download_url

echo "Protobuf ${latest_version} has been downloaded"