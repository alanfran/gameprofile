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

	It("Stores profiles", func() {
		p := Profile{
			ID:        "some_user",
			Coins:     999,
			Inventory: map[string]string{},
			Equipment: map[string]string{},
		}

		Expect(s.PutProfile(p)).Should(Succeed())

		p2, err := s.GetProfile(p.ID)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(p2).To(Equal(p))
	})

	Context("When retrieving a nonexistent key", func() {
		It("Returns an error", func() {
			_, err := s.GetProfile("this_does_not_exist")
			Expect(err).To(HaveOccurred())
		})
	})

	It("Stores and deletes punishments", func() {
		p := Punishment{
			ID:       1234,
			PlayerID: "some_user",
			By:       "some_admin",
			Type:     "ban",
			Reason:   "reason goes here",
			Date:     time.Now(),
			Expires:  time.Now().Add(time.Minute * 10),
		}

		// Store

		Expect(s.PutPunishment(p)).Should(Succeed())

		// Get

		ps, err := s.GetPunishments(p.PlayerID)
		Expect(err).ShouldNot(HaveOccurred())

		p2 := ps[p.Type]
		Expect(p2).To(Equal(p))

		// Delete

		Expect(s.DelPunishment(p2.ID)).Should(Succeed())

		ps2, err := s.GetPunishments(p.PlayerID)
		Expect(err).ShouldNot(HaveOccurred())

		_, ok := ps2[p.Type]
		Expect(ok).To(BeFalse())
	})

})
