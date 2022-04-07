package sigo_test

// func TestCountUniqueValues(t *testing.T) {
// 	t.Parallel()

// 	qi := []string{"x", "y"}
// 	kdtree := sigo.NewKDTreeFactory().New(3, 1, 2, qi)

// 	kdtree.AddValues(createRow(4, 1, qi))
// 	kdtree.AddValues(createRow(3, 2, qi))
// 	kdtree.AddValues(createRow(4, 3, qi))

// 	res := kdtree.CountUniqueValues()

// 	assert.Equal(t, 2, res["x"])
// 	assert.Equal(t, 3, res["y"])
// }

// func TestOrderMap(t *testing.T) {
// 	t.Parallel()

// 	values := map[string]int{"x": 2, "y": 3}

// 	res := sigo.Order(values)

// 	assert.Equal(t, []string{"y", "x"}, res)
// }
