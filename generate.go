package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func CreatePuzzle(inputFile string, outputDir string, pieceSize int, seed int64) error {
	log.Info(fmt.Sprintf("%s %s %d %d", inputFile, outputDir, pieceSize, seed))
	return nil
}
