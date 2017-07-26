package profile

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mock", func() {
	var s Storer

	BeforeEach(func() {
		s = NewMockStore()
	})

	Context("Profiles", func() {
		It("Stores profiles", func() {
			p := Profile{
				ID:        "some_user",
				Coins:     999,
				Inventory: map[string]string{},
				Equipment: map[string]string{},
			}

			Expect(s.PutProfile(p)).Should(Succeed())

			p2, err := s.GetProfile(p.ID)
			Expect(err).ToNot(HaveOccurred())
			Expect(p2).To(Equal(p))
		})

		Context("When retrieving a nonexistent profile", func() {
			It("Returns an error", func() {
				_, err := s.GetProfile("this_does_not_exist")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("Punishments", func() {
		var testPunishment Punishment
		BeforeEach(func() {
			testPunishment = Punishment{
				ID:       1234,
				PlayerID: "some_user",
				By:       "some_admin",
				Type:     "ban",
				Reason:   "reason goes here",
				Date:     time.Now(),
				Expires:  time.Now().Add(time.Minute * 10),
			}
		})

		It("Stores and deletes punishments", func() {

			// Store

			Expect(s.PutPunishment(testPunishment)).Should(Succeed())

			// Get

			ps, err := s.GetPunishments(testPunishment.PlayerID)
			Expect(err).ShouldNot(HaveOccurred())

			p2 := ps[testPunishment.Type]
			Expect(p2).To(Equal(testPunishment))

			// Delete

			Expect(s.DelPunishment(p2.ID)).Should(Succeed())

			_, err = s.GetPunishments(testPunishment.PlayerID)
			Expect(err).Should(HaveOccurred())
		})

		Context("When a player has no punishments", func() {
			It("GetPunishments returns an error", func() {
				_, err := s.GetPunishments("this_user_has_no_punishments")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("When deleting a nonexistent punishment", func() {
			It("Fails", func() {
				Expect(s.DelPunishment(9001)).ToNot(Succeed())
			})
		})
	})

})
