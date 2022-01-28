# Boundary Tests on SIGO

This document presents the limitations of sigo regarding the size of the datasets
We will generate different datasets by varying the number of rows and the number of columns.
We will test these parameters on the different anonymization methods.

## Number of rows

We generate datasets of different sizes using `pimo`.
Below the `masking.yml` file allowing to generate a flow of jsonline with random float.

```yaml
version: "1"
seed: 42
masking:
  - selector:
      jsonpath: "A"
    mask:
      randomDecimal:
        min: 0
        max: 100.00
        precision: 2
  - selector:
      jsonpath: "B"
    mask:
      randomDecimal:
        min: 0
        max: 100.00
        precision: 2
```

We change the size of the datasets by using the `--repeat,-r` flag of `pimo` (***N = [100, 1000, 10000, 100000, 1000000]***).
And we anonymize the data with `sigo` using the anonymization method of our choice with the `--anonymizer,-a` flag.

```console
pimo < test.json -c masking.yml -r 100 > test1_1.json
sigo -q A,B -a general < test1_1.json > output.json
```

A bash script is written to automate the tests in the `rows.sh` file.

```console
cd rows
sudo chmod u+x rows.sh
. ./rows.sh
```

The results are listed in the `log.txt` file.

| NoAnonymizer |    Size    | Execution  time  (sec) | Results |
|-------|:----------:|:----------------------:|:-------:|
| Test1 |     100    |          0.00          | SUCCESS |
| Test2 |    1 000   |          0.00          | SUCCESS |
| Test3 |   10 000   |          2.00          | SUCCESS |
| Test4 |   100 000  |         27.00          | SUCCESS |
| Test5 |  1 000 000 |        418.00          | SUCCESS |
| Test6 | 10 000 000 |                        |  FAILED |

<table>
<tr><th> Generalization </th><th> Aggregation </th><th> Top Bottom Coding </th><th> Random Noise </th></tr>
<tr><td>

|       |    Size    | Execution  time  (sec) | Results |
|-------|:----------:|:----------------------:|:-------:|
| Test1 |     100    |          1.00          | SUCCESS |
| Test2 |    1 000   |          0.00          | SUCCESS |
| Test3 |   10 000   |          3.00          | SUCCESS |
| Test4 |   100 000  |         30.00          | SUCCESS |
| Test5 |  1 000 000 |        395.00          | SUCCESS |

</td><td>

|       |    Size    | Execution  time  (sec) | Results |
|-------|:----------:|:----------------------:|:-------:|
| Test1 |     100    |          0.00          | SUCCESS |
| Test2 |    1 000   |          0.00          | SUCCESS |
| Test3 |   10 000   |          3.00          | SUCCESS |
| Test4 |   100 000  |         29.00          | SUCCESS |
| Test5 |  1 000 000 |        386.00          | SUCCESS |

</td><td>

|       |    Size    | Execution  time  (sec) | Results |
|-------|:----------:|:----------------------:|:-------:|
| Test1 |     100    |          0.00          | SUCCESS |
| Test2 |    1 000   |          0.00          | SUCCESS |
| Test3 |   10 000   |          3.00          | SUCCESS |
| Test4 |   100 000  |         28.00          | SUCCESS |
| Test5 |  1 000 000 |        398.00          | SUCCESS |

</td><td>

|       |    Size    | Execution  time  (sec) | Results |
|-------|:----------:|:----------------------:|:-------:|
| Test1 |     100    |          0.00          | SUCCESS |
| Test2 |    1 000   |          1.00          | SUCCESS |
| Test3 |   10 000   |          3.00          | SUCCESS |
| Test4 |   100 000  |         37.00          | SUCCESS |
| Test5 |  1 000 000 |        420.00          | SUCCESS |

</td></tr> </table>

![rows](rows/rows.png)

## Number of columns

Now we generate a dataset of 1000 rows using `pimo` and we change the number of attributes.
To do this we take the `masking.yml` file from the **rows folder** and add additional masks for each new attribute.
Here an example for the test with 4 attributes,

```yaml
version: "1"
seed: 42
masking:
  - selector:
      jsonpath: "A"
    mask:
      randomDecimal:
        min: 0
        max: 100.00
        precision: 2
  - selector:
      jsonpath: "B"
    mask:
      randomDecimal:
        min: 0
        max: 100.00
        precision: 2
  - selector:
      jsonpath: "C"
    mask:
      randomDecimal:
        min: 0
        max: 100.00
        precision: 2
  - selector:
      jsonpath: "D"
    mask:
      randomDecimal:
        min: 0
        max: 100.00
        precision: 2
```

So we get 5 masking.yml files :

- `masking1.yml` for 2 attributes.
- `masking2.yml` for 4 attributes.
- `masking3.yml` for 8 attributes.
- `masking4.yml` for 16 attributes.
- `masking5.yml` for 32 attributes.

```console
pimo < test2.json -c masking2.yml -r 1000 > test2_2.json
sigo -q A,B,C,D -a general < test2_2.json > output.json
```

The bash script for test automation is in the `columns.sh` file and the results are in the `log.txt` file.

```console
cd columns
sudo chmod u+x columns.sh
. ./columns.sh
```

| NoAnonymizer | Attributes | Execution  time  (sec) | Results |
|-------|:----------:|:----------------------:|:-------:|
| Test1 |      2     |          0.00          | SUCCESS |
| Test2 |      4     |          1.00          | SUCCESS |
| Test3 |      8     |          0.00          | SUCCESS |
| Test4 |     16     |          1.00          | SUCCESS |
| Test5 |     32     |          3.00          | SUCCESS |

<table>
<tr><th> Generalization </th><th> Aggregation </th><th> Top Bottom Coding </th><th> Random Noise </th></tr>
<tr><td>

|       | Attributes | Execution  time  (sec) | Results |
|-------|:----------:|:----------------------:|:-------:|
| Test1 |      2     |          0.00          | SUCCESS |
| Test2 |      4     |          1.00          | SUCCESS |
| Test3 |      8     |          1.00          | SUCCESS |
| Test4 |     16     |          2.00          | SUCCESS |
| Test5 |     32     |          4.00          | SUCCESS |

</td><td>

|       | Attributes | Execution  time  (sec) | Results |
|-------|:----------:|:----------------------:|:-------:|
| Test1 |      2     |          0.00          | SUCCESS |
| Test2 |      4     |          1.00          | SUCCESS |
| Test3 |      8     |          1.00          | SUCCESS |
| Test4 |     16     |          2.00          | SUCCESS |
| Test5 |     32     |          4.00          | SUCCESS |

</td><td>

| Top Bottom Coding | Attributes | Execution  time  (sec) | Results |
|-------|:----------:|:----------------------:|:-------:|
| Test1 |      2     |          0.00          | SUCCESS |
| Test2 |      4     |          0.00          | SUCCESS |
| Test3 |      8     |          0.00          | SUCCESS |
| Test4 |     16     |          1.00          | SUCCESS |
| Test5 |     32     |          4.00          | SUCCESS |

</td><td>

|       | Attributes | Execution  time  (sec) | Results |
|-------|:----------:|:----------------------:|:-------:|
| Test1 |      2     |          0.00          | SUCCESS |
| Test2 |      4     |          0.00          | SUCCESS |
| Test3 |      8     |          0.00          | SUCCESS |
| Test4 |     16     |          2.00          | SUCCESS |
| Test5 |     32     |          4.00          | SUCCESS |

</td></tr> </table>

![columns](columns/columns.png)
