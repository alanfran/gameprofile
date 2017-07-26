package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/alanfran/gameprofile/profile"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Game profile microservice", func() {
	var app *App
	var resp *httptest.ResponseRecorder
	var testProfile profile.Profile

	BeforeSuite(func() {
		gin.SetMode(gin.TestMode)
	})

	Context("Using mock store", func() {
		BeforeEach(func() {
			app = NewApp(profile.NewMockStore())
			resp = httptest.NewRecorder()
			testProfile = profile.Profile{
				ID:        "test_profile",
				Coins:     1234,
				Inventory: map[string]string{},
				Equipment: map[string]string{},
			}
		})

		Context("/:steamid", func() {
			Context("GET", func() {
				Context("When the profile exists", func() {
					BeforeEach(func() {
						Expect(app.profiles.PutProfile(testProfile)).To(Succeed())
					})

					It("returns 200 Success and the ProfileWithHash.", func() {
						req, err := http.NewRequest("GET", "/"+testProfile.ID, nil)
						Expect(err).ToNot(HaveOccurred())
						app.engine.ServeHTTP(resp, req)

						result := resp.Result()
						Expect(result.StatusCode).To(Equal(http.StatusOK))

						body, err := ioutil.ReadAll(resp.Body)
						Expect(err).ToNot(HaveOccurred())

						var profileWithHash ProfileWithHash
						var profileWithoutHash profile.Profile

						err = json.Unmarshal(body, &profileWithHash)
						Expect(err).ToNot(HaveOccurred())
						Expect(profileWithHash.Hash).ToNot(BeZero())

						err = json.Unmarshal(body, &profileWithoutHash)
						Expect(err).ToNot(HaveOccurred())
						Expect(profileWithoutHash).To(Equal(testProfile))
					})
				})

				Context("When the profile does not exist", func() {
					It("returns 404 Not Found", func() {
						req, err := http.NewRequest("GET", "/this_profile_should_not_exist", nil)
						Expect(err).ToNot(HaveOccurred())
						app.engine.ServeHTTP(resp, req)

						result := resp.Result()
						Expect(result.StatusCode).To(Equal(http.StatusNotFound))
					})
				})
			})

			Context("POST", func() {
				Context("When there is no profile with that SteamID", func() {
					It("returns 201 Created along with the new ProfileWithHash", func() {
						postJSON, err := json.Marshal(testProfile)
						Expect(err).ToNot(HaveOccurred())

						req, err := http.NewRequest("POST", "/"+testProfile.ID, bytes.NewBuffer(postJSON))
						Expect(err).ToNot(HaveOccurred())
						req.Header.Set("Content-Type", "application/json")
						app.engine.ServeHTTP(resp, req)

						result := resp.Result()
						Expect(result.StatusCode).To(Equal(http.StatusCreated))

						body, err := ioutil.ReadAll(resp.Body)
						Expect(err).ToNot(HaveOccurred())

						var profileWithHash ProfileWithHash
						var profileWithoutHash profile.Profile

						err = json.Unmarshal(body, &profileWithHash)
						Expect(err).ToNot(HaveOccurred())
						Expect(profileWithHash.Hash).ToNot(BeZero())

						err = json.Unmarshal(body, &profileWithoutHash)
						Expect(err).ToNot(HaveOccurred())
						Expect(profileWithoutHash).To(Equal(testProfile))
					})
				})

				Context("When a profile with that SteamID already exists", func() {
					BeforeEach(func() {
						Expect(app.profiles.PutProfile(testProfile)).To(Succeed())
					})

					PIt("returns 490 Conflict along with the current ProfileWithHash", func() {

					})
				})
			})

			Context("PUT", func() {
				BeforeEach(func() {
					Expect(app.profiles.PutProfile(testProfile)).To(Succeed())
				})

				Context("When the hash is valid", func() {
					It("returns 200 OK and returns the new ProfileWithHash", func() {
						pwh := NewProfileWithHash(testProfile)
						pwh.Coins = 9999

						postJSON, err := json.Marshal(pwh)
						Expect(err).ToNot(HaveOccurred())

						req, err := http.NewRequest("PUT", "/"+testProfile.ID, bytes.NewBuffer(postJSON))
						Expect(err).ToNot(HaveOccurred())
						req.Header.Set("Content-Type", "application/json")
						app.engine.ServeHTTP(resp, req)

						result := resp.Result()
						Expect(result.StatusCode).To(Equal(http.StatusOK))

						body, err := ioutil.ReadAll(resp.Body)
						Expect(err).ToNot(HaveOccurred())

						var profileWithHash ProfileWithHash
						err = json.Unmarshal(body, &profileWithHash)
						Expect(err).ToNot(HaveOccurred())

						Expect(profileWithHash.Hash).ToNot(BeZero())
						Expect(profileWithHash.Hash).ToNot(Equal(pwh.Hash))
						Expect(profileWithHash.Coins).To(Equal(pwh.Coins))
					})
				})

				Context("When the hash is invalid", func() {
					It("returns 409 Conflict along with the new ProfileWithHash", func() {
						postJSON, err := json.Marshal(testProfile)
						Expect(err).ToNot(HaveOccurred())

						req, err := http.NewRequest("PUT", "/"+testProfile.ID, bytes.NewBuffer(postJSON))
						Expect(err).ToNot(HaveOccurred())
						req.Header.Set("Content-Type", "application/json")
						app.engine.ServeHTTP(resp, req)

						result := resp.Result()
						Expect(result.StatusCode).To(Equal(http.StatusConflict))

						body, err := ioutil.ReadAll(resp.Body)
						Expect(err).ToNot(HaveOccurred())

						var profileWithHash ProfileWithHash
						err = json.Unmarshal(body, &profileWithHash)
						Expect(err).ToNot(HaveOccurred())

						Expect(profileWithHash.Hash).ToNot(BeZero())
					})
				})
			})
		})

		Context("/:steamid/punishments", func() {
			var testPunishment profile.Punishment

			BeforeEach(func() {
				testPunishment = profile.Punishment{
					ID:       12345,
					PlayerID: testProfile.ID,
					By:       "an_admin",
					Type:     "ban",
					Reason:   "testing",
					Date:     time.Now(),
					Expires:  time.Now().Add(time.Minute * 10),
				}
			})

			Context("GET", func() {
				Context("there are no punishments for that steam id", func() {
					It("returns 404 Not Found", func() {
						req, err := http.NewRequest("GET", "/"+testProfile.ID+"/punishments", nil)
						Expect(err).ToNot(HaveOccurred())
						app.engine.ServeHTTP(resp, req)

						result := resp.Result()
						Expect(result.StatusCode).To(Equal(http.StatusNotFound))
					})
				})

				Context("there are punishments for that steam id", func() {

					BeforeEach(func() {
						Expect(app.profiles.PutPunishment(testPunishment)).To(Succeed())
					})

					It("returns 200 Success along with the Punishments", func() {
						req, err := http.NewRequest("GET", "/"+testProfile.ID+"/punishments", nil)
						Expect(err).ToNot(HaveOccurred())
						app.engine.ServeHTTP(resp, req)

						result := resp.Result()
						Expect(result.StatusCode).To(Equal(http.StatusOK))

						body, err := ioutil.ReadAll(resp.Body)
						Expect(err).ToNot(HaveOccurred())

						var punishmentResults map[string]profile.Punishment
						Expect(json.Unmarshal(body, &punishmentResults)).To(Succeed())

						Expect(punishmentResults).To(Equal(map[string]profile.Punishment{testPunishment.Type: testPunishment}))
					})
				})
			})

			Context("POST", func() {
				It("returns 204 No Content and stores the punishment object", func() {
					postJSON, err := json.Marshal(testPunishment)
					Expect(err).ToNot(HaveOccurred())

					req, err := http.NewRequest("POST", "/"+testProfile.ID+"/punishments", bytes.NewBuffer(postJSON))
					Expect(err).ToNot(HaveOccurred())
					req.Header.Set("Content-Type", "application/json")
					app.engine.ServeHTTP(resp, req)

					result := resp.Result()
					Expect(result.StatusCode).To(Equal(http.StatusNoContent))

					p, err := app.profiles.GetPunishments(testPunishment.PlayerID)
					Expect(err).ToNot(HaveOccurred())

					Expect(p).To(Equal(testPunishment))
				})
			})

			Context("PUT", func() {
				PIt("returns 204 No Content and updates the punishments list", func() {

				})
			})
		})

	})

})
