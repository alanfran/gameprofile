package profile

import (
	"io/ioutil"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bolt", func() {
	var s *BoltStore
	var file string

	BeforeEach(func() {
		f, err := ioutil.TempFile("", "test-boltdb")
		if err != nil {
			panic("Error creating temp db.")
		}
		file = f.Name()
		f.Close()

		s, err = NewBoltStore(file)
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		// close bolt db
		s.db.Close()
		os.Remove(file)
		// delete file
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

		Context("When storing a punishment", func() {
			Context("and all required fields are present", func() {
				It("succeeds", func() {
					Expect(s.PutPunishment(testPunishment)).Should(Succeed())
				})
			})

			Context("and either the PlayerID, By, or Type are missing", func() {
				It("fails", func() {
					incompletePunishments := []Punishment{
						Punishment{
							PlayerID: "someone",
							By:       "an_admin",
						},
						Punishment{
							By:   "an_admin",
							Type: "ban",
						},
						Punishment{
							PlayerID: "someone",
							Type:     "ban",
						},
					}

					for _, v := range incompletePunishments {
						Expect(s.PutPunishment(v)).ToNot(Succeed())
					}
				})
			})
		})

		Context("When retrieving punishments", func() {
			Context("and that player has punishments", func() {
				BeforeEach(func() {
					Expect(s.PutPunishment(testPunishment)).To(Succeed())
				})

				It("succeeds", func() {
					p, err := s.GetPunishments(testPunishment.PlayerID)
					Expect(err).ToNot(HaveOccurred())
					Expect(p).ToNot(BeZero())
				})
			})

			Context("and that player has no punishments", func() {
				It("returns an error", func() {
					_, err := s.GetPunishments("this_user_has_no_punishments")
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("When deleting a punishment", func() {
			Context("that exists", func() {
				BeforeEach(func() {
					Expect(s.PutPunishment(testPunishment)).To(Succeed())
				})

				It("succeeds", func() {
					Expect(s.DelPunishment(testPunishment.ID)).To(Succeed())
					// Verify
					_, err := s.GetPunishments(testPunishment.PlayerID)
					Expect(err).To(HaveOccurred())
				})
			})

			Context("that does not exist", func() {
				It("fails", func() {
					Expect(s.DelPunishment(9001)).ToNot(Succeed())
				})
			})
		})
	})

})
