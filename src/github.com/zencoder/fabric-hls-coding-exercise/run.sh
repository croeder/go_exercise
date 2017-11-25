#!/bin/bash

if ! cd "$(dirname "${BASH_SOURCE[0]}")"; then
	echo 'Unable to cd to run.sh dir'
	exit 1
fi

# For GoLang applications:
exec go run main.go

# For Ruby applications:
# bundle exec ruby main.rb

# For Java applications:
# javac Main.java && java Main
