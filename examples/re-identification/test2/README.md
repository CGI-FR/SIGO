# Re-identification test

## 1st example
### Datasets

There are 2 files:

- `arbres.json` : the file on the trees of Paris to be anonymized containing the sensitive data **remarquable**.
- `arbres_openData.json` : a file containing information on the trees of Paris that can be found on the open data.

> arbres.json

|   genre  | circonference | hauteur | remarquable |     x     |     y    |
|:--------:|:-------------:|:-------:|:-----------:|:---------:|:--------:|
|  Prunus  |       76      |    5    |     NON     | 48.889788 | 2.319906 |
|   Tilia  |      111      |    12   |     NON     | 48.894999 | 2.363968 |
| Carpinus |       28      |    7    |     NON     | 48.879623 | 2.363230 |
| Platanus |      110      |    10   |     NON     | 48.839062 | 2.391879 |
|   Tilia  |      105      |    10   |     NON     | 48.832918 | 2.446663 |

> arbres_openData.json

|   genre  | circonference | hauteur |     x     |     y    |
|:--------:|:-------------:|:-------:|:---------:|:--------:|
|  Prunus  |       76      |    5    | 48.889788 | 2.319906 |
|   Tilia  |      111      |    12   | 48.894999 | 2.363968 |
| Carpinus |       28      |    7    | 48.879623 | 2.363230 |
| Platanus |      110      |    10   | 48.839062 | 2.391879 |
|   Tilia  |      105      |    10   | 48.832918 | 2.446663 |

### 1st "bad anonymization"

Forgot to anonymize the **genre** column and use the method `meanAggregation`.

```console
sigo -q x,y,circonference,hauteur -s remarquable -a meanAggregation < arbres.json > arbres-sigo.json
```

> arbres-sigo.json

|        |  genre | circonference | hauteur | remarquable |   x   |  y  |
|--------|:------:|:-------------:|:-------:|:-----------:|:-----:|:---:|
| 322004 | Prunus |      0.0      |   0.0   |     NON     | 48.82 | 2.3 |
| 322003 | Prunus |      0.0      |   0.0   |     NON     | 48.82 | 2.3 |
| 322002 | Prunus |      0.0      |   0.0   |     NON     | 48.82 | 2.3 |
| 322001 | Prunus |      0.0      |   0.0   |     NON     | 48.83 | 2.3 |
| 201022 | Prunus |      0.0      |   0.0   |     NON     | 48.83 | 2.3 |

By grouping data with the same values for the attributes **circonference**, **hauteur**, **x** and **y**, we can re-identify some individuals by linking to the attribute **gender**.

Take for example the cluster formed by the 3 individuals below, the *Punica* is a tree noted as *remarquable*.

|        |     genre    | circonference |   hauteur  | remarquable |      x      |      y     |
|--------|:------------:|:-------------:|:----------:|:-----------:|:-----------:|:----------:|
| 501004 |    Prunus    |     30.67     |    4.33    |     NON     |    48.89    |    2.35    |
| 701006 |    Ostrya    |     30.67     |    4.33    |     NON     |    48.89    |    2.35    |
| 404003 | **_Punica_** |  **_30.67_**  | **_4.33_** |  **_OUI_**  | **_48.89_** | **_2.35_** |

If we look at the data collected from the open data, there are only 3 trees in Paris that are *Punicas*.

|         |     genre    | circonference | hauteur |        x        |        y       |
|---------|:------------:|:-------------:|:-------:|:---------------:|:--------------:|
| 404003  | **_Punica_** |    **_30_**   | **_3_** | **_48.885642_** | **_2.343820_** |
| 250012  | Punica |       0       |    0    | 48.835915 | 2.446839 |
| 101010  | Punica |       5       |    1    | 48.871901 | 2.275000 |

We can easily make the link that the tree `{genre:Punica, circonference:30, hauteur:3, x:48.885642, y:2.343820`} is notes as `remarquable`, so find its sensitive data.

Another example with another group of subjects.

|       |      genre     | circonference |   hauteur  | remarquable |      x      |      y     |
|-------|:--------------:|:-------------:|:----------:|:-----------:|:-----------:|:----------:|
| 60050 | **_Pistacia_** |  **_189.67_** | **_10.0_** |  **_OUI_**  | **_48.85_** | **_2.25_** |
| 30027 |     Quercus    |     189.67    |    10.0    |     NON     |    48.85    |    2.25    |
| 40020 |    Magnolia    |     189.67    |    10.0    |     NON     |    48.85    |    2.25    |

|        |      genre     | circonference |  hauteur |        x        |        y       |
|--------|:--------------:|:-------------:|:--------:|:---------------:|:--------------:|
| 60050  | **_Pistacia_** |   **_171_**   | **_10_** | **_48.845904_** | **_2.253027_** |
| 104001 |    Pistacia    |       50      |     6    |    48.841918    |    2.297990    |

We can easily make the link that the tree `{genre:Pistacia, circonference:171, hauteur:10, x:48.845904, y:2.253027`} is notes as `remarquable`.

## 2nd example

### Datasets

This time we use a sample of Paris trees:

- `trees.json` : the file on the sample of trees in Paris to be anonymized containing the sensitive data **remarquable**.
- `trees-paris.json` : a file containing information on the trees of Paris that can be found on the open data.

> trees.json

| id |  hauteur  | circonference | arrondissement | remarquable |
|----|:---------:|:-------------:|:--------------:|:-----------:|
| 1  | 48.850732 |      2.406460 |        1       |     OUI     |
| 2  | 48.863923 |      2.338329 |        1       |     NON     |
| 3  | 48.830706 |      2.356600 |        3       |     NON     |
| 4  | 48.837150 |      2.436883 |        2       |     OUI     |
| 5  | 48.873035 |      2.274325 |        2       |     NON     |

> trees-paris.json

| id |  hauteur  | arrondissement |
|----|:---------:|:--------------:|
| 1  | 48.850732 |        1       |
| 2  | 48.863923 |        1       |
| 3  | 48.830706 |        3       |
| 4  | 48.837150 |        2       |
| 5  | 48.873035 |        2       |

### 2nd "bad anonymization"

We use the `outlier` method to anonymize the dataset.

```console
sigo -q circonference,hauteur,arrondissement -s remarquable -a outlier < trees.json > trees-sigo.json
```

> trees-sigo.json

| id |  hauteur  | circonference | arrondissement | remarquable |
|----|:---------:|:-------------:|:--------------:|:-----------:|
| 11 | 48.801988 |    2.307882   |        1       |     OUI     |
| 16 | 48.808755 |    2.306808   |        2       |     OUI     |
| 6  | 48.868977 |    2.285416   |        2       |     OUI     |
| 23 | 48.847782 |    2.275808   |        3       |     OUI     |
| 28 | 48.858214 |    2.321236   |        3       |     NON     |

If we compare the anonymized data of **circonference**, **hauteur** and **arrondissement** with the data of **hauteur** and **arrondissement** from the open data, we can make the link with 17 trees (more than half: 17/30).

Here is an example with 3 trees that can be easily identified.

The trees with identifiers ***2***, ***11*** and ***30***, which can be found in the open data, have values of **hauteur** and **arrondissement** identical to the anonymized data.

> open data

| id |  hauteur  | arrondissement |
|----|:---------:|:--------------:|
| 2  | 48.863923 |        1       |
| 11 | 48.801988 |        1       |
| 30 | 48.887968 |        4       |

> anonymized data

| id |  hauteur  | circonference | arrondissement | remarquable |
|----|:---------:|:-------------:|:--------------:|:-----------:|
| 2  | 48.863923 |    2.345846   |        1       |     NON     |
| 11 | 48.801988 |    2.307882   |        1       |     OUI     |
| 30 | 48.887968 |    2.367091   |        4       |     OUI     |

We can therefore deduce that tree ***2*** is not "remarquable" and that trees ***11*** and ***30*** are.

> This anonymization method should be used with care. It is to be used with another method or on only a few columns and not on the full dataset.
