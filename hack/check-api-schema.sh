#!/usr/bin/env sh

if [ -z "$(git status -s src/ui/lib/api-types.ts)" ]; then
    echo "Success!"
    exit 0
else
    git diff src/ui/lib/api-types.ts
    exit 1
fi
