#!/bin/bash
set -eou pipefail

mkdir -p results

for i in $(seq 1 5); do
    echo "counter feedback k=$i"
    ./bin/rumormongering --peers 1000 --k $i --strategy cf > results/cf_${i}.csv 
    echo "blind random k=$i"
    ./bin/rumormongering --peers 1000 --k $i --strategy br > results/br_${i}.csv
done
