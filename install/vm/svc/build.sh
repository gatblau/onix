#!/usr/bin/env bash
git clone https://github.com/gatblau/onix.git
cd onix
mvn package
cd ..
cp onix/target/onix*.jar .
rm -rf onix
