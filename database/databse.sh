#!/bin/bash -e

sqlite3 forum.db < create.sql
sqlite3 forum.db < insert.sql
sqlite3 forum.db < select.sql
