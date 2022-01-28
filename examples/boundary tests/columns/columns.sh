#!/bin/bash

rm -f log.txt

size_array=("2" "4" "8" "16" "32")
qi_array=("A,B" "A,B,C,D" "A,B,C,D,E,F,G,H" "A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P" "A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,AA,BB,CC,DD,EE,FF")

echo "NOANONYMIZER" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    pimo < test${i}.json -c masking${i}.yml -r 1000 > test${i}_2.json
    START=$(date +%s)
    sigo -q ${qi_array[$i]} < test${i}_2.json > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test2_$i: (n=1000, x=${size_array[$i]}, method=NoAnomymizer) Execution time was $DIFF seconds." >> log.txt
    rm -f test${i}_2.json
    rm -f output.json
done

echo "" >> log.txt
echo "GENERALIZATION" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    pimo < test${i}.json -c masking${i}.yml -r 1000 > test${i}_2.json
    START=$(date +%s)
    sigo -q ${qi_array[$i]} -a general < test${i}_2.json > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test2_$i: (n=1000, x=${size_array[$i]}, method=Generalization) Execution time was $DIFF seconds." >> log.txt
    rm -f test${i}_2.json
    rm -f output.json
done

echo "" >> log.txt
echo "AGGREGATION" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    pimo < test${i}.json -c masking${i}.yml -r 1000 > test${i}_2.json
    START=$(date +%s)
    sigo -q ${qi_array[$i]} -a meanAggregation < test${i}_2.json > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test2_$i: (n=1000, x=${size_array[$i]}, method=Aggregation) Execution time was $DIFF seconds." >> log.txt
    rm -f test${i}_2.json
    rm -f output.json
done

echo "" >> log.txt
echo "TOPBOTTOMCODING" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    pimo < test${i}.json -c masking${i}.yml -r 1000 > test${i}_2.json
    START=$(date +%s)
    sigo -q ${qi_array[$i]} -a outlier < test${i}_2.json > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test2_$i: (n=1000, x=${size_array[$i]}, method=TopBottomCoding) Execution time was $DIFF seconds." >> log.txt
    rm -f test${i}_2.json
    rm -f output.json
done

echo "" >> log.txt
echo "RANDOMNOISE" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    pimo < test${i}.json -c masking${i}.yml -r 1000 > test${i}_2.json
    START=$(date +%s)
    sigo -q ${qi_array[$i]} -a laplaceNoise < test${i}_2.json > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test2_$i: (n=1000, x=${size_array[$i]}, method=RandomNoise) Execution time was $DIFF seconds." >> log.txt
    rm -f test${i}_2.json
    rm -f output.json
done

