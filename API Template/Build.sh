#!/bin/bash

# This is a script that encapsulated the "make" funcationality so
# that our build process can be cross platform.

# Stop the script if it errors at any phase
set -e

# Map the command to a variable
CMD=""
ENV_FILE="dev.env"

# Directory Structure Variables
PROJECT_NAME=${PWD##*/} #This will grab the directory name which needs to be the same as the project name
PROJECT_VERSION=$(git describe --always --dirty)
PROJECT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
PROJECT_BUILD_TIME="$(date +'%Y-%m-%d@%H:%M:%S')"
PROJECT_SHA_1=$(git rev-list -1 HEAD)

DIST_DIR="./dist"
DIST_CONFIG_DIR="$DIST_DIR/config"
DIST_PUBLIC_DIR="$DIST_DIR/public"
LOCAL_PUBLIC_DIR="public"
LOCAL_CONFIG_DIR="config"

# Linker Flags
VERSION_FLAG="-X $PROJECT_NAME/env.Version=$PROJECT_VERSION"
BRANCH_FLAG="-X $PROJECT_NAME/env.Branch=$PROJECT_BRANCH"
BUILD_TIME_FLAG="-X $PROJECT_NAME/env.BuildTime=$PROJECT_BUILD_TIME"
ENCRYPTION_FLAG="-X $PROJECT_NAME/env.SHA1=$PROJECT_SHA_1"

# Shared Functions
makeDeps() {
    dep ensure -v
}

makeDist() {
    if [ ! -d $DIST_DIR ]; then
        mkdir $DIST_DIR
    fi

    if [ ! -d $DIST_PUBLIC_DIR ]; then
        mkdir $DIST_PUBLIC_DIR
    fi

    cp -r "$LOCAL_PUBLIC_DIR"/* $DIST_PUBLIC_DIR

    if [ ! -d $DIST_CONFIG_DIR ];then
        mkdir $DIST_CONFIG_DIR
    fi

    cp -r "$LOCAL_CONFIG_DIR"/* $DIST_CONFIG_DIR
}

help() {
    # Help is built based on lines that start with ## in a <command>:<description> format. As commands are added you can
    # make sure they are added to help by adding the comment line similar to those below.
	echo
	echo "Usage: ./Build.sh optional:<Environment File> [<command1> <command2> ...]"
	echo
	echo "Select a command to run for '$PROJECT_NAME':"
	echo
	sed  -n 's/^##//p' ${PWD}/$0 | column -t -s ':' | sed -e 's/^/ /'
	echo
}

runGoCommand(){
## clean: Remove all build artifacts and generated files.
    if [ "$CMD" == "clean" ]; then
        go clean -x
        rm -rf dist/
        rm -rf vendor/
        rm -rf Gopkg.lock

## deps: Create the build container and pull down necessary dependencies.
    elif [ "$CMD" == "deps" ]; then
        makeDeps

## server: Build the application.
    elif [ "$CMD" == "server" ]; then
        makeDeps
        makeDist
        go build -v -ldflags "$VERSION_FLAG $BRANCH_FLAG $BUILD_TIME_FLAG $ENCRYPTION_FLAG" -o $DIST_DIR/$PROJECT_NAME

## start: Start the application.
    elif [ "$CMD" == "start" ]; then
        echo "Loading environment '$ENV_FILE':"
        export LOG_FACILITY=$PROJECT_NAME
        cat $ENV_FILE
        eval $(sed 's/^/export /' $ENV_FILE)
        echo
        $DIST_DIR/$PROJECT_NAME

## test: Test the Go code.
    elif [ "$CMD" == "test" ]; then
        makeDeps
        go test -v ./...
    elif [ "$CMD" == "help" ]; then
        help
    else
        echo
        echo "Command '$CMD' not supported"
        help
    fi
}

# Process the environment variable
if [[ $1 == *".env" ]]; then
    ENV_FILE="$1"
fi

# Process the commands
for arg in "$@"
do
    if [[ "$arg" != *".env" ]]; then
        CMD=$arg
        runGoCommand
    fi
done