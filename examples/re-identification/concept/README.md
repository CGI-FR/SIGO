# Reidentification key concepts

**Definition of distance :** (*wikipédia*)
We call distance, on a set ![equation](https://latex.codecogs.com/svg.image?E), any application `d` defined on ![equation](https://latex.codecogs.com/svg.image?E%5E%7B2%7D) and with values in the set of positive or zero real numbers (![equation](https://latex.codecogs.com/svg.image?%5Cmathbb%7BR%7D&plus;)),

> ![equation](https://latex.codecogs.com/svg.image?d%20:%20E%20%5Ctimes%20E%20%5Cto%20%5Cmathbb%7BR%7D&plus;)

verifying the following properties :

- symmetry : ![equation](https://latex.codecogs.com/svg.image?%5Cforall%20(a,b)%20%5Cin%20E%5E%7B2%7D,%20d(a,b)%20=%20d(b,a)%20)
- separation : ![equation](https://latex.codecogs.com/svg.image?%5Cforall%20(a,b)%20%5Cin%20E%5E%7B2%7D,%20d(a,b)%20=%200%20%5CLeftrightarrow%20a%20=%20b%20)
- triangular inequality : ![equation](https://latex.codecogs.com/svg.image?%5Cforall%20(a,b,c)%20%5Cin%20E%5E%7B3%7D,%20d(a,c)%20%5Cleq%20d(a,b)%20&plus;%20d(b,c))

The best known distances are the **Euclidean distance** and the **Manhattan distance**.

![image](manhattan.png)

|![equation](https://latex.codecogs.com/svg.image?d(a,b)%20=%20%5Csqrt%7B%5Csum_%7Bi=1%7D%5E%7Bn%7D%5Cleft%20(%20a_%7Bi%7D%20-%20b_%7Bi%7D%20%5Cright%20)%5E%7B2%7D%7D%20) | ![equation](https://latex.codecogs.com/svg.image?%5Cinline%20d(a,b)%20=%20%5Csum_%7Bi=1%7D%5E%7Bn%7D%20%5Cleft%7C%20a_%7Bi%7D%20-%20b_%7Bi%7D%20%5Cright%7C) |
|:----------------------:|:---------------------:|
| **Euclidean** distance | **Manhattan** distance|

![image](distances.png)

> https://towardsdatascience.com/17-types-of-similarity-and-dissimilarity-measures-used-in-data-science-3eb914d2681

We can use a distance to define a similarity: the more two points are distant, the less similar they are, and inversely. We can then transform a distance `d` into a similarity `s` in the following way :

> Similarity
> ![equation](https://latex.codecogs.com/svg.image?s%20=%20%5Cfrac%7B1%7D%7B1%20&plus;%20d%7D)

A similarity index is between 0 and 1 where 0 indicates a total difference between the points and 1 indicates a total similarity.

We have the following properties :

![equation](https://latex.codecogs.com/svg.image?%5Cleft%5C%7B%5Cbegin%7Bmatrix%7Dd%20=%200%20%5CLeftrightarrow%20s%20=%201%20%5C%5Cd%20=%20&plus;%5Cinfty%20%5CLeftrightarrow%20s%20=%200%20%20%20%5Cend%7Bmatrix%7D%5Cright.)

![image](sim.png)
