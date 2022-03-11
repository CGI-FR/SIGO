# Re-identification issue

With the evolution of information technologies that make it possible to link data from different sources, it is almost impossible to guarantee an anonymization that would offer a zero risk of re-identification.

**Re-identification Definition :** A process (or algorithm) that takes an anonymized dataset and related knowledge as input and seeks to match the anonymized data with real-world individuals.

Let's take as an example a very simple dataset that you can find in the `data.json` file.

```json
{"x": 5, "y": 6, "z":"a"}
{"x": 3, "y": 7, "z":"a"}
{"x": 4, "y": 4, "z":"c"}
{"x": 2, "y": 10, "z":"b"}
{"x": 8, "y": 4, "z":"a"}
...
```

And suppose that we have 2 quasi-identifiers: `x` and `y` and as sensitive data the variable `z`. Anonymize the dataset using sigo, we use **k=4** and **l=3** with the **meanAggregation** method :

```console
< data.json | sigo -q x,y -s z -k 4 -l 3 -a meanAggregation > data-sigo.json
```

```json
{"x":5,"y":6.83,"z":"c"}
{"x":5,"y":6.83,"z":"a"}
{"x":5,"y":6.83,"z":"a"}
{"x":5,"y":6.83,"z":"a"}
{"x":5,"y":6.83,"z":"b"}
...
```

**Objective :** Identify for each individual of the original dataset if an anonymized individual is similar to him.

[define similarity]

Let's assume in the worst case that the attacker has the original dataset but does not have the sensitive data.

```python
14, 6
           x_sigo       sim
6   14.17, 6.0, a  0.999991
7   14.17, 6.0, c  0.999991
8   14.17, 6.0, c  0.999991
9   14.17, 6.0, c  0.999991
10  14.17, 6.0, b  0.999991
11  14.17, 6.0, b  0.999991


7, 19
           x_sigo       sim
0  6.17, 16.67, c  0.999999
1  6.17, 16.67, c  0.999999
2  6.17, 16.67, a  0.999999
3  6.17, 16.67, a  0.999999
4  6.17, 16.67, b  0.999999
5  6.17, 16.67, a  0.999999
```

We can see that the individual `{'x':14, 'y':6}` is similar to a cluster but it is impossible for the attacker to define its sensitive data and it is the same for the individual `{'x':7, 'y':19}`.

Let's try with another anonymization using as parameter **k=3** and **l=1**.

```console
< data.json | sigo -q x,y -s z -k 3 -l 1 -a meanAggregation > data2-sigo.json
```

```json
{"x":3,"y":7,"z":"b"}
{"x":3,"y":7,"z":"a"}
{"x":3,"y":7,"z":"c"}
{"x":7,"y":6.67,"z":"a"}
{"x":7,"y":6.67,"z":"a"}
...
```

```python
20, 18
            x_sigo  sim
3  19.67, 17.67, b  1.0
4  19.67, 17.67, b  1.0
5  19.67, 17.67, b  1.0


3, 7
        x_sigo  sim
0  3.0, 7.0, b  1.0
1  3.0, 7.0, a  1.0
2  3.0, 7.0, c  1.0
```

We can see that in this case the attacker can deduce the sensitive data for the individual `{'x':20, 'y':18}`.

This shows the importance of **k-anonymity** and **l-diversity** but also of the choice of parameters.

**TO DO**

- Choose the most appropriate distance/similarity metric
- Define similarity in README
- Add notebook
