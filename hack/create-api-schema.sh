#!/usr/bin/env sh

# Create the json schema for the API and use it to create the typescript definitions
go run main.go internal api-schema | npx quicktype -s schema -o src/ui/lib/api-types.ts
