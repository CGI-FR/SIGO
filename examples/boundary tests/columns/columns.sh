#!/bin/bash

rm -f log.txt

size_array=("2" "4" "8" "16" "32")
qi_array=("A,B" "A,B,C,D" "A,B,C,D,E,F,G,H" "A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P" "A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,AA,BB,CC,DD,EE,FF")

echo "NOANONYMIZER" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    START=$(date +%s)
    < test${i}_2.json | jq -c '.[]' | sigo -q ${qi_array[$i]} | jq -s > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test2_$i: (n=1000, x=${size_array[$i]}, method=NoAnomymizer) Execution time was $DIFF seconds." >> log.txt
done

echo "" >> log.txt
echo "GENERALIZATION" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    START=$(date +%s)
    < test${i}_2.json | jq -c '.[]' | sigo -q ${qi_array[$i]} -a general | jq -s > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test2_$i: (n=1000, x=${size_array[$i]}, method=Generalization) Execution time was $DIFF seconds." >> log.txt
done

echo "" >> log.txt
echo "AGGREGATION" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    START=$(date +%s)
    < test${i}_2.json | jq -c '.[]' | sigo -q ${qi_array[$i]} -a meanAggregation | jq -s > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test2_$i: (n=1000, x=${size_array[$i]}, method=Aggregation) Execution time was $DIFF seconds." >> log.txt
done

echo "" >> log.txt
echo "TOPBOTTOMCODING" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    START=$(date +%s)
    < test${i}_2.json | jq -c '.[]' | sigo -q ${qi_array[$i]} -a outlier | jq -s > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test2_$i: (n=1000, x=${size_array[$i]}, method=TopBottomCoding) Execution time was $DIFF seconds." >> log.txt
done

echo "" >> log.txt
echo "RANDOMNOISE" >> log.txt
echo "--------------------------------------------------------------------------------------------------------------" >> log.txt

for i in 1 2 3 4 5
do
    START=$(date +%s)
    < test${i}_2.json | jq -c '.[]' | sigo -q ${qi_array[$i]} -a laplaceNoise | jq -s > output.json
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "Test2_$i: (n=1000, x=${size_array[$i]}, method=RandomNoise) Execution time was $DIFF seconds." >> log.txt
done

rm -f output.json
