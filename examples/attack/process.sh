#!/bin/bash

set -e

echo "version: v1
masking:
  - selector:
      jsonpath: \"individual_id\"
    masks:
      - add: \"\"
      - incremental:
          start: 1
          increment: 1
  - selector:
      jsonpath: firstname
    masks:
      - add: \"\"
      - randomChoiceInUri: \"pimo://nameFR\"
  - selector:
      jsonpath: city
    masks:
      - add: \"\"
      - randomChoice:
          - \"Paris\"
          - \"Nantes\"
          - \"Toulouse\"
          - \"Marseille\"
          - \"Lyon\"
  - selector:
      jsonpath: age
    masks:
      - add: \"\"
      - randomInt:
          min: 20
          max: 70
  - selector:
      jsonpath: salary
    masks:
      - add: \"\"
      - randomChoice:
          - 1100
          - 1200
          - 1300
          - 1400
          - 1500
          - 1600
          - 1700

" > masking.yml

pimo --empty-input --repeat 1000 > data/originaldata.jsonl

echo "version: v1
masking:
  - selector:
      jsonpath: salary
    mask:
      remove: true
  - selector:
      jsonpath: individual_id
    mask:
      remove: true
" > masking.yml

cat data/originaldata.jsonl | pimo > data/foundData.jsonl

echo "version: v1
masking:
  - selector:
      jsonpath: firstname
    mask:
      randomChoiceInUri: \"pimo://nameFR\"
" > masking.yml

cat data/originaldata.jsonl | pimo > data/nameAnonymized.jsonl

echo "version: v1
masking:
  - selector:
      jsonpath: city
    mask:
      incremental:
        start: 1
        increment: 1
    cache: \"hash\"
caches:
    hash:
        unique: true
" > masking.yml

cat data/nameAnonymized.jsonl | pimo --dump-cache hash=data/hash.jsonl | sigo -q city,age -k 1 -i id > data/aftersigo.jsonl

echo "version: v1
masking:
  - selector:
      jsonpath: temp
    masks:
      - add-transient: \"\"
      - replacement: key
  - selector:
      jsonpath: key
    masks:
      - replacement: \"value\"
  - selector:
      jsonpath: value
    masks:
      - replacement: \"temp\"
" > masking.yml

cat data/hash.jsonl | pimo > data/hash.jsonl

echo "version: v1
masking:
  - selector:
      jsonpath: city
    mask:
      fromCache: hash
caches:
    hash: {}
" > masking.yml

cat data/aftersigo.jsonl | pimo --load-cache hash=data/hash.jsonl > data/aftersigo.jsonl





