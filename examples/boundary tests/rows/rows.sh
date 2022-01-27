#!/bin/bash

rm -f log.txt

size_array=("100" "1000" "10000" "100000" "1000000")

echo "NOANONYMIZER" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    START=$(date +%s)
    < test${i}_1.json | jq -c '.[]' | sigo -q A,B | jq -s > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test1_$i: (n=${size_array[$i]}, x=2, method=NoAnomymizer) Execution time was $DIFF seconds." >> log.txt
done

echo "" >> log.txt
echo "GENERALIZATION" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    START=$(date +%s)
    < test${i}_1.json | jq -c '.[]' | sigo -q A,B -a general | jq -s > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test1_$i: (n=${size_array[$i]}, x=2, method=Generalization) Execution time was $DIFF seconds." >> log.txt
done

echo "" >> log.txt
echo "AGGREGATION" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    START=$(date +%s)
    < test${i}_1.json | jq -c '.[]' | sigo -q A,B -a meanAggregation | jq -s > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test1_$i: (n=${size_array[$i]}, x=2, method=Aggregation) Execution time was $DIFF seconds." >> log.txt
done

echo "" >> log.txt
echo "TOPBOTTOMCODING" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    START=$(date +%s)
    < test${i}_1.json | jq -c '.[]' | sigo -q A,B -a outlier | jq -s > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test1_$i: (n=${size_array[$i]}, x=2, method=TopBottomCoding) Execution time was $DIFF seconds." >> log.txt
done

echo "" >> log.txt
echo "RANDOMNOISE" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    START=$(date +%s)
    < test${i}_1.json | jq -c '.[]' | sigo -q A,B -a laplaceNoise | jq -s > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test1_$i: (n=${size_array[$i]}, x=2, method=RandomNoise) Execution time was $DIFF seconds." >> log.txt
done

rm -f output.json
