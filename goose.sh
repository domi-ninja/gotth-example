#!/bin/bash

# helper for running goose from command line
goose -dir ./sql/migrations sqlite3 ./data.db $@