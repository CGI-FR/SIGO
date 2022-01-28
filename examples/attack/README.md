# SIGO USE FOR REIDENTIFICATION

This document presents how to use **SIGO** to test reidentification.

## State of the problem

The data published **after de-identification** by **SIGO**:

```json
{"x":4.67,"y":8.33,"z":"a"}
{"x":4.67,"y":8.33,"z":"b"}
{"z":"c","x":4.67,"y":8.33}
{"x":9.33,"y":7.33,"z":"b"}
{"x":9.33,"y":7.33,"z":"a"}
{"x":9.33,"y":7.33,"z":"b"}
{"z":"b","x":3.67,"y":16.67}
{"y":16.67,"z":"a","x":3.67}
{"y":16.67,"z":"a","x":3.67}
{"x":10,"y":16,"z":"c"}
{"z":"c","x":10,"y":16}
{"x":10,"y":16,"z":"b"}
{"x":10,"y":16,"z":"b"}
{"x":16.4,"y":3.6,"z":"b"}
{"x":16.4,"y":3.6,"z":"c"}
{"x":16.4,"y":3.6,"z":"b"}
{"x":16.4,"y":3.6,"z":"c"}
{"x":16.4,"y":3.6,"z":"c"}
{"x":15.83,"y":14.83,"z":"a"}
{"z":"c","x":15.83,"y":14.83}
{"x":15.83,"y":14.83,"z":"c"}
{"x":15.83,"y":14.83,"z":"b"}
{"x":15.83,"y":14.83,"z":"a"}
{"x":15.83,"y":14.83,"z":"c"}
```

In addition, we have the dataset without the sensitive attribute:

```json
{"x":15,"y":18,"id":1}
{"x":10,"y":20,"id":2}
{"x":6,"y":7,"id":3}
{"x":12,"y":20,"id":4}
{"x":2,"y":19,"id":5}
{"x":18,"y":6,"id":6}
{"x":2,"y":16,"id":7}
{"x":4,"y":9,"id":8}
{"x":18,"y":7,"id":9}
{"x":9,"y":7,"id":10}
{"x":13,"y":0,"id":11}
{"x":17,"y":2,"id":12}
{"x":8,"y":13,"id":13}
{"x":14,"y":14,"id":14}
{"x":12,"y":10,"id":15}
{"x":4,"y":9,"id":16}
{"x":7,"y":5,"id":17}
{"x":18,"y":8,"id":18}
{"x":15,"y":20,"id":19}
{"x":16,"y":3,"id":20}
{"x":10,"y":11,"id":21}
{"x":7,"y":15,"id":22}
{"x":19,"y":20,"id":23}
{"x":14,"y":9,"id":24}

```

The objective is to identifie each quasi-identifier (here associated with an id) to the sensitive attribute

## Some leads

Supposing we know the quasi-identifiers values, we can determine clusters, and the size of them. We can suppose a lower boundary for **k**. Same for **l** (number of different attribute in each cluster).

Now, we don't know the anonymisation methode used. Could be meanAggregation, medianAggregation, TopBottomCoding, or other methode not (yet) implemented in Sigo.

## First Step

The first step is to find clusters. It might be easier for some anonymisation methodes: meanAggregation, medianAggregation or Generalizer display visible cluster in the anonymized dataSet. Might be trickier with randomNoise or topBottom coding for which the clusters are not clearly visible in the output dataset.

@[vega](plot.vg.json)

# Situation example:

Anonymized DataSet:

```json
{"firstname":"Marie-Élise","city":"Lyon","age":40,"salary":1700}
{"firstname":"Wilson","city":"Nantes","age":30,"salary":1200}
{"firstname":"Romane","city":"Paris","age":22,"salary":1100}
{"firstname":"Laurent","city":"Lyon","age":40,"salary":1700}
```

Data retrieved with other sources by attackers:

```json
{"firstname":"Laura","city":"Lyon","age":40}
{"firstname":"Arlette","city":"Nantes","age":30}
{"firstname":"Joseph","city":"Paris","age":22}
{"firstname":"Gaëtan","city":"Lyon","age":40}
```

Supposing we know that only the firstName is anonymized (with **PIMO** for instance), as it is an identifier. (the tuple [city,age] is a quasi-identifier). The objective is to associate each individual to their salary.

With only 4 records, besides ordonned, it's easy to identifie individuals without programs.
However, with thousands of records, not even ordored, it would ask for a lot of time to do it manually.

Let's find out how SIGO would be able to identify a portion of the individuals.

## K-Anonymity

First thing to observe is that the tuple [city,age] is a quasi-identifier, and if k-anonymity is not preserved, some individuals can be easily identified.
In the example above, the individuals with pseudonyme `Wilson` and `Romane` are unique for their set of quasi-identifier. As we can intersect with our data found with other sources, we now know the salary of `Arlette` and `Joseph`.

## L-Diversity

Some individuals will share the same tuple of quasi-identifiers. For instance, we see that the real `Gaëtan` and `Laura` are 40yo and live in Lyon. We say that this pair of records (aka cluster) respects the 2-Anonymity. We imagine that we won't be able to find their salary. Luckily for us, attackers, the data does not respects l-diversity: the two individuals share the same salary. We know that `Gaëtan` and `Laura` have a salary of `1700(ua)`.

## Using SIGO:

In the example above, we were able to identify each indivudal with a confidence of 100%. We can imagine a solution using SIGO that will propose an identification for each records, with the level of confidence.

The objective of this demo is to show that pseudo-anonymisation of data cannot entirely protect your data. Attributes and tables association
