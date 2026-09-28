package main

import (
	"log/slog"
	"os"

	"github.com/khwajalabs/regionctl/internal/spec"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// A real region spec would be loaded from a signed file; this proves the
	// content address is stable.
	s := spec.RegionSpec{Name: "au-sensitive-1", Classification: "sensitive", Isolated: true}
	log.Info("loaded region spec", "name", s.Name, "address", s.ContentAddress())
}
