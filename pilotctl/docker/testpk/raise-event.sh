#!/bin/bash

echo -n "This is a test alert" | nc -4u -w1 localhost 1514
echo -n "This is another test alert" | nc -4u -w1 localhost 1514
echo -n "This is a final test alert" | nc -4u -w1 localhost 1514
