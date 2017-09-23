package moka_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMoka(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Moka Suite")
}
