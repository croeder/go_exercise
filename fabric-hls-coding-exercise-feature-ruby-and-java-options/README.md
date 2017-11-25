# Overview

Today we're going to build a new microservice which will be part of our next generation media delivery platform.

We strongly recommend that you implement this service in Go (golang). We've included a [barebones Go application](./main.go) to get you started.  There are also barebones [Ruby](./main.rb) and [Java](./Main.java) applications if you prefer to implement in Ruby or Java.

You may use any online resource you want during the exercise.

# Background

Most online video delivery today uses a variant of a "segmented" (chunked) delivery technology.

In segmented delivery, the video player first requests a Manifest file, which contains a list of the URLs of small chunks of media for the player to play, and the duration of each of those segments.

The delivery technology we're using is called "HTTP Live Streaming" (HLS). Here are some examples of some manifests.

1. [2 Second segments](./cucumber/fixtures/expected-manifest-2s-seg.m3u8)
2. [10 Second segments](./cucumber/fixtures/expected-manifest-10s-seg.m3u8).

We're going to build a microservice which generates these HLS manifest files based on some metadata.

<sup>If you want to know more about HLS, [read the Apple developer's guide, here](https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/StreamingMediaGuide/Introduction/Introduction.html), or [the latest specification here](https://tools.ietf.org/html/draft-pantos-http-live-streaming-19).</sup>

# Task 1

Get the code, put it in place, get the BDD tests running.

NOTE: You will likely be given this repo as a zip file. Place it here if you're implementing in Go:

    $GOPATH/src/github.com/zencoder/fabric-hls-coding-exercise

Install Ruby and Bundler (Ubuntu, YMMV) for running the BDD test suite (depending on your Ruby setup, you may need to use sudo for `gem` as well):

    sudo apt-get install ruby ruby-dev
    gem install bundler

If you're implementing in Ruby, configure bundler in the application directory:

    bundle install

Update the `run.sh` file to uncomment the run command for the language you're implementing in (Go is the default).

Installi Cucumber and its dependencies:

    cd cucumber
    bundle install

Running all the tests:

    bundle exec cucumber

Running a specific task:

    bundle exec cucumber -t @task2

NOTE: The Cucumber test runner will run your application using the `run.sh` file.  You do not need the application running ahead of time.

# Task 2

The BDD tests will time out right now waiting for the application to start. We wait 10 seconds for the healthcheck endpoint on the application to start.

Implement a healthcheck handler in the application at `http://localhost:1337/healthcheck` which returns a 200

    GIVEN the application is running
    WHEN call GET on "/healthcheck"
    THEN the result should be a "200"

# Task 3

Our HTTP microservice will use JSON metadata files stored on disk to provide the information which will allow us to structure the manifest.

Before we get dive head first into generating manifests, lets start by exposing the metadata files via our API.

Add a new route to the Microservice which returns the contents of the JSON file, on the endpoint as follows:

    GET /metadata/${ID}

The microservice should behave as follows:

    GIVEN the metadata file "simple.json"
    WHEN call GET on "/metadata/simple"
    THEN the result should be a "200"
    AND the response should be JSON equivalent to "simple.json"
    AND the content type should be "application/json"

You can run the behavioural tests for this task by running

    bundle exec cucumber -t @task3

# Task 4

To generate the Manifest, we'll need to parse the JSON file, and use the data in it.

As an intermediary step, we're going to create a new API which returns the full duration of the media in milliseconds.

This information is obtained by adding together each of the "duration" entries in the "atoms" list in the JSON metadata file.

    GIVEN the metadata file "simple.json"
    WHEN call GET on "/duration/simple"
    THEN the result should be a "200"
    AND the response should match "simple-duration.txt"
    AND the content type should be "text/plain"

You can run the behavioural tests for this task by running

    bundle exec cucumber -t @task4

# Task 5

Build a HTTP microservice which can generate these manifest files based on a JSON metadata file stored on disk. The request URL should look something like this:

    GET /manifest/${ID}/2s.m3u8

Where $ID is the name of the metadata on disk, with the ".json" extension removed. EG:

    GET /manifest/simple/2s.m3u8

Would use the metadata file `simple.json`

In the first use case, for each entry in the "atoms" list, a segment should be added to the manifest.

    GIVEN the metadata file "simple.json"
    WHEN call GET on "/manifest/simple/2s.m3u8"
    THEN the result should be a "200"
    AND the response should match "expected-manifest-2s-seg.m3u8"
    AND the content type should be "application/x-mpegURL"

You can run the behavioural tests for this task by running

    bundle exec cucumber -t @task5

# Task 6

A request for a non-existent metadata file should result in 404.

    GIVEN the metadata file "idonotexist.json"
    WHEN call GET on "/manifest/simple/2s.m3u8"
    THEN the result should be a "404"

You can run the behavioural tests for this task by running

    bundle exec cucumber -t @task6

# Task 7

This is only part of the story. Our dynamic media server is able change the duration of the media segments on the fly. Our manifest server needs to be able to take a requested segment duration, and render a manifest with that duration set for the segments. For this, the URL structure would need to contain an extra parameter:

    GET /manifest/${ID}/${SEGMENT_DURATION}s.m3u8

So for example,

    GET /manifest/simple/10s.m3u8

This would produce a manifest where the segments were roughly 10 seconds in duration.

    GIVEN the metadata file "simple.json"
    WHEN call GET on "/manifest/simple/10s.m3u8"
    THEN the result should be a "200"
    AND the response should match "expected-manifest-10s-seg.m3u8"

You can run the behavioural tests for this task by running

    bundle exec cucumber -t @task7
