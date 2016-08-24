# MoonLegend
    Powered by Hyper2Stack, special thanks to HP7

## build
    ./build.sh

## unit test
    go test controller/...

## start moonlegend
    ./ctrl start

## cmd line client
    ./bin/ml help

## example
Make sure following conditions are satisfied:

    [moon agent](https://github.com/Hyper2Stack/Moon/blob/master/README.md) is installed
    docker is installed
    ports 3306, 13306, 8000, 8080 are available
    moonlegend server is not running

Run below script to start the example

    ./example/init.sh
    ./example/run.sh env
    ./example/run.sh config_file
    ./example/destroy.sh
