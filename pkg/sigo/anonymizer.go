// Copyright (C) 2022 CGI France
//
// This file is part of SIGO.
//
// SIGO is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// SIGO is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with SIGO.  If not, see <http://www.gnu.org/licenses/>.

package sigo

import (
	"fmt"
	"os"

	"github.com/cgi-fr/jsonline/pkg/cast"
	"github.com/rs/zerolog/log"
)

const (
	laplace  = "laplace"
	gaussian = "gaussian"
)

type bounds struct {
	down, up float64
}

func NewNoAnonymizer() NoAnonymizer { return NoAnonymizer{} }

func NewGeneralAnonymizer() GeneralAnonymizer {
	return GeneralAnonymizer{boundsValues: make(map[string]map[string]bounds)}
}

func NewAggregationAnonymizer(typeAgg string) AggregationAnonymizer {
	return AggregationAnonymizer{typeAggregation: typeAgg, values: make(map[string]map[string]float64)}
}

func NewCodingAnonymizer() CodingAnonymizer {
	return CodingAnonymizer{}
}

func NewNoiseAnonymizer(mechanism string) NoiseAnonymizer {
	return NoiseAnonymizer{typeNoise: mechanism}
}

func NewSwapAnonymizer() SwapAnonymizer {
	return SwapAnonymizer{swapValues: make(map[string]map[string][]float64)}
}

func NewReidentification(args []string) Reidentification {
	return Reidentification{
		masked:           make(map[string][]map[string]interface{}),
		unique:           make(map[string]bool),
		sensitive:        make(map[string]map[string]interface{}),
		stats:            make(map[string]map[string]map[string]float64),
		sensitivesFields: args,
	}
}

type (
	NoAnonymizer      struct{}
	GeneralAnonymizer struct {
		// groupMap map[Cluster]map[string]string
		// map of cluster -> qi -> bounds
		boundsValues map[string]map[string]bounds
	}
	AggregationAnonymizer struct {
		typeAggregation string
		values          map[string]map[string]float64
	}
	CodingAnonymizer struct{}
	NoiseAnonymizer  struct {
		typeNoise string
	}
	SwapAnonymizer struct {
		swapValues map[string]map[string][]float64
	}
	AnonymizedRecord struct {
		original Record
		mask     map[string]interface{}
	}

	Reidentification struct {
		masked           map[string][]map[string]interface{}
		unique           map[string]bool
		sensitive        map[string]map[string]interface{}
		stats            map[string]map[string]map[string]float64
		sensitivesFields []string
	}
)

func (ar AnonymizedRecord) QuasiIdentifer() ([]float64, error) {
	return ar.original.QuasiIdentifer()
}

func (ar AnonymizedRecord) Sensitives() []interface{} {
	return ar.original.Sensitives()
}

func (ar AnonymizedRecord) Row() map[string]interface{} {
	original := ar.original.Row()
	for k, v := range ar.mask {
		original[k] = v
	}

	return original
}

// Anonymize returns the original record, there is no anonymization.
func (a NoAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	mask := map[string]interface{}{}
	for _, q := range qi {
		mask[q] = rec.Row()[q]
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

// Anonymize returns the record anonymize with the method general
// the record takes the bounds of the cluster.
func (a GeneralAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	mask := map[string]interface{}{}

	if a.boundsValues[clus.ID()] == nil {
		a.ComputeGeneralization(clus, qi)
	}

	for _, key := range qi {
		mask[key] = []float64{a.boundsValues[clus.ID()][key].down, a.boundsValues[clus.ID()][key].up}
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

// ComputeGeneralization calculates the min and max values of the cluster for each qi.
func (a GeneralAnonymizer) ComputeGeneralization(clus Cluster, qi []string) {
	values, _ := listValues(clus, qi)

	boundsVal := make(map[string]bounds)

	for _, key := range qi {
		var b bounds
		b.down = Min(values[key])
		b.up = Max(values[key])
		boundsVal[key] = b
	}

	a.boundsValues[clus.ID()] = boundsVal
}

// Anonymize returns the record anonymized with the method meanAggregarion or medianAggregation
// the record takes the aggregated values of the cluster.
func (a AggregationAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	mask := map[string]interface{}{}

	if a.values[clus.ID()] == nil {
		a.ComputeAggregation(clus, qi)
	}

	for _, key := range qi {
		mask[key] = a.values[clus.ID()][key]
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

// ComputeAggregation calculates the mean (method meanAggreagtion)
// or median (method medianAggregation) value of the cluster for each qi.
func (a AggregationAnonymizer) ComputeAggregation(clus Cluster, qi []string) {
	values, _ := listValues(clus, qi)

	valAggregation := make(map[string]float64)

	for _, key := range qi {
		switch a.typeAggregation {
		case "mean":
			valAggregation[key] = Mean(values[key])
		case "median":
			valAggregation[key] = Median(values[key])
		}
	}

	a.values[clus.ID()] = valAggregation
}

// Anonymize returns the record anonymized with the method outlier
// if the record is in the interval [Q1;Q3] then we don't change its value
// if the record is > Q3 then it takes the Q3 value
// if the record is < Q1 then it takes the Q1 value.
func (a CodingAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	values, _ := listValues(clus, qi)
	mask := map[string]interface{}{}

	for i, key := range qi {
		vals := values[key]
		q := Quartile(vals)
		bottom := q.Q1
		top := q.Q3

		recVals, err := rec.QuasiIdentifer()
		if err != nil {
			fmt.Println(err)
			break
		}
		val := recVals[i]

		switch {
		case val < bottom:
			mask[key] = bottom
		case val > top:
			mask[key] = top
		default:
			mask[key] = val
		}
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

// Anonymize returns the record anonymized with the method laplaceNoise or gaussianNoise
// the record takes as value the original value added to a Laplacian or Gaussian noise
// the anonymized value stays within the bounds of the cluster.
func (a NoiseAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	values, _ := listValues(clus, qi)
	mask := map[string]interface{}{}

	for i, key := range qi {
		recVals, err := rec.QuasiIdentifer()
		if err != nil {
			fmt.Println(err)
			break
		}
		val := recVals[i]

		laplaceVal := Scaling(val, values[key], laplace)
		gaussianVal := Scaling(val, values[key], gaussian)

		var randomVal float64

		for {
			switch a.typeNoise {
			case laplace:
				randomVal = Rescaling(laplaceVal+LaplaceNumber(), values[key], laplace)
			case gaussian:
				randomVal = Rescaling(gaussianVal+GaussianNumber(0, 1), values[key], gaussian)
			}

			if (randomVal > Min(values[key]) && randomVal < Max(values[key])) || Min(values[key]) == Max(values[key]) {
				break
			}
		}

		mask[key] = randomVal
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

func (a SwapAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	mask := map[string]interface{}{}

	// cluster value swapping
	if a.swapValues[clus.ID()] == nil {
		a.Swap(clus, qi)
	}

	var idx int

	// retrieve the position (idx) of the record in the cluster
	for i, r := range clus.Records() {
		if rec == r {
			idx = i
		}
	}

	for _, key := range qi {
		// retrieve the swapped value
		mask[key] = a.swapValues[clus.ID()][key][idx]
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

func (a SwapAnonymizer) Swap(clus Cluster, qi []string) {
	// retrieve the cluster values for each qi
	values, _ := listValues(clus, qi)

	swapVal := make(map[string][]float64)

	for _, key := range qi {
		// values permutation
		swapVal[key] = Shuffle(values[key])
	}

	a.swapValues[clus.ID()] = swapVal
}

// Anonymize on object Reidentification re-identifies the original data using the anonymized data.
func (r Reidentification) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	mask := map[string]interface{}{}

	// initialize re-identification object: groups the anonymized data,
	// checks if the sensitive data is not unique in the cluster,
	// checks if the anonymized data of the cluster have the same qi value,
	// and computes the mean and the standard deviation of the cluster
	if r.masked[clus.ID()] == nil {
		r.InitReidentification(clus, qi, s)
	}

	original, err := cast.ToInt64(rec.Sensitives()[0])
	if err != nil {
		log.Err(err).Msg("Cannot cast original value")
		log.Warn().Int("return", 1).Msg("End SIGO")
		os.Exit(1)
	}

	// re-identification of the original data
	if original.(int64) == 1 {
		for _, q := range qi {
			mask[q] = rec.Row()[q]
		}

		for _, sensitive := range r.sensitivesFields {
			// if in a cluster the sensitive data is unique then we can re-identify the individuals
			if r.sensitive[clus.ID()][sensitive] != nil {
				mask[r.sensitivesFields[0]] = r.sensitive[clus.ID()][sensitive]
				mask["similarity"] = 1
			} else if !r.unique[clus.ID()] {
				// else if the masked data are not unique, then the distances are computed
				// (if they are unique, impossible to re-identify)
				scores := r.ComputeSimilarity(rec, clus, qi, s)
				sim, sens := TopSimilarity(scores)
				mask[r.sensitivesFields[0]] = sens
				mask["similarity"] = sim
			}
		}
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

// InitReidentification initialize the re-identification object.
func (r Reidentification) InitReidentification(clus Cluster, qi []string, s []string) {
	// map containing the anonymized data
	maskedData := []map[string]interface{}{}
	// slice containing the records of cluster
	data := []map[string]interface{}{}

	// map containing for each sensitive attribute the list of sensitive data
	sensitivesData := make(map[string][]interface{})

	for _, rec := range clus.Records() {
		data = append(data, rec.Row())

		original, err := cast.ToInt64(rec.Sensitives()[0])
		if err != nil {
			log.Err(err).Msg("Cannot cast original value")
			log.Warn().Int("return", 1).Msg("End SIGO")
			os.Exit(1)
		}

		if original.(int64) == 0 {
			maskedData = append(maskedData, rec.Row())

			for _, s := range r.sensitivesFields {
				sensitivesData[s] = append(sensitivesData[s], rec.Row()[s])
			}
		}
	}

	if len(maskedData) == 0 {
		log.Error().Msg("Clusters with only original data, pay attention to the l-diversity parameter ")
		log.Warn().Int("return", 1).Msg("End SIGO")
		os.Exit(1)
	}

	// groups all the anonymized records
	r.masked[clus.ID()] = maskedData
	// indicates if the cluster contains unique masked data
	r.unique[clus.ID()] = Unique(maskedData, qi)
	// checks if the sensitive data is well represented
	uniqueSensitive := IsUnique(sensitivesData)

	tmp := make(map[string]interface{})

	for _, s := range r.sensitivesFields {
		if uniqueSensitive[s] {
			tmp[s] = sensitivesData[s][0]
		} else {
			tmp[s] = nil
		}
	}

	r.sensitive[clus.ID()] = tmp

	// computes the mean and standard deviation of each cluster
	r.ComputeStatistics(data, clus, s)
}

// ComputeStatistics computes the mean and standart deviation for cluster clus.
func (r Reidentification) ComputeStatistics(data []map[string]interface{}, clus Cluster, s []string) {
	statistics := make(map[string]map[string]float64)

	for key, val := range ListValues(data, append(s, r.sensitivesFields...)) {
		stats := make(map[string]float64)
		stats["mean"] = Mean(SliceToFloat64(val))
		stats["std"] = Std(SliceToFloat64(val))

		statistics[key] = stats
	}

	r.stats[clus.ID()] = statistics
}

// Statistics returns the statistics of the q attribute of the cluster with path idCluster.
func (r Reidentification) Statistics(idCluster string, q string) (mean float64, std float64) {
	return r.stats[idCluster][q]["mean"], r.stats[idCluster][q]["std"]
}

// ComputeSimilarity computes the similarity score between the record rec and the anonymized cluster data.
func (r Reidentification) ComputeSimilarity(rec Record, clus Cluster,
	qi []string, s []string,
) map[float64]interface{} {
	scores := make(map[float64]interface{})

	x := make(map[string]interface{})

	for _, q := range qi {
		mean, std := r.Statistics(clus.ID(), q)
		x[q] = Scale(rec.Row()[q], mean, std)
	}

	X := MapItoMapF(x)

	for _, row := range r.masked[clus.ID()] {
		y := make(map[string]interface{})

		for _, q := range qi {
			mean, std := r.Statistics(clus.ID(), q)
			y[q] = Scale(row[q], mean, std)
		}

		Y := MapItoMapF(y)

		// Compute similarity
		score := Similarity(ComputeDistance("", X, Y))

		scores[score] = row[r.sensitivesFields[0]]
	}

	return scores
}

// Returns the list of values present in the cluster for each qi.
func listValues(clus Cluster, qi []string) (mapValues map[string][]float64, err error) {
	mapValues = make(map[string][]float64)

	for _, record := range clus.Records() {
		for i, key := range qi {
			vals, _ := record.QuasiIdentifer()
			// if err != nil {
			// 	fmt.Println(err)
			// 	return map[string][]float64{}, err
			// }
			val := vals[i]
			mapValues[key] = append(mapValues[key], val)
		}
	}

	return mapValues, nil
}
