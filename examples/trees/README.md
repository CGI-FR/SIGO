# Tree data test

This example is based on the **trees.json** dataset containing information about trees in Paris.
The dataset is downloaded from Paris Open Data Repository and we have removed some attributes.

So for each tree we have the following data :

- ***Id*** tree id
- ***genre*** type of tree (175 types)
- ***circonference*** circumference of the tree (cm.)
- ***hauteur*** height of the tree (cm)
- ***x*** coordinate x, geolocation of the tree (latitude)
- ***y*** coordinate y, geolocation of the tree (longitude)

Consider that the ***genre*** of the tree is a sensitive data and given the original data,

![original](cars.png)

```console
< cars.json | jq -c '.[]' | sigo -q Id,Miles_per_Gallon,Cylinders,Displacement,Horsepower,Weight_in_lbs,Acceleration | jq -s > cars_sigo.json
```

![masked](cars-sigo.png)

We can see that the anonymisation of the dataset has not changed the correlation of the attributes.
