# Distributed Algorithms project 2015

## Compilation from sources

Given a working (go installation)[http://golang.org/doc/install] the compilation procedure is just

    ~ $ go get github.com/dallagi/rumormongering

This downloads and compiles the project. The resulting binary will be bin/rumormongering
Note that this requires a working (GOPATH)[https://golang.org/doc/code.html#GOPATH].
A precompiled binary version is available in the release tab.

## Generating the data

The ```gen_data.sh``` script is provided to generate the data used in the report.
The script expects to be placed in the same folder as the program executable.
The resulting csv files will be put in the results folder and will be named according to the {strategy}_{k}.csv patter.
Due to the non-determinist nature of the simulation the data generated will be slightly different from the data used in the report.

    ~/test $ ls
    rumormongering gen_data.sh

    ~/test $ ./gen_data.sh
    counter feedback k=1
    blind random k=1
    counter feedback k=2
    blind random k=2
    counter feedback k=3
    blind random k=3
    counter feedback k=4
    blind random k=4
    counter feedback k=5
    blind random k=5

    ~/test $ ls results/
    br_1.csv  br_2.csv  br_3.csv  br_4.csv  br_5.csv  cf_1.csv  cf_2.csv  cf_3.csv  cf_4.csv  cf_5.csv
    
