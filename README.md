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

### V1 Implemantetion

- First download csv files
- Calculate diff based on local csv files
- Generate array with added values and removed values

### V2 Implemantetion

- Open conection with first remote csv file
- Proccess on demand and save data on diff map
- Open conection with second remote csv file
- Proccess on demand and save data on diff map
- Generate array with added values and removed values

## Results

V1 is the winner! Used less memory during the process.

### V1

- With csv with 600k lines:

  - Tempo de execução: 2m28.475000412s
  - Consumo de memória:

        Alloc = 258 MiB
        TotalAlloc = 412 MiB
        Sys = 338 MiB
        NumGC = 15

- With csv with 4.5kk lines:

  - Tempo de execução: 28m5.439666081s
  - Consumo de memória:

        Alloc = 2352 MiB
        TotalAlloc = 3269 MiB
        Sys = 2567 MiB
        NumGC = 31

### V2

- With csv with 600k lines:

  - Tempo de execução: 2m10.968240411s
  - Consumo de memória:

        Alloc = 263 MiB
        TotalAlloc = 412 MiB
        Sys = 338 MiB
        NumGC = 15

- With csv with 4.5kk lines:

  - panic: read: connection reset by peer