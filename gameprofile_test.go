package main

import (
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
	"./profile"
)

var _ = Describe("Game profile microservice", func() {
	var s profile.Storer
	Context("Using mock store", func() {
		BeforeEach(func() {
			s = profile.NewMockStore()
		})

		Context("When trying to update an object that has been changed since the client's last request", func() {
			PIt("Returns status code 409 Conflict along with the new state of the object", func() {

			})
		})

		Context("When updating an untouched object", func() {
			PIt("Returns status code 204 No Content", func() {

			})
		})
	})

})
