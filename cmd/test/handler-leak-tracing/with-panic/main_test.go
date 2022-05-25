package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	main()
	os.Exit(m.Run())
}
