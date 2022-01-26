# Random Noise

It consists in adding random noise on quantitative variables.

This demo is based on the **cars.json** dataset and the added noise comes from a **Laplace** distribution.

## 2-dimension

This part focuses on 2 quasi-identifiers : ***Miles_per_Gallon*** and ***Horsepower***.
The sensitive data is always the attribute ***Origin***.

The data is recorded in the cars2.json file, below you will find an overview of the data :

```json
    {
      "Miles_per_Gallon": 18,
      "Horsepower": 130,
      "Origin": "USA"
    },
    {
      "Miles_per_Gallon": 15,
      "Horsepower": 165,
      "Origin": "USA"
    },
    {
      "Miles_per_Gallon": 18,
      "Horsepower": 150,
      "Origin": "USA"
    },
```

### 1- Original data

![original](cars2.png)

|                  | Miles_per_Gallon | Horsepower |
|------------------|------------------|------------|
| Miles_per_Gallon |     1.000000     |  -0.778427 |
| Horsepower       |     -0.778427    |  1.000000  |

### 1- De-identified data

```console
< cars2.json | jq -c '.[]' | sigo -q Miles_per_Gallon,Horsepower -s Origin -a laplaceNoise | jq -s > cars2_sigo.json
```

![masked](cars2-sigo.png)

|                  | Miles_per_Gallon | Horsepower |
|------------------|------------------|------------|
| Miles_per_Gallon |     1.000000     |  -0.778478 |
| Horsepower       |     -0.778478    |  1.000000  |

## n-dimension

This part focuses on 6 quasi-identifiers : ***Miles_per_Gallon***, ***Cylinders***, ***Displacement***, ***Horsepower***, ***Weight_in_lbs*** and ***Acceleration***.
The sensitive data is always the attribute ***Origin***.

The data is recorded in the carsn.json file, below you will find an overview of the data :

```json
   {
      "Miles_per_Gallon": 18,
      "Cylinders": 8,
      "Displacement": 307,
      "Horsepower": 130,
      "Weight_in_lbs": 3504,
      "Acceleration": 12,
      "Origin": "USA"
    },
    {
      "Miles_per_Gallon": 15,
      "Cylinders": 8,
      "Displacement": 350,
      "Horsepower": 165,
      "Weight_in_lbs": 3693,
      "Acceleration": 11.5,
      "Origin": "USA"
    },
```

### 2- Original data

![original](carsn.png)

|                  | Miles_per_Gallon | Cylinders | Displacement | Horsepower | Weight_in_lbs | Acceleration |
|------------------|:----------------:|:---------:|:------------:|:----------:|:-------------:|:------------:|
| Miles_per_Gallon |     1.000000     | -0.777618 |   -0.805127  |  -0.778427 |   -0.832244   |   0.423329   |
| Cylinders        |     -0.777618    |  1.000000 |   0.950823   |  0.842983  |    0.897527   |   -0.504683  |
| Displacement     |     -0.805127    |  0.950823 |   1.000000   |  0.897257  |    0.932994   |   -0.543800  |
| Horsepower       |     -0.778427    |  0.842983 |   0.897257   |  1.000000  |    0.864538   |   -0.689196  |
| Weight_in_lbs    |     -0.832244    |  0.897527 |   0.932994   |  0.864538  |    1.000000   |   -0.416839  |
| Acceleration     |     0.423329     | -0.504683 |   -0.543800  |  -0.689196 |   -0.416839   |   1.000000   |

### 2- De-identified data

```console
< carsn.json | jq -c '.[]' | sigo -q Miles_per_Gallon,Cylinders,Displacement,Horsepower,Weight_in_lbs,Acceleration -s Origin -a laplaceNoise | jq -s > carsn_sigo.json
```

![masked](carsn-sigo.png)

|                  | Miles_per_Gallon | Cylinders | Displacement | Horsepower | Weight_in_lbs | Acceleration |
|------------------|:----------------:|:---------:|:------------:|:----------:|:-------------:|:------------:|
| Miles_per_Gallon |     1.000000     | -0.686867 |   -0.673229  |  -0.673453 |   -0.713713   |   0.355517   |
| Cylinders        |     -0.686867    |  1.000000 |   0.629645   |  0.592229  |    0.704464   |   -0.307798  |
| Displacement     |     -0.673229    |  0.629645 |   1.000000   |  0.627392  |    0.670037   |   -0.342196  |
| Horsepower       |     -0.673453    |  0.592229 |   0.627392   |  1.000000  |    0.664508   |   -0.464786  |
| Weight_in_lbs    |     -0.713713    |  0.704464 |   0.670037   |  0.664508  |    1.000000   |   -0.293661  |
| Acceleration     |     0.355517     | -0.307798 |   -0.342196  |  -0.464786 |   -0.293661   |   1.000000   |

The correlation after anonymization is in the range ![equation](https://latex.codecogs.com/svg.image?%5Cinline%20%5Cleft%20%5B%20%5Cpm%200.07%20;%20%5Cpm%200.3%20%5Cright%20%5D)

### Bibliography

***Brand, Ruth.***, **"Microdata Protection through Noise Addition"**,
[in Inference Control in Statistical Databases, From Theory to Practice, 2002, 97‑116](<https://link.springer.com/chapter/10.1007/3-540-47804-3_8?code=d7da801e-b5d7-4f86-8820-3547ba948938>).
