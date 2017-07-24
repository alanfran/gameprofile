package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGameprofile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gameprofile Suite")
}
