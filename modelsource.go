package main

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
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
	filename string
	options  CsvFileModelSourceOptions
}

type CsvFileModelSourceOptions struct {
	Comma rune
}

func NewCsvFileModelSource(filename string, options CsvFileModelSourceOptions) CsvFileModelSource {
	return CsvFileModelSource{
		filename: filename,
		options:  options,
	}
}

// Describes the source
func (s CsvFileModelSource) String() string {
	return filepath.Base(s.filename)
}

// Read the model from the given source
func (s CsvFileModelSource) Read() (Model, error) {
	// Check if the file exists.  If not, return an empty model
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		return NewSingleCellStdModel(), nil
	}

	f, err := os.Open(s.filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	model := new(StdModel)
	r := csv.NewReader(f)
	r.Comma = s.options.Comma
	r.FieldsPerRecord = -1
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		model.appendStr(record)
	}

	model.dirty = false
	return model, nil
}

func (s CsvFileModelSource) Write(m Model) error {
	f, err := os.Create(s.filename)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	w.Comma = s.options.Comma

	rows, cols := m.Dimensions()

	for r := 0; r < rows; r++ {
		record := make([]string, cols) // Reuse the record slice
		for c := 0; c < cols; c++ {
			record[c] = m.CellValue(r, c)
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
