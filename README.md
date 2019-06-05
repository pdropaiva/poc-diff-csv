# poc-diff-csv

PoC (proof of concept) to make difference of two csv file with same structure.

## Dependencies

- Go version 1.12 installed
- Make

## How to execute

- Create env file based on env.sample
- `make poc`

## Introduction

This app will access one API to get presigned URL to access two csv files and make difference between they. This difference will be two arrays one with added values and another with removed values.
