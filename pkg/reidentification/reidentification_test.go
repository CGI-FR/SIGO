package reidentification_test

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	t.Parallel()

	row := jsonline.NewRow()
	row.Set("x", 3)
	row.Set("y", 6)
	record := infra.NewJSONLineRecord(&row, &[]string{"x", "y"}, &[]string{})
	test := []reidentification.Similarity{}

	for i := 0; i < 3; i++ {
		row1 := jsonline.NewRow()
		row1.Set("x", 3)
		row1.Set("y", 7)
		row1.Set("z", "c")

		record1 := infra.NewJSONLineRecord(&row1, &[]string{"x", "y"}, &[]string{"z"})
		sim := reidentification.NewSimilarity(i, record1, []string{"x", "y"}, []string{"z"})
		sim.ComputeSimilarity(record, []string{"x", "y"}, reidentification.NewCosineSimilarity())

		test = append(test, sim)
	}

	res, risk := reidentification.Recover(test)

	assert.True(t, risk)
	assert.Equal(t, []string{"c"}, res)
}

func TestCountValues(t *testing.T) {
	t.Parallel()

	values := []string{"a", "a", "b", "a", "c", "c", "a", "b"}
	count := reidentification.CountValues(values)

	assert.Equal(t, 4, count["a"])
	assert.Equal(t, 2, count["b"])
	assert.Equal(t, 2, count["c"])
}

func TestRisk(t *testing.T) {
	t.Parallel()

	row := jsonline.NewRow()
	row.Set("x", 11)
	row.Set("y", 9)
	record := infra.NewJSONLineRecord(&row, &[]string{"x", "y"}, &[]string{})
	test := []reidentification.Similarity{}

	for i := 0; i < 3; i++ {
		row1 := jsonline.NewRow()
		row1.Set("x", 19.67)
		row1.Set("y", 17.67)
		row1.Set("z", "b")

		record1 := infra.NewJSONLineRecord(&row1, &[]string{"x", "y"}, &[]string{"z"})
		sim := reidentification.NewSimilarity(i, record1, []string{"x", "y"}, []string{"z"})
		sim.ComputeSimilarity(record, []string{"x", "y"}, reidentification.NewCosineSimilarity())

		test = append(test, sim)
	}

	risk := reidentification.Risk(test)

	assert.Equal(t, float64(1), risk)

	test2 := []reidentification.Similarity{}
	z := []string{"a", "b", "b", "b"}

	for i := range z {
		row1 := jsonline.NewRow()
		row1.Set("x", 19.67)
		row1.Set("y", 17.67)
		row1.Set("z", z[i])

		record1 := infra.NewJSONLineRecord(&row1, &[]string{"x", "y"}, &[]string{"z"})
		sim := reidentification.NewSimilarity(i, record1, []string{"x", "y"}, []string{"z"})
		sim.ComputeSimilarity(record, []string{"x", "y"}, reidentification.NewCosineSimilarity())

		test2 = append(test2, sim)
	}

	risk2 := reidentification.Risk(test2)

	assert.Equal(t, float64(0.5), risk2)
}

//nolint: funlen
func TestReidentification(t *testing.T) {
	t.Parallel()

	// Importation dataset original
	originalFile, err := os.Open("../../examples/re-identification/data.json")
	assert.Nil(t, err)

	original, err := infra.NewJSONLineSource(bufio.NewReader(originalFile), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	// Liste des invidus du dataset original avec liste de similarités
	var res reidentification.Original

	i := 0

	// Pour chaque indivu du dataset original
	for original.Next() {
		qi := make(map[string]interface{})

		// On récupère les valeurs QI
		for _, q := range original.QuasiIdentifer() {
			qi[q] = original.Value().Row()[q]
		}

		// On créé un objet Similarities qui contient la liste des similarités avec les données anonymisées
		cosine := reidentification.NewCosineSimilarity()
		// euclidean := reidentification.NewEuclideanDistance()
		sims := reidentification.NewSimilarities(cosine)

		// Importation dataset Anonymisé
		sigoFile, err := os.Open("../../examples/re-identification/data2-sigo.json")
		assert.Nil(t, err)

		sigo, err := infra.NewJSONLineSource(bufio.NewReader(sigoFile), []string{"x", "y"}, []string{"z"})
		assert.Nil(t, err)

		j := 0

		// Pour chaque individu du dataset anonymisé
		for sigo.Next() {
			// Calcul la similarité avec l'individu original
			sim := reidentification.NewSimilarity(j, sigo.Value(), sigo.QuasiIdentifer(), sigo.Sensitive())
			sim.ComputeSimilarity(original.Value(), original.QuasiIdentifer(), sims.Metric())

			// Ajout à la liste des Similarités
			sims.Add(sim)
			j++
		}

		ind := reidentification.NewIndividu(i, qi, sims)
		i++

		// Ajout de l'individu original
		res.Add(ind)
	}

	// Calcul ré-identification sur les individu du dataset original avec en paramètre k
	riskInd := res.Reidenfication(3)
	expected := map[string]interface{}{
		"x": json.Number("20"), "y": json.Number("18"), "sensitive": []string{"b"},
	}
	assert.Contains(t, riskInd, expected)

	log.Println(riskInd)
}
