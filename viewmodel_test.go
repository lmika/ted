package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModelViewCtrl_OpenRight(t *testing.T) {
	t.Run("should move cols to the right within the model", func(t *testing.T) {
		rwModel := NewStdModelFromSlice([][]string{
			{"letters", "numbers", "greek"},
			{"a", "1", "alpha"},
			{"b", "2", "bravo"},
			{"c", "3", "charlie"},
		})

		mvc := NewGridViewModel(rwModel)
		err := mvc.OpenRight(1)
		assert.NoError(t, err)

		assertModel(t, rwModel, [][]string{
			{"letters", "numbers", "", "greek"},
			{"a", "1", "", "alpha"},
			{"b", "2", "", "bravo"},
			{"c", "3", "", "charlie"},
		})
	})

	t.Run("should move cols to the right at the left of the model", func(t *testing.T) {
		rwModel := NewStdModelFromSlice([][]string{
			{"letters", "numbers", "greek"},
			{"a", "1", "alpha"},
			{"b", "2", "bravo"},
			{"c", "3", "charlie"},
		})

		mvc := NewGridViewModel(rwModel)
		err := mvc.OpenRight(0)
		assert.NoError(t, err)

		assertModel(t, rwModel, [][]string{
			{"letters", "", "numbers", "greek"},
			{"a", "", "1", "alpha"},
			{"b", "", "2", "bravo"},
			{"c", "", "3", "charlie"},
		})
	})

	t.Run("should move cols to the right at the right of the model", func(t *testing.T) {
		rwModel := NewStdModelFromSlice([][]string{
			{"letters", "numbers", "greek"},
			{"a", "1", "alpha"},
			{"b", "2", "bravo"},
			{"c", "3", "charlie"},
		})

		mvc := NewGridViewModel(rwModel)
		err := mvc.OpenRight(2)
		assert.NoError(t, err)

		assertModel(t, rwModel, [][]string{
			{"letters", "numbers", "greek", ""},
			{"a", "1", "alpha", ""},
			{"b", "2", "bravo", ""},
			{"c", "3", "charlie", ""},
		})
	})

	t.Run("should return error if row out of bounds", func(t *testing.T) {
		scenario := []int{-1, 3, 12}
		for _, scenario := range scenario {
			t.Run(fmt.Sprint(scenario), func(t *testing.T) {
				rwModel := NewStdModelFromSlice([][]string{
					{"letters", "numbers", "greek"},
					{"a", "1", "alpha"},
					{"b", "2", "bravo"},
					{"c", "3", "charlie"},
				})
				mvc := NewGridViewModel(rwModel)
				err := mvc.OpenRight(scenario)
				assert.Error(t, err)
			})
		}
	})
}

func assertModel(t *testing.T, actual Model, expected [][]string) {
	dr, dc := actual.Dimensions()
	assert.Equalf(t, len(expected), dr, "number of rows in model")

	for r, row := range expected {
		assert.Equalf(t, len(row), dc, "number of cols in row %v", r)
		for c, cell := range row {
			assert.Equalf(t, cell, actual.CellValue(r, c), "cell value at row %v, col %v", r, c)
		}
	}
}
