# Tests

```console
sigo -q x,y -s z --load-original "examples/re-identification/test1/data.json" < examples/re-identification/test1/data2-sigo.json
```

## Test1

### Dataset original

```json
{"x":2,"y":10,"z":"b","id":1}
{"id":1,"x":3,"y":7,"z":"a"}
{"id":1,"x":4,"y":4,"z":"c"}
{"z":"a","id":2,"x":5,"y":6}
{"x":8,"y":4,"z":"a","id":2}
{"id":2,"x":8,"y":10,"z":"a"}
{"x":3,"y":16,"z":"a","id":3}
{"x":4,"y":19,"z":"b","id":3}
{"x":6,"y":18,"z":"a","id":3}
{"id":4,"x":7,"y":14,"z":"c"}
{"x":7,"y":19,"z":"a","id":4}
{"z":"c","id":4,"x":10,"y":14}
{"x":11,"y":9,"z":"b","id":5}
{"x":12,"y":3,"z":"a","id":5}
{"x":14,"y":6,"z":"c","id":5}
{"x":15,"y":5,"z":"c","id":6}
{"z":"b","id":6,"x":15,"y":7}
{"id":6,"x":18,"y":6,"z":"c"}
{"x":14,"y":18,"z":"b","id":7}
{"x":18,"y":18,"z":"c","id":7}
{"z":"c","id":7,"x":18,"y":19}
{"z":"b","id":8,"x":19,"y":15}
{"z":"b","id":8,"x":20,"y":18}
{"x":20,"y":20,"z":"b","id":8}
```

### Dataset anonymisé

```json
{"x":3,"y":7,"z":"b","id":1}
{"x":3,"y":7,"z":"a","id":1}
{"x":3,"y":7,"z":"c","id":1}
{"id":2,"x":7,"y":6.67,"z":"a"}
{"x":7,"y":6.67,"z":"a","id":2}
{"y":6.67,"z":"a","id":2,"x":7}
{"z":"a","id":3,"x":4.33,"y":17.67}
{"id":3,"x":4.33,"y":17.67,"z":"b"}
{"x":4.33,"y":17.67,"z":"a","id":3}
{"x":8,"y":15.67,"z":"c","id":4}
{"x":8,"y":15.67,"z":"a","id":4}
{"x":8,"y":15.67,"z":"c","id":4}
{"x":12.33,"y":6,"z":"b","id":5}
{"id":5,"x":12.33,"y":6,"z":"a"}
{"x":12.33,"y":6,"z":"c","id":5}
{"x":16,"y":6,"z":"c","id":6}
{"x":16,"y":6,"z":"b","id":6}
{"y":6,"z":"c","id":6,"x":16}
{"x":16.67,"y":18.33,"z":"b","id":7}
{"x":16.67,"y":18.33,"z":"c","id":7}
{"x":16.67,"y":18.33,"z":"c","id":7}
{"x":19.67,"y":17.67,"z":"b","id":8}
{"x":19.67,"y":17.67,"z":"b","id":8}
{"z":"b","id":8,"x":19.67,"y":17.67}
```

On devrait pouvoir ré-identifier le groupe `"id":2` et `"id":8`, donc 6 individus au total

```json
{"z":"a","id":2,"x":5,"y":6}
{"x":8,"y":4,"z":"a","id":2}
{"id":2,"x":8,"y":10,"z":"a"}

{"z":"b","id":8,"x":19,"y":15}
{"z":"b","id":8,"x":20,"y":18}
{"x":20,"y":20,"z":"b","id":8}
```

#### Cosine_Similarity

```console
reidentification.Reidenfication(original, sigo, "cosine", 3)

[map[sensitive:[a] x:4 y:4]
map[sensitive:[b] x:11 y:9]
map[sensitive:[a] x:20 y:20]
map[sensitive:[b] x:20 y:18] VP
map[sensitive:[a] x:18 y:18]
map[sensitive:[b] x:19 y:15]] VP
```

- VP = 2
- FP = 4

Dans les FP :

    - 2 individus qui ne peuvent pas être identifié avec une valeur de donnée sensible fausse
    - 1 individus est bien identifié mais mauvaise valeur donnée sensible
    - 1 indivudus qui ne peut pas être identifié mais avec une valeur de donnée sensible exacte

#### Distance_Euclidienne

```console
reidentification.Reidenfication(original, sigo, "euclidean", 3)

[map[sensitive:[a] x:5 y:6]
map[sensitive:[a] x:8 y:4]
map[sensitive:[a] x:8 y:10]
map[sensitive:[b] x:20 y:20]
map[sensitive:[b] x:20 y:18]
map[sensitive:[b] x:19 y:15]]
```

- VP = 6

## Test2

### Dataset2 original

```json
{"location":1,"weight":56,"height":1.59,"sickness":"flu","id":1,"age":18}
{"weight":41,"height":1.58,"sickness":"cancer","id":1,"age":20,"location":1}
{"age":36,"location":1,"weight":69,"height":1.48,"sickness":"heart disease","id":1}
{"age":39,"location":1,"weight":62,"height":1.59,"sickness":"flu","id":2}
{"age":43,"location":1,"weight":58,"height":1.72,"sickness":"cancer","id":2}
{"age":52,"location":1,"weight":52,"height":1.47,"sickness":"cancer","id":2}
{"id":3,"age":19,"location":1,"weight":46,"height":1.95,"sickness":"flu"}
{"age":19,"location":1,"weight":56,"height":2.03,"sickness":"cancer","id":3}
{"weight":41,"height":1.96,"sickness":"flu","id":3,"age":42,"location":1}
{"age":52,"location":1,"weight":40,"height":1.99,"sickness":"flu","id":4}
{"height":2.01,"sickness":"heart disease","id":4,"age":56,"location":1,"weight":60}
{"sickness":"heart disease","id":4,"age":60,"location":1,"weight":40,"height":1.99}
{"location":1,"weight":53,"height":1.74,"sickness":"cancer","id":4,"age":61}
{"height":1.47,"sickness":"heart disease","id":5,"age":31,"location":1,"weight":77}
{"height":1.62,"sickness":"flu","id":5,"age":42,"location":1,"weight":102}
{"age":42,"location":1,"weight":85,"height":1.75,"sickness":"flu","id":5}
{"sickness":"heart disease","id":6,"age":47,"location":1,"weight":103,"height":1.56}
{"age":48,"location":1,"weight":90,"height":1.86,"sickness":"cancer","id":6}
{"id":6,"age":52,"location":1,"weight":81,"height":1.85,"sickness":"cancer"}
{"height":1.71,"sickness":"cancer","id":6,"age":53,"location":1,"weight":70}
{"age":19,"location":1,"weight":107,"height":2,"sickness":"heart disease","id":7}
{"age":22,"location":1,"weight":104,"height":2.09,"sickness":"cancer","id":7}
{"weight":77,"height":2.03,"sickness":"flu","id":7,"age":25,"location":1}
{"height":1.97,"sickness":"cancer","id":8,"age":26,"location":1,"weight":110}
{"id":8,"age":50,"location":1,"weight":80,"height":1.99,"sickness":"cancer"}
{"id":8,"age":51,"location":1,"weight":100,"height":1.97,"sickness":"cancer"}
{"weight":105,"height":2.08,"sickness":"heart disease","id":8,"age":54,"location":1}
{"id":9,"age":21,"location":4,"weight":77,"height":1.46,"sickness":"cancer"}
{"age":23,"location":2,"weight":79,"height":1.67,"sickness":"heart disease","id":9}
{"age":31,"location":2,"weight":42,"height":1.46,"sickness":"heart disease","id":9}
{"age":41,"location":2,"weight":55,"height":1.69,"sickness":"cancer","id":10}
{"id":10,"age":46,"location":2,"weight":42,"height":1.87,"sickness":"cancer"}
{"age":61,"location":4,"weight":80,"height":1.67,"sickness":"heart disease","id":10}
{"height":2.01,"sickness":"cancer","id":11,"age":26,"location":3,"weight":81}
{"age":43,"location":4,"weight":80,"height":1.88,"sickness":"flu","id":11}
{"location":4,"weight":46,"height":2.04,"sickness":"cancer","id":11,"age":43}
{"weight":50,"height":1.9,"sickness":"heart disease","id":12,"age":48,"location":2}
{"location":2,"weight":62,"height":1.89,"sickness":"flu","id":12,"age":51}
{"age":56,"location":3,"weight":55,"height":2.04,"sickness":"cancer","id":12}
{"age":19,"location":2,"weight":87,"height":1.67,"sickness":"cancer","id":13}
{"age":26,"location":3,"weight":99,"height":1.53,"sickness":"cancer","id":13}
{"age":44,"location":4,"weight":94,"height":1.57,"sickness":"flu","id":13}
{"age":50,"location":2,"weight":109,"height":1.8,"sickness":"flu","id":14}
{"age":52,"location":2,"weight":84,"height":1.82,"sickness":"heart disease","id":14}
{"sickness":"heart disease","id":14,"age":56,"location":3,"weight":90,"height":1.64}
{"location":2,"weight":102,"height":1.87,"sickness":"flu","id":15,"age":20}
{"age":22,"location":3,"weight":90,"height":1.87,"sickness":"flu","id":15}
{"weight":99,"height":1.97,"sickness":"flu","id":15,"age":27,"location":4}
{"sickness":"flu","id":16,"age":38,"location":3,"weight":104,"height":1.88}
{"sickness":"flu","id":16,"age":53,"location":4,"weight":90,"height":2.06}
{"age":58,"location":2,"weight":107,"height":1.85,"sickness":"flu","id":16}
{"age":65,"location":1,"weight":73,"height":1.66,"sickness":"heart disease","id":17}
{"age":75,"location":2,"weight":60,"height":1.46,"sickness":"heart disease","id":17}
{"id":17,"age":77,"location":2,"weight":43,"height":1.62,"sickness":"cancer"}
{"age":81,"location":1,"weight":50,"height":1.73,"sickness":"flu","id":17}
{"location":1,"weight":71,"height":1.57,"sickness":"heart disease","id":17,"age":87}
{"age":88,"location":2,"weight":46,"height":1.75,"sickness":"flu","id":18}
{"id":18,"age":92,"location":1,"weight":79,"height":1.75,"sickness":"heart disease"}
{"age":93,"location":2,"weight":62,"height":1.73,"sickness":"cancer","id":18}
{"weight":62,"height":1.47,"sickness":"heart disease","id":18,"age":94,"location":2}
{"weight":59,"height":1.46,"sickness":"heart disease","id":18,"age":98,"location":1}
{"age":64,"location":2,"weight":53,"height":1.85,"sickness":"flu","id":19}
{"location":1,"weight":76,"height":2.08,"sickness":"cancer","id":19,"age":74}
{"location":2,"weight":55,"height":1.82,"sickness":"flu","id":19,"age":78}
{"id":19,"age":78,"location":1,"weight":57,"height":2.02,"sickness":"heart disease"}
{"height":1.94,"sickness":"cancer","id":20,"age":80,"location":1,"weight":70}
{"height":1.77,"sickness":"heart disease","id":20,"age":84,"location":2,"weight":61}
{"weight":67,"height":1.9,"sickness":"heart disease","id":20,"age":88,"location":1}
{"age":93,"location":1,"weight":79,"height":1.81,"sickness":"flu","id":20}
{"sickness":"cancer","id":20,"age":99,"location":1,"weight":49,"height":1.87}
{"id":21,"age":62,"location":1,"weight":86,"height":1.6,"sickness":"cancer"}
{"id":21,"age":64,"location":2,"weight":82,"height":1.48,"sickness":"flu"}
{"location":1,"weight":83,"height":1.59,"sickness":"heart disease","id":21,"age":64}
{"height":1.82,"sickness":"flu","id":21,"age":70,"location":1,"weight":104}
{"age":81,"location":2,"weight":89,"height":1.54,"sickness":"heart disease","id":22}
{"sickness":"cancer","id":22,"age":81,"location":2,"weight":108,"height":1.65}
{"location":2,"weight":87,"height":1.48,"sickness":"flu","id":22,"age":98}
{"height":1.48,"sickness":"heart disease","id":22,"age":98,"location":1,"weight":91}
{"age":98,"location":2,"weight":94,"height":1.8,"sickness":"flu","id":22}
{"age":62,"location":1,"weight":106,"height":1.86,"sickness":"cancer","id":23}
{"age":64,"location":2,"weight":95,"height":2.03,"sickness":"heart disease","id":23}
{"location":1,"weight":102,"height":1.88,"sickness":"flu","id":23,"age":67}
{"age":69,"location":2,"weight":83,"height":1.92,"sickness":"heart disease","id":23}
{"height":1.93,"sickness":"cancer","id":24,"age":84,"location":2,"weight":106}
{"weight":103,"height":1.96,"sickness":"heart disease","id":24,"age":90,"location":1}
{"age":90,"location":1,"weight":97,"height":2.07,"sickness":"heart disease","id":24}
{"age":94,"location":2,"weight":109,"height":2.02,"sickness":"cancer","id":24}
{"height":1.98,"sickness":"cancer","id":24,"age":98,"location":1,"weight":100}
{"age":97,"location":3,"weight":73,"height":1.6,"sickness":"flu","id":25}
{"age":88,"location":4,"weight":65,"height":1.83,"sickness":"flu","id":25}
{"height":1.95,"sickness":"heart disease","id":25,"age":62,"location":3,"weight":53}
{"id":26,"age":68,"location":4,"weight":53,"height":2,"sickness":"heart disease"}
{"location":3,"weight":83,"height":2.05,"sickness":"heart disease","id":26,"age":84}
{"age":85,"location":4,"weight":49,"height":2.08,"sickness":"cancer","id":26}
{"id":27,"age":77,"location":3,"weight":86,"height":1.5,"sickness":"heart disease"}
{"sickness":"cancer","id":27,"age":76,"location":3,"weight":94,"height":1.57}
{"weight":104,"height":1.62,"sickness":"cancer","id":27,"age":65,"location":3}
{"height":1.84,"sickness":"cancer","id":28,"age":67,"location":4,"weight":84}
{"weight":101,"height":1.88,"sickness":"cancer","id":28,"age":92,"location":3}
{"age":97,"location":3,"weight":92,"height":1.95,"sickness":"cancer","id":28}
```

### Dataset2 anonymisé

```json
{"height":1.55,"sickness":"flu","id":1,"age":24.67,"location":1,"weight":55.33}
{"age":24.67,"location":1,"weight":55.33,"height":1.55,"sickness":"cancer","id":1}
{"age":24.67,"location":1,"weight":55.33,"height":1.55,"sickness":"heart disease","id":1}
{"id":2,"age":44.67,"location":1,"weight":57.33,"height":1.59,"sickness":"flu"}
{"height":1.59,"sickness":"cancer","id":2,"age":44.67,"location":1,"weight":57.33}
{"location":1,"weight":57.33,"height":1.59,"sickness":"cancer","id":2,"age":44.67}
{"age":26.67,"location":1,"weight":47.67,"height":1.98,"sickness":"flu","id":3}
{"age":26.67,"location":1,"weight":47.67,"height":1.98,"sickness":"cancer","id":3}
{"height":1.98,"sickness":"flu","id":3,"age":26.67,"location":1,"weight":47.67}
{"weight":48.25,"height":1.93,"sickness":"flu","id":4,"age":57.25,"location":1}
{"age":57.25,"location":1,"weight":48.25,"height":1.93,"sickness":"heart disease","id":4}
{"age":57.25,"location":1,"weight":48.25,"height":1.93,"sickness":"heart disease","id":4}
{"sickness":"cancer","id":4,"age":57.25,"location":1,"weight":48.25,"height":1.93}
{"id":5,"age":38.33,"location":1,"weight":88,"height":1.61,"sickness":"heart disease"}
{"age":38.33,"location":1,"weight":88,"height":1.61,"sickness":"flu","id":5}
{"weight":88,"height":1.61,"sickness":"flu","id":5,"age":38.33,"location":1}
{"sickness":"heart disease","id":6,"age":50,"location":1,"weight":86,"height":1.75}
{"location":1,"weight":86,"height":1.75,"sickness":"cancer","id":6,"age":50}
{"height":1.75,"sickness":"cancer","id":6,"age":50,"location":1,"weight":86}
{"sickness":"cancer","id":6,"age":50,"location":1,"weight":86,"height":1.75}
{"age":22,"location":1,"weight":96,"height":2.04,"sickness":"heart disease","id":7}
{"id":7,"age":22,"location":1,"weight":96,"height":2.04,"sickness":"cancer"}
{"age":22,"location":1,"weight":96,"height":2.04,"sickness":"flu","id":7}
{"location":1,"weight":98.75,"height":2,"sickness":"cancer","id":8,"age":45.25}
{"age":45.25,"location":1,"weight":98.75,"height":2,"sickness":"cancer","id":8}
{"location":1,"weight":98.75,"height":2,"sickness":"cancer","id":8,"age":45.25}
{"height":2,"sickness":"heart disease","id":8,"age":45.25,"location":1,"weight":98.75}
{"age":25,"location":2.67,"weight":66,"height":1.53,"sickness":"cancer","id":9}
{"id":9,"age":25,"location":2.67,"weight":66,"height":1.53,"sickness":"heart disease"}
{"location":2.67,"weight":66,"height":1.53,"sickness":"heart disease","id":9,"age":25}
{"age":49.33,"location":2.67,"weight":59,"height":1.74,"sickness":"cancer","id":10}
{"height":1.74,"sickness":"cancer","id":10,"age":49.33,"location":2.67,"weight":59}
{"age":49.33,"location":2.67,"weight":59,"height":1.74,"sickness":"heart disease","id":10}
{"age":37.33,"location":3.67,"weight":69,"height":1.98,"sickness":"cancer","id":11}
{"age":37.33,"location":3.67,"weight":69,"height":1.98,"sickness":"flu","id":11}
{"age":37.33,"location":3.67,"weight":69,"height":1.98,"sickness":"cancer","id":11}
{"height":1.94,"sickness":"heart disease","id":12,"age":51.67,"location":2.33,"weight":55.67}
{"id":12,"age":51.67,"location":2.33,"weight":55.67,"height":1.94,"sickness":"flu"}
{"age":51.67,"location":2.33,"weight":55.67,"height":1.94,"sickness":"cancer","id":12}
{"location":3,"weight":93.33,"height":1.59,"sickness":"cancer","id":13,"age":29.67}
{"height":1.59,"sickness":"cancer","id":13,"age":29.67,"location":3,"weight":93.33}
{"weight":93.33,"height":1.59,"sickness":"flu","id":13,"age":29.67,"location":3}
{"height":1.75,"sickness":"flu","id":14,"age":52.67,"location":2.33,"weight":94.33}
{"location":2.33,"weight":94.33,"height":1.75,"sickness":"heart disease","id":14,"age":52.67}
{"weight":94.33,"height":1.75,"sickness":"heart disease","id":14,"age":52.67,"location":2.33}
{"height":1.9,"sickness":"flu","id":15,"age":23,"location":3,"weight":97}
{"height":1.9,"sickness":"flu","id":15,"age":23,"location":3,"weight":97}
{"location":3,"weight":97,"height":1.9,"sickness":"flu","id":15,"age":23}
{"age":49.67,"location":3,"weight":100.33,"height":1.93,"sickness":"flu","id":16}
{"weight":100.33,"height":1.93,"sickness":"flu","id":16,"age":49.67,"location":3}
{"age":49.67,"location":3,"weight":100.33,"height":1.93,"sickness":"flu","id":16}
{"sickness":"heart disease","id":17,"age":77,"location":1.4,"weight":59.4,"height":1.61}
{"age":77,"location":1.4,"weight":59.4,"height":1.61,"sickness":"heart disease","id":17}
{"sickness":"cancer","id":17,"age":77,"location":1.4,"weight":59.4,"height":1.61}
{"location":1.4,"weight":59.4,"height":1.61,"sickness":"flu","id":17,"age":77}
{"age":77,"location":1.4,"weight":59.4,"height":1.61,"sickness":"heart disease","id":17}
{"location":1.6,"weight":61.6,"height":1.63,"sickness":"flu","id":18,"age":93}
{"sickness":"heart disease","id":18,"age":93,"location":1.6,"weight":61.6,"height":1.63}
{"height":1.63,"sickness":"cancer","id":18,"age":93,"location":1.6,"weight":61.6}
{"weight":61.6,"height":1.63,"sickness":"heart disease","id":18,"age":93,"location":1.6}
{"id":18,"age":93,"location":1.6,"weight":61.6,"height":1.63,"sickness":"heart disease"}
{"sickness":"flu","id":19,"age":73.5,"location":1.5,"weight":60.25,"height":1.94}
{"id":19,"age":73.5,"location":1.5,"weight":60.25,"height":1.94,"sickness":"cancer"}
{"age":73.5,"location":1.5,"weight":60.25,"height":1.94,"sickness":"flu","id":19}
{"age":73.5,"location":1.5,"weight":60.25,"height":1.94,"sickness":"heart disease","id":19}
{"age":88.8,"location":1.2,"weight":65.2,"height":1.86,"sickness":"cancer","id":20}
{"age":88.8,"location":1.2,"weight":65.2,"height":1.86,"sickness":"heart disease","id":20}
{"id":20,"age":88.8,"location":1.2,"weight":65.2,"height":1.86,"sickness":"heart disease"}
{"sickness":"flu","id":20,"age":88.8,"location":1.2,"weight":65.2,"height":1.86}
{"height":1.86,"sickness":"cancer","id":20,"age":88.8,"location":1.2,"weight":65.2}
{"age":65,"location":1.25,"weight":88.75,"height":1.62,"sickness":"cancer","id":21}
{"sickness":"flu","id":21,"age":65,"location":1.25,"weight":88.75,"height":1.62}
{"height":1.62,"sickness":"heart disease","id":21,"age":65,"location":1.25,"weight":88.75}
{"weight":88.75,"height":1.62,"sickness":"flu","id":21,"age":65,"location":1.25}
{"weight":93.8,"height":1.59,"sickness":"heart disease","id":22,"age":91.2,"location":1.8}
{"age":91.2,"location":1.8,"weight":93.8,"height":1.59,"sickness":"cancer","id":22}
{"location":1.8,"weight":93.8,"height":1.59,"sickness":"flu","id":22,"age":91.2}
{"sickness":"heart disease","id":22,"age":91.2,"location":1.8,"weight":93.8,"height":1.59}
{"id":22,"age":91.2,"location":1.8,"weight":93.8,"height":1.59,"sickness":"flu"}
{"age":65.5,"location":1.5,"weight":96.5,"height":1.92,"sickness":"cancer","id":23}
{"weight":96.5,"height":1.92,"sickness":"heart disease","id":23,"age":65.5,"location":1.5}
{"id":23,"age":65.5,"location":1.5,"weight":96.5,"height":1.92,"sickness":"flu"}
{"weight":96.5,"height":1.92,"sickness":"heart disease","id":23,"age":65.5,"location":1.5}
{"weight":103,"height":1.99,"sickness":"cancer","id":24,"age":91.2,"location":1.4}
{"age":91.2,"location":1.4,"weight":103,"height":1.99,"sickness":"heart disease","id":24}
{"weight":103,"height":1.99,"sickness":"heart disease","id":24,"age":91.2,"location":1.4}
{"height":1.99,"sickness":"cancer","id":24,"age":91.2,"location":1.4,"weight":103}
{"age":91.2,"location":1.4,"weight":103,"height":1.99,"sickness":"cancer","id":24}
{"weight":63.67,"height":1.79,"sickness":"flu","id":25,"age":82.33,"location":3.33}
{"weight":63.67,"height":1.79,"sickness":"flu","id":25,"age":82.33,"location":3.33}
{"age":82.33,"location":3.33,"weight":63.67,"height":1.79,"sickness":"heart disease","id":25}
{"age":79,"location":3.67,"weight":61.67,"height":2.04,"sickness":"heart disease","id":26}
{"weight":61.67,"height":2.04,"sickness":"heart disease","id":26,"age":79,"location":3.67}
{"location":3.67,"weight":61.67,"height":2.04,"sickness":"cancer","id":26,"age":79}
{"location":3,"weight":94.67,"height":1.56,"sickness":"heart disease","id":27,"age":72.67}
{"age":72.67,"location":3,"weight":94.67,"height":1.56,"sickness":"cancer","id":27}
{"weight":94.67,"height":1.56,"sickness":"cancer","id":27,"age":72.67,"location":3}
{"age":85.33,"location":3.33,"weight":92.33,"height":1.89,"sickness":"cancer","id":28}
{"age":85.33,"location":3.33,"weight":92.33,"height":1.89,"sickness":"cancer","id":28}
{"sickness":"cancer","id":28,"age":85.33,"location":3.33,"weight":92.33,"height":1.89}
```

On devrait pouvoir ré-identifier le groupe `"id":15`, `"id":16` et `"id":28`, donc 6 individus au total

```json
{"location":2,"weight":102,"height":1.87,"sickness":"flu","id":15,"age":20}
{"age":22,"location":3,"weight":90,"height":1.87,"sickness":"flu","id":15}
{"weight":99,"height":1.97,"sickness":"flu","id":15,"age":27,"location":4}

{"sickness":"flu","id":16,"age":38,"location":3,"weight":104,"height":1.88}
{"sickness":"flu","id":16,"age":53,"location":4,"weight":90,"height":2.06}
{"age":58,"location":2,"weight":107,"height":1.85,"sickness":"flu","id":16}

{"height":1.84,"sickness":"cancer","id":28,"age":67,"location":4,"weight":84}
{"weight":101,"height":1.88,"sickness":"cancer","id":28,"age":92,"location":3}
{"age":97,"location":3,"weight":92,"height":1.95,"sickness":"cancer","id":28}
```

#### Cosine_Similarity2

```console
reidentification.Reidenfication(original, sigo, "cosine", 3)

[map[age:56 height:2.01 location:1 sensitive:[cancer] weight:60]
map[age:44 height:1.57 location:4 sensitive:[flu] weight:94]
map[age:54 height:2.08 location:1 sensitive:[flu] weight:105]
map[age:60 height:1.99 location:1 sensitive:[heart disease] weight:40]
map[age:42 height:1.75 location:1 sensitive:[flu] weight:85]
map[age:65 height:1.66 location:1 sensitive:[cancer] weight:73]
map[age:81 height:1.54 location:2 sensitive:[cancer] weight:89]
map[age:22 height:1.87 location:3 sensitive:[flu] weight:90]
map[age:90 height:2.07 location:1 sensitive:[cancer] weight:97]
map[age:88 height:1.75 location:2 sensitive:[heart disease] weight:46]
map[age:51 height:1.97 location:1 sensitive:[flu] weight:100]
map[age:26 height:1.53 location:3 sensitive:[flu] weight:99]
map[age:21 height:1.46 location:4 sensitive:[flu] weight:77]
map[age:50 height:1.8 location:2 sensitive:[cancer] weight:109]
map[age:92 height:1.88 location:3 sensitive:[cancer] weight:101]
map[age:20 height:1.58 location:1 sensitive:[flu] weight:41]
map[age:27 height:1.97 location:4 sensitive:[flu] weight:99]
map[age:77 height:1.5 location:3 sensitive:[cancer] weight:86]
map[age:36 height:1.48 location:1 sensitive:[flu] weight:69]]
```

- Bien identifié : 3

```json
{"weight":99,"height":1.97,"sickness":"flu","id":15,"age":27,"location":4}
{"age":22,"location":3,"weight":90,"height":1.87,"sickness":"flu","id":15}
{"weight":101,"height":1.88,"sickness":"cancer","id":28,"age":92,"location":3}
```

- Identifié mais mauvaise valeur donnée sensible : 0
- Pas identifié : 6

```json
{"location":2,"weight":102,"height":1.87,"sickness":"flu","id":15,"age":20}
{"age":58,"location":2,"weight":107,"height":1.85,"sickness":"flu","id":16}
{"sickness":"flu","id":16,"age":38,"location":3,"weight":104,"height":1.88}
{"sickness":"flu","id":16,"age":53,"location":4,"weight":90,"height":2.06}
{"height":1.84,"sickness":"cancer","id":28,"age":67,"location":4,"weight":84}
{"age":97,"location":3,"weight":92,"height":1.95,"sickness":"cancer","id":28}
```

- Mauvaise identification : 10

#### Distance_Euclidienne2

```console
reidentification.Reidenfication(original, sigo, "euclidean", 3)

[map[age:44 height:1.57 location:4 sensitive:[cancer] weight:94]
map[age:42 height:1.62 location:1 sensitive:[cancer] weight:102]
map[age:54 height:2.08 location:1 sensitive:[flu] weight:105]
map[age:19 height:2 location:1 sensitive:[flu] weight:107]
map[age:81 height:1.54 location:2 sensitive:[cancer] weight:89]
map[age:38 height:1.88 location:3 sensitive:[cancer] weight:104]
map[age:20 height:1.87 location:2 sensitive:[flu] weight:102]
map[age:22 height:2.09 location:1 sensitive:[flu] weight:104]
map[age:26 height:1.97 location:1 sensitive:[flu] weight:110]
map[age:47 height:1.56 location:1 sensitive:[flu] weight:103]
map[age:58 height:1.85 location:2 sensitive:[flu] weight:107]
map[age:51 height:1.97 location:1 sensitive:[flu] weight:100]
map[age:26 height:1.53 location:3 sensitive:[flu] weight:99]
map[age:50 height:1.8 location:2 sensitive:[flu] weight:109]
map[age:27 height:1.97 location:4 sensitive:[flu] weight:99]
map[age:84 height:2.05 location:3 sensitive:[cancer] weight:83]]
```

- Bien identifié : 3

```json
{"location":2,"weight":102,"height":1.87,"sickness":"flu","id":15,"age":20}
{"weight":99,"height":1.97,"sickness":"flu","id":15,"age":27,"location":4}
{"age":58,"location":2,"weight":107,"height":1.85,"sickness":"flu","id":16}
```

- Identifié mais mauvaise valeur donnée sensible : 1

```json
{"sickness":"flu","id":16,"age":38,"location":3,"weight":104,"height":1.88}
```

- Pas identifié : 5

```json
{"age":22,"location":3,"weight":90,"height":1.87,"sickness":"flu","id":15}
{"sickness":"flu","id":16,"age":53,"location":4,"weight":90,"height":2.06}
{"height":1.84,"sickness":"cancer","id":28,"age":67,"location":4,"weight":84}
{"weight":101,"height":1.88,"sickness":"cancer","id":28,"age":92,"location":3}
{"age":97,"location":3,"weight":92,"height":1.95,"sickness":"cancer","id":28}
```

- Mauvaise identification : 7

#### Distance_Manhattan2

```console
reidentification.Reidenfication(original, sigo, "manhattan", 3)

[map[age:44 height:1.57 location:4 sensitive:[cancer] weight:94]
map[age:42 height:1.62 location:1 sensitive:[cancer] weight:102]
map[age:54 height:2.08 location:1 sensitive:[flu] weight:105]
map[age:81 height:1.54 location:2 sensitive:[cancer] weight:89]
map[age:38 height:1.88 location:3 sensitive:[cancer] weight:104]
map[age:20 height:1.87 location:2 sensitive:[flu] weight:102]
map[age:22 height:1.87 location:3 sensitive:[flu] weight:90]
map[age:47 height:1.56 location:1 sensitive:[cancer] weight:103]
map[age:58 height:1.85 location:2 sensitive:[flu] weight:107]
map[age:51 height:1.97 location:1 sensitive:[flu] weight:100]
map[age:26 height:1.53 location:3 sensitive:[flu] weight:99]
map[age:50 height:1.8 location:2 sensitive:[flu] weight:109]
map[age:27 height:1.97 location:4 sensitive:[flu] weight:99]
map[age:84 height:2.05 location:3 sensitive:[cancer] weight:83]]
```

- Bien identifié : 4

```json
{"location":2,"weight":102,"height":1.87,"sickness":"flu","id":15,"age":20}
{"weight":99,"height":1.97,"sickness":"flu","id":15,"age":27,"location":4}
{"age":58,"location":2,"weight":107,"height":1.85,"sickness":"flu","id":16}
{"age":22,"location":3,"weight":90,"height":1.87,"sickness":"flu","id":15}
```

- Identifié mais mauvaise valeur donnée sensible : 1

```json
{"sickness":"flu","id":16,"age":38,"location":3,"weight":104,"height":1.88}
```

- Pas identifié : 5

```json
{"sickness":"flu","id":16,"age":53,"location":4,"weight":90,"height":2.06}
{"height":1.84,"sickness":"cancer","id":28,"age":67,"location":4,"weight":84}
{"weight":101,"height":1.88,"sickness":"cancer","id":28,"age":92,"location":3}
{"age":97,"location":3,"weight":92,"height":1.95,"sickness":"cancer","id":28}
```

- Mauvaise identification : 5

#### Distance_Canberra2

```console
reidentification.Reidenfication(original, sigo, "canberra", 3)

[map[age:19 height:1.67 location:2 sensitive:[flu] weight:87]
map[age:44 height:1.57 location:4 sensitive:[flu] weight:94]
map[age:26 height:2.01 location:3 sensitive:[flu] weight:81]
map[age:54 height:2.08 location:1 sensitive:[cancer] weight:105]
map[age:38 height:1.88 location:3 sensitive:[flu] weight:104]
map[age:20 height:1.87 location:2 sensitive:[flu] weight:102]
map[age:22 height:1.87 location:3 sensitive:[flu] weight:90]
map[age:67 height:1.84 location:4 sensitive:[cancer] weight:84]
map[age:53 height:2.06 location:4 sensitive:[flu] weight:90]
map[age:47 height:1.56 location:1 sensitive:[cancer] weight:103]
map[age:51 height:1.97 location:1 sensitive:[cancer] weight:100]
map[age:97 height:1.95 location:3 sensitive:[cancer] weight:92]
map[age:62 height:1.86 location:1 sensitive:[cancer] weight:106]
map[age:92 height:1.88 location:3 sensitive:[cancer] weight:101]
map[age:27 height:1.97 location:4 sensitive:[flu] weight:99]
map[age:84 height:2.05 location:3 sensitive:[cancer] weight:83]]
```

- Bien identifié : 8

```json
{"location":2,"weight":102,"height":1.87,"sickness":"flu","id":15,"age":20}
{"weight":99,"height":1.97,"sickness":"flu","id":15,"age":27,"location":4}
{"age":22,"location":3,"weight":90,"height":1.87,"sickness":"flu","id":15}
{"sickness":"flu","id":16,"age":38,"location":3,"weight":104,"height":1.88}
{"sickness":"flu","id":16,"age":53,"location":4,"weight":90,"height":2.06}
{"height":1.84,"sickness":"cancer","id":28,"age":67,"location":4,"weight":84}
{"weight":101,"height":1.88,"sickness":"cancer","id":28,"age":92,"location":3}
{"age":97,"location":3,"weight":92,"height":1.95,"sickness":"cancer","id":28}
```

- Pas identifié : 1

```json
{"age":58,"location":2,"weight":107,"height":1.85,"sickness":"flu","id":16}
```

- Mauvaise identification : 7

#### Distance_Chebyshev2

```console
reidentification.Reidenfication(original, sigo, "chebyshev", 3)

[map[age:93 height:1.81 location:1 sensitive:[cancer] weight:79]
map[age:44 height:1.57 location:4 sensitive:[cancer] weight:94]
map[age:42 height:1.62 location:1 sensitive:[cancer] weight:102]
map[age:54 height:2.08 location:1 sensitive:[flu] weight:105]
map[age:19 height:2 location:1 sensitive:[flu] weight:107]
map[age:81 height:1.54 location:2 sensitive:[cancer] weight:89]
map[age:38 height:1.88 location:3 sensitive:[cancer] weight:104]
map[age:20 height:1.87 location:2 sensitive:[flu] weight:102]
map[age:22 height:2.09 location:1 sensitive:[flu] weight:104]
map[age:92 height:1.75 location:1 sensitive:[cancer] weight:79]
map[age:26 height:1.97 location:1 sensitive:[flu] weight:110]
map[age:47 height:1.56 location:1 sensitive:[flu] weight:103]
map[age:58 height:1.85 location:2 sensitive:[flu] weight:107]
map[age:51 height:1.97 location:1 sensitive:[flu] weight:100]
map[age:26 height:1.53 location:3 sensitive:[flu] weight:99]
map[age:50 height:1.8 location:2 sensitive:[flu] weight:109]
map[age:27 height:1.97 location:4 sensitive:[flu] weight:99]
map[age:77 height:1.5 location:3 sensitive:[cancer] weight:86]
map[age:84 height:2.05 location:3 sensitive:[cancer] weight:83]]
```

- Bien identifié : 3

```json
{"location":2,"weight":102,"height":1.87,"sickness":"flu","id":15,"age":20}
{"weight":99,"height":1.97,"sickness":"flu","id":15,"age":27,"location":4}
{"age":58,"location":2,"weight":107,"height":1.85,"sickness":"flu","id":16}
```

- Pas identifié : 5

```json
{"age":22,"location":3,"weight":90,"height":1.87,"sickness":"flu","id":15}
{"sickness":"flu","id":16,"age":53,"location":4,"weight":90,"height":2.06}
{"height":1.84,"sickness":"cancer","id":28,"age":67,"location":4,"weight":84}
{"weight":101,"height":1.88,"sickness":"cancer","id":28,"age":92,"location":3}
{"age":97,"location":3,"weight":92,"height":1.95,"sickness":"cancer","id":28}
```

- Identifié mais mauvaise valeur donnée sensible : 1

```json
{"sickness":"flu","id":16,"age":38,"location":3,"weight":104,"height":1.88}
```

- Mauvaise identification : 10

#### Distance_Minkowski2 p=3

```console
reidentification.Reidenfication(original, sigo, "minkowski", 3)

[map[age:44 height:1.57 location:4 sensitive:[cancer] weight:94]
map[age:42 height:1.62 location:1 sensitive:[cancer] weight:102]
map[age:54 height:2.08 location:1 sensitive:[flu] weight:105]
map[age:19 height:2 location:1 sensitive:[flu] weight:107]
map[age:81 height:1.54 location:2 sensitive:[cancer] weight:89]
map[age:38 height:1.88 location:3 sensitive:[cancer] weight:104]
map[age:20 height:1.87 location:2 sensitive:[flu] weight:102]
map[age:22 height:2.09 location:1 sensitive:[flu] weight:104]
map[age:26 height:1.97 location:1 sensitive:[flu] weight:110]
map[age:47 height:1.56 location:1 sensitive:[flu] weight:103]
map[age:58 height:1.85 location:2 sensitive:[flu] weight:107]
map[age:51 height:1.97 location:1 sensitive:[flu] weight:100]
map[age:26 height:1.53 location:3 sensitive:[flu] weight:99]
map[age:50 height:1.8 location:2 sensitive:[flu] weight:109]
map[age:27 height:1.97 location:4 sensitive:[flu] weight:99]
map[age:84 height:2.05 location:3 sensitive:[cancer] weight:83]]
```

- Bien identifié : 3

```json
{"location":2,"weight":102,"height":1.87,"sickness":"flu","id":15,"age":20}
{"weight":99,"height":1.97,"sickness":"flu","id":15,"age":27,"location":4}
{"age":58,"location":2,"weight":107,"height":1.85,"sickness":"flu","id":16}
```

- Pas identifié : 5

```json
{"age":22,"location":3,"weight":90,"height":1.87,"sickness":"flu","id":15}
{"sickness":"flu","id":16,"age":53,"location":4,"weight":90,"height":2.06}
{"height":1.84,"sickness":"cancer","id":28,"age":67,"location":4,"weight":84}
{"weight":101,"height":1.88,"sickness":"cancer","id":28,"age":92,"location":3}
{"age":97,"location":3,"weight":92,"height":1.95,"sickness":"cancer","id":28}
```

- Identifié mais mauvaise valeur donnée sensible : 1

```json
{"sickness":"flu","id":16,"age":38,"location":3,"weight":104,"height":1.88}
```

- Mauvaise identification : 7

## Resumé

- 6 individus étaient ré-identifiables dans le `Dataset 1`.
- 9 individus étaient ré-identifiables dans le `Dataset 2`.

|                   | Dataset 1 | Dataset 2 |
|-------------------|:---------:|:---------:|
| Cosine_Similarity |    2/6    |    3/9    |
| Euclidean         |    6/6    |    3/9    |
| Manhattan         |    6/6    |    4/9    |
| Canberra          |    5/6    |    8/9    |
| Chebyshev         |    5/6    |    3/9    |
| Minkowski (p=3)   |    6/6    |    3/9    |

<https://towardsdatascience.com/17-types-of-similarity-and-dissimilarity-measures-used-in-data-science-3eb914d2681>
