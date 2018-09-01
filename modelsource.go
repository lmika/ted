package main

import (
	"path/filepath"
	"os"
	"encoding/csv"
	"io"
)

// ModelSource is a source of models.  At a minimum, it must be able to read models.
type ModelSource interface {
	// Describes the source
	String() string

	// Read the model from the given source
	Read() (Model, error)
}

// Writable models take a model and write it to the source
type WritableModelSource interface {
	ModelSource

	// Write writes a model to the source
	Write(m Model) error
}

// A model source backed by a CSV file
type CsvFileModelSource struct {
	Filename string
}

// Describes the source
func (s CsvFileModelSource) String() string {
	return filepath.Base(s.Filename)
}

// Read the model from the given source
func (s CsvFileModelSource) Read() (Model, error) {
	f, err := os.Open(s.Filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	model := new(StdModel)
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		model.appendStr(record)
	}

	model.dirty = false
	return model, nil
}

func (s CsvFileModelSource) Write(m Model) error {
	f, err := os.Create(s.Filename)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)

	rows, cols := m.Dimensions()

	for r := 0; r < rows; r++ {
		record := make([]string, cols)		// Reuse the record slice
		for c := 0; c < cols; c++ {
			record[c] =  m.CellValue(r, c)
		}
		if err := w.Write(record); err != nil {
			f.Close()
			return err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		f.Close()
		return err
	}

	return f.Close()
}