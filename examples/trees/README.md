# Tree data test

This example is based on the **trees.json** dataset containing information about trees in Paris.
The dataset is downloaded from Paris Open Data Repository and we have removed some attributes.

So for each tree we have the following data :

- ***id*** tree id
- ***circonference*** circumference of the tree (cm.)
- ***hauteur*** height of the tree (m.)
- ***remarquable*** whether the tree is remarkable or not (OUI/NON)
- ***x*** coordinate x, geolocation of the tree (latitude)
- ***y*** coordinate y, geolocation of the tree (longitude)

We removed the ***id*** attribute because it is an identifier.
Consider that the ***remarquable*** attribute is a sensitive data and given the original data,

```json
{
    "circonference":20,
    "hauteur":5,
    "remarquable":"NON",
    "x":48.9002546593994,
    "y":2.334152828878867
  },
  {
    "circonference":115,
    "hauteur":10,
    "remarquable":"NON",
    "x":48.84935636396974,
    "y":2.3957233289766773
  },
```

This part is intended to show you different anonymization techniques and to show you that this does not significantly affect the correlation of attributes.

## Correlation

To calculate the correlation between each variable of the dataset we use the pearson correlation.
Pearson correlation measures the strength of the linear relationship between two continuous variables. It has a value between -1 to 1, with a value of -1 meaning a total negative linear correlation, 0 being no correlation, and + 1 meaning a total positive correlation.

Pearson Correlation Coefficient :
$$ \rho \left( x, y \right) = \frac{\sum \left[ \left( x_i - \overline x \right) \times \left( y_i - \overline y \right) \right]}{\sigma_x \times \sigma_y}  $$

With,

$$ \overline x  \text : \space \text mean \space \text of \space \text x \space \text variable. \\
\overline y  \text : \space \text mean \space \text of \space \text y \space \text variable. \\
\sigma_x  \text : \space \text standart \space \text deviation \space \text of \space \text x \space \text variable. \\
\sigma_y  \text : \space \text standart \space \text deviation \space \text of \space \text y \space \text variable. $$

```python
import pandas as pd
import numpy as np
import json

input_file = open(r'trees.json')
jsondata = json.load(input_file)
df = pd.DataFrame(jsondata)

df.corr(method='pearson')
```

|               | circonference |  hauteur  |     x     |     y    |
|---------------|:-------------:|:---------:|:---------:|:--------:|
| circonference |    1.000000   |  0.848523 | -0.045860 | 0.017326 |
| hauteur       |    0.848523   |  1.000000 | -0.032621 | 0.168414 |
| x             |   -0.045860   | -0.032621 |  1.000000 | 0.001270 |
| y             |    0.017326   |  0.168414 |  0.001270 | 1.000000 |

```python
import seaborn as sb

mask = np.triu(np.ones(df.corr().shape)).astype(np.bool)
plot = sb.heatmap(df.corr(), mask = mask, cmap="YlGnBu", annot=True)
```

![correlation](correlation_trees.png)
