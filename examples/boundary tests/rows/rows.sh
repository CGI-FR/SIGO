#!/bin/bash

rm -f log.txt

size_array=("100" "1000" "10000" "100000" "1000000")

echo "NOANONYMIZER" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    pimo < test.json -c masking.yml -r ${size_array[$i]} > test${i}_1.json
    START=$(date +%s)
    sigo -q A,B < test${i}_1.json > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test1_$i: (n=${size_array[$i]}, x=2, method=NoAnomymizer) Execution time was $DIFF seconds." >> log.txt
    rm -f test${i}_1.json
    rm -f output.json
done

echo "" >> log.txt
echo "GENERALIZATION" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    pimo < test.json -c masking.yml -r ${size_array[$i]} > test${i}_1.json
    START=$(date +%s)
    sigo -q A,B -a general < test${i}_1.json > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test1_$i: (n=${size_array[$i]}, x=2, method=Generalization) Execution time was $DIFF seconds." >> log.txt
    rm -f test${i}_1.json
    rm -f output.json
done

echo "" >> log.txt
echo "AGGREGATION" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    pimo < test.json -c masking.yml -r ${size_array[$i]} > test${i}_1.json
    START=$(date +%s)
    sigo -q A,B -a meanAggregation < test${i}_1.json > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test1_$i: (n=${size_array[$i]}, x=2, method=Aggregation) Execution time was $DIFF seconds." >> log.txt
    rm -f test${i}_1.json
    rm -f output.json
done

echo "" >> log.txt
echo "TOPBOTTOMCODING" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    pimo < test.json -c masking.yml -r ${size_array[$i]} > test${i}_1.json
    START=$(date +%s)
    sigo -q A,B -a outlier < test${i}_1.json > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test1_$i: (n=${size_array[$i]}, x=2, method=TopBottomCoding) Execution time was $DIFF seconds." >> log.txt
    rm -f test${i}_1.json
    rm -f output.json
done

echo "" >> log.txt
echo "RANDOMNOISE" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    pimo < test.json -c masking.yml -r ${size_array[$i]} > test${i}_1.json
    START=$(date +%s)
    sigo -q A,B -a laplaceNoise < test${i}_1.json > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test1_$i: (n=${size_array[$i]}, x=2, method=RandomNoise) Execution time was $DIFF seconds." >> log.txt
    rm -f test${i}_1.json
    rm -f output.json
done

