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
	Pieces    []Piece
	NumRows   int
	NumCols   int
	PieceSize int
	Image     *image.NRGBA
}

type PuzzleMetadata struct {
	NumRows   int
	NumCols   int
	PieceSize int
}

func loadNrgbaImage(imagePath string) (*image.NRGBA, error) {
	var nrgba_image *image.NRGBA
	f, err := os.Open(imagePath)
	if err != nil {
		return nrgba_image, errors.New(fmt.Sprintf("Error loading input file %s", err.Error()))
	}

	r := bufio.NewReader(f)

	raw_image, _, err := image.Decode(r)
	if err != nil {
		return nrgba_image, errors.New(fmt.Sprintf("Error decoding image file %s", err.Error()))
	}

	nrgba_image = image.NewNRGBA(raw_image.Bounds())
	for x := raw_image.Bounds().Min.X; x < raw_image.Bounds().Max.X; x++ {
		for y := raw_image.Bounds().Min.Y; y < raw_image.Bounds().Max.Y; y++ {
			nrgba_image.Set(x, y, raw_image.At(x, y))
		}
	}

	return nrgba_image, nil
}

func CreatePuzzleFromFile(inputFile string, pieceSize int, seed int64) (Puzzle, error) {
	var puz Puzzle
	var err error
	puz.Image, err = loadNrgbaImage(inputFile)
	if err != nil {
		return puz, err
	}

	// figure out how many pieces we will have in each dimension - crop from the bottom right (for now)
	dimensionPoint := puz.Image.Bounds().Max.Sub(puz.Image.Bounds().Min)
	puz.NumRows = dimensionPoint.Y / pieceSize
	puz.NumCols = dimensionPoint.X / pieceSize
	puz.PieceSize = pieceSize
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

func CreatePuzzleFromDirectory(inputDir string) (Puzzle, error) {
	var puz Puzzle

	puzzleMetadataJson, err := ioutil.ReadFile(path.Join(inputDir, "puzzle.json"))
	if err != nil {
		return puz, err
	}
	var puzzleMetadata PuzzleMetadata
	err = json.Unmarshal(puzzleMetadataJson, &puzzleMetadata)
	if err != nil {
		return puz, err
	}
	puz.NumRows = puzzleMetadata.NumRows
	puz.NumCols = puzzleMetadata.NumCols
	puz.PieceSize = puzzleMetadata.PieceSize

	puz.Image, err = loadNrgbaImage(path.Join(inputDir, "original_image.png"))

	for row := 0; row < puz.NumRows; row++ {
		rowDir := path.Join(inputDir, fmt.Sprintf("%d", row))
		for col := 0; col < puz.NumCols; col++ {
			piecePath := path.Join(rowDir, fmt.Sprintf("%d.png", col))
			pieceImage, err := loadNrgbaImage(piecePath)
			if err != nil {
				return puz, err
			}
			piece := Piece{
				PuzzleRow: row,
				PuzzleCol: col,
				Image:     pieceImage,
			}
			puz.Pieces = append(puz.Pieces, piece)
		}
	}

	return puz, nil
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
	err := os.MkdirAll(outputDir, 0777)
	if err != nil && !os.IsExist(err) {
		return errors.New(fmt.Sprintf("Error creating output directory %s: %s", outputDir, err.Error()))
	}
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
	puzMetaJson, err := json.Marshal(PuzzleMetadata{NumRows: puz.NumRows, NumCols: puz.NumCols, PieceSize: puz.PieceSize})
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
