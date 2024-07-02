#!/bin/bash

# turn on bash's job control
set -m

echo  "Running Migration Scripts"

./migrate-db up

echo  "Starting Service"

./service
