package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path"

	_ "image/jpeg"

	log "github.com/sirupsen/logrus"
)

type Piece struct {
	PuzzleRow int
	PuzzleCol int
	Image     image.Image
}

type Puzzle struct {
	Pieces  []Piece
	NumRows int
	NumCols int
	Image   *image.NRGBA
}

type PuzzleMetadata struct {
	NumRows int
	NumCols int
}

func CreatePuzzleFromFile(inputFile string, pieceSize int, seed int64) (Puzzle, error) {
	var puz Puzzle
	// load the image file
	f, err := os.Open(inputFile)
	if err != nil {
		return puz, errors.New(fmt.Sprintf("Error loading input file %s", err.Error()))
	}

	r := bufio.NewReader(f)

	raw_image, _, err := image.Decode(r)
	if err != nil {
		return puz, errors.New(fmt.Sprintf("Error decoding image file %s", err.Error()))
	}

	// convert to nrgba image so we can keep things lossless
	puz.Image = image.NewNRGBA(raw_image.Bounds())
	for x := raw_image.Bounds().Min.X; x < raw_image.Bounds().Max.X; x++ {
		for y := raw_image.Bounds().Min.Y; y < raw_image.Bounds().Max.Y; y++ {
			puz.Image.Set(x, y, raw_image.At(x, y))
		}
	}
	log.WithField("image_bounds", raw_image.Bounds()).Debug("Raw image bounds")

	// figure out how many pieces we will have in each dimension - crop from the bottom right (for now)
	dimensionPoint := puz.Image.Bounds().Max.Sub(puz.Image.Bounds().Min)
	puz.NumRows = dimensionPoint.Y / pieceSize
	puz.NumCols = dimensionPoint.X / pieceSize
	log.Debug(fmt.Sprintf("Dimensions %v", dimensionPoint))
	log.Debug(fmt.Sprintf("%d rows %d cols", puz.NumRows, puz.NumCols))

	// cut up the image to create the pieces
	for row := 0; row < puz.NumRows; row++ {
		for col := 0; col < puz.NumCols; col++ {
			pieceMin := image.Point{
				X: puz.Image.Bounds().Min.X + (col * pieceSize),
				Y: puz.Image.Bounds().Min.Y + (row * pieceSize),
			}
			pieceRect := image.Rectangle{
				Min: pieceMin,
				Max: image.Point{
					X: pieceMin.X + pieceSize,
					Y: pieceMin.Y + pieceSize,
				},
			}
			piece := Piece{
				PuzzleRow: row,
				PuzzleCol: col,
				Image:     puz.Image.SubImage(pieceRect),
			}
			puz.Pieces = append(puz.Pieces, piece)
		}
	}

	return puz, nil
}

func CreateFromDirectory(inputDir string) error {
	rowFileInfos, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return err
	}
	for _, file := range rowFileInfos {

	}

}

func (puz Puzzle) Get(row int, column int) (Piece, error) {
	var p Piece
	for _, p := range puz.Pieces {
		if p.PuzzleRow == row && p.PuzzleCol == column {
			return p, nil
		}
	}

	return p, errors.New(fmt.Sprintf("No piece with row %d and col %d", row, column))
}

func (puz Puzzle) WriteToDirectory(outputDir string) error {
	// save the original puzzle image to a file
	puzFile, err := os.Create(path.Join(outputDir, "original_image.png"))
	if err != nil {
		return err
	}
	err = png.Encode(puzFile, puz.Image)
	err = puzFile.Close()
	if err != nil {
		return err
	}
	// save some puzzle metadata to a file
	puzMetaJson, err := json.Marshal(PuzzleMetadata{NumRows: puz.NumRows, NumCols: puz.NumCols})
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(outputDir, "puzzle.json"), puzMetaJson, 0644)
	if err != nil {
		return err
	}

	for row := 0; row < puz.NumRows; row++ {
		// create the dir for the row
		rowDir := path.Join(outputDir, fmt.Sprintf("%d", row))

		err := os.MkdirAll(rowDir, 0777)
		if err != nil && !os.IsExist(err) {
			return errors.New(fmt.Sprintf("Error creating row %d directory %s: %s", row, rowDir, err.Error()))
		}

		for col := 0; col < puz.NumCols; col++ {
			file_location := path.Join(rowDir, fmt.Sprintf("%d.png", col))

			errorTags := fmt.Sprintf("row: %d col: %d output: %s", row, col, outputDir)

			piece, err := puz.Get(row, col)
			if err != nil {
				return errors.New(fmt.Sprintf("%s Error getting piece from puzzle: %s", errorTags, err.Error()))
			}

			f, err := os.Create(file_location)
			if err != nil {
				return errors.New(fmt.Sprintf("%s Error opening piece file: %s", errorTags, err.Error()))
			}

			err = png.Encode(f, piece.Image)
			if err != nil {
				return errors.New(fmt.Sprintf("%s Error writing piece to file: %s", errorTags, err.Error()))
			}

			err = f.Close()
			if err != nil {
				return errors.New(fmt.Sprintf("%s Error closing piece file: %s", errorTags, err.Error()))
			}
		}
	}
	return nil
}
