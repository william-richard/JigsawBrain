package main

import (
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePuzzleFromFile(t *testing.T) {
	cases := []string{"1.png", "2.jpeg"}

	for _, caseInputFile := range cases {
		puzzleOutputDir := path.Join("testdata", "output", strings.Split(caseInputFile, ".")[0])
		puzzleFromDir, err := CreatePuzzleFromDirectory(puzzleOutputDir)
		assert.Nil(t, err)

		caseInputPath := path.Join("testdata", "input", caseInputFile)
		puzzleFromFile, err := CreatePuzzleFromFile(caseInputPath, puzzleFromDir.PieceSize, 123456)
		assert.Nil(t, err)

		assert.Equal(t, puzzleFromDir.NumRows, puzzleFromFile.NumRows, "puzzle row count should be equal")
		assert.Equal(t, puzzleFromDir.NumCols, puzzleFromFile.NumCols, "puzzle col count should be equal")
		assert.Equal(t, puzzleFromFile.PieceSize, puzzleFromFile.PieceSize, "puzzle piece size should be equal")

		for row := 0; row < puzzleFromDir.NumRows; row++ {
			for col := 0; col < puzzleFromDir.NumCols; col++ {
				pieceFromDir, err := puzzleFromDir.Get(row, col)
				assert.Nil(t, err)
				pieceFromFile, err := puzzleFromFile.Get(row, col)
				assert.Nil(t, err)
				assert.Equal(t, pieceFromDir.PuzzleRow, row, "piece row should be right")
				assert.Equal(t, pieceFromFile.PuzzleRow, row, "piece row should be right")
				// TODO check that pieces sizes are correct
				// TODO check that pieces colors match
			}
		}
	}
}
