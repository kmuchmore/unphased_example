#!/bin/bash

# Unzip the downloaded file
unzip protoc.zip -d protoc

# Move the binaries to /usr/local/bin
mv protoc/bin/* /usr/local/bin/

# # Clean up
# rm -rf protoc protoc.zip

echo "Protobuf has been installed"
