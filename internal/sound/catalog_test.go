// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package sound

import "testing"

func TestCatalogValid(t *testing.T) {
	if err := Catalog().Validate(); err != nil {
		t.Fatalf("Catalog().Validate() returned %v", err)
	}
}
