#!/bin/bash

rm -f log_test.txt

START=$(date +%s)
< test_2.json | jq -c '.[]' | sigo -q A,B | jq -s > output.json
END=$(date +%s)
DIFF=$(( $END - $START ))
echo "Test: (n=2000000000, x=2, method=NoAnomymizer) Execution time was $DIFF seconds." >> log_test.txt

START=$(date +%s)
< test_3.json | jq -c '.[]' | sigo -q A,B | jq -s > output.json
END=$(date +%s)
DIFF=$(( $END - $START ))
echo "Test: (n=3000000000, x=2, method=NoAnomymizer) Execution time was $DIFF seconds." >> log_test.txt

START=$(date +%s)
< test_4.json | jq -c '.[]' | sigo -q A,B | jq -s > output.json
END=$(date +%s)
DIFF=$(( $END - $START ))
echo "Test: (n=4000000000, x=2, method=NoAnomymizer) Execution time was $DIFF seconds." >> log_test.txt

START=$(date +%s)
< test_5.json | jq -c '.[]' | sigo -q A,B | jq -s > output.json
END=$(date +%s)
DIFF=$(( $END - $START ))
echo "Test: (n=5000000000, x=2, method=NoAnomymizer) Execution time was $DIFF seconds." >> log_test.txt

rm -f output.json
