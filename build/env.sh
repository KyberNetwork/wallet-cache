#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
appdir="$workspace/src/github.com/KyberNetwork"
if [ ! -L "$appdir/cache" ]; then
    mkdir -p "$appdir"
    cd "$appdir"
    ln -s ../../../../../. cache
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$appdir/cache"
PWD="$appdir/cache"

# Launch the arguments with the configured environment.
exec "$@"
