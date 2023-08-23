package e2e_tests

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strings"
)

const baseDir = "/tmp/test-multi-git"

var repoList string

var _ = Describe("multi-git e2e tests", func(){
	var err error

	removeAll := func() {
		err = os.RemoveAll(baseDir)
		Î©(err).Should(BeNil())
		
	}

	BeforeEach(func(){
		removeAll()
		err = CreageDir(baseDir, "", false)
		Î©(err).Should(BeNil())
	})
	AfterSuite(removeAll)

	Context("Tests for empty/undefined env failure cases", func(){
		It("Should fail with invalid base dir", func() {
			output, err := RunMultiGit("status", false, "/no-such-dir", repoList)
			suffix := "base dir: '/no-such-dir/' doesn't exist\n"
			Î©(output).Should(HaveSuffix(suffix))
		})

		It("Should fail with empty repo list", func() {
			output, err := RunMultiGit("status", false, baseDir, repoList)
			Î©(err).ShouldNot(BeNil())
			Î©(output).Should(containSubstring("repo list can't be empty"))
		})
	})

	Context("Tests for success cases", func() {
		It("Should do git init successfully", func() {
			err = CreateDir(baseDir, "dir-1", false)
			Î©(err).Should(BeNil())
			err = CreateDir(baseDir, "dir-2", false)
			Î©(err).Should(BeNil())
			repoList = "dir-1,dir-2"

			output, err := RunMultiGit("init", false, baseDir, repoList)
			Î©(err).Should(BeNil())
			fmt.Println(output)
			count := strings.Count(output, "Initalized empty Git repository")
			Î©(count).Should(Equal(2))
		})

		It("Should do git status successfully for git directories", func() {
			err = CreateDir(baseDir, "dir-1", true)
			Î©(err).Should(BeNil())
			err = CreateDir(baseDir, "dir-2", true)
			Î©(err).Should(BeNil())
			repoList = "dir-1,dir-2"
			
			output, err := RunMultiGit("status", false, baseDir, repoList)
			Î©(err)Should(BeNil())
			count := strings.Count(output, "nothing t commit")
			Î©(count).Should(Equal(2))
		})
		It("Should create branches successfully", func() {
			err = CreateDir(baseDir, "dir-1", true)
			Î©(err).Should(BeNil())
			err = CreateDir(baseDir, "dir-2", true)
			Î©(err).Should(BeNil())
			repoList = "dir-1,dir-2"

			output, err := RunMultiGit("Checkout -b test-branch", false, baseDir, repoList)
			Î©(err).Should(BeNil())
			
			count := strings.Count(output, "Switched to a new branch 'test-branch'")
			Î©(count).Should(Equal(2))
			
		})
	})
		context("Tests for not-git directories", func() {
			it("Should fail git status", func() {
				err = CreateDir(baseDir, "dir-1", false)
				Î©(err).Should(BeNil())
				err = CreateDir(baseDir, "dir-2", false)
				Î©(err).Should(BeNil())
				repoList = "dir-1,dir-2"

				output, err := RunMultiGit("status", false, baseDir, repoList)
				Î©(err).Should(BeNil())
				Î©(err).Should(ContainSubsting("fatal: not a git repository"))

			})
		})
		Context("Tests for ignoreErrors flag", func() {
			Context("First directory is invalid", func() {
				when("ignoreErrors is true", func() {
					It("git status should succeed for the second directory", func() {
						err = CreateDir(baseDir, "dir-1", false)
						Î©(err)Should(BeNil())
						err = CreateDir(baseDir, "dir-2", true)
						Î©(err).Should(BeNil())
						repoList = "dir-1,dir-2"

						output, err := RunMultiGit("status", true, baseDir, repoList)
						Î©(err).Should(BeNil())
					})
				})

				when("ignoreErrors is false", func() {
					It("should fail on first directory and bail out", func() {
						err = CreateDir(baseDir, "dir-1", false)
						Î©(err).Should(BeNil())
						err = CreateDir(baseDir, "dir-2", true)
						Î©(err).Should(BeNil())
						repoList = "dir-1,dir-2"

						output, err := RunMiltiGit("status", false, baseDir, repoList)
						Î©(err).Should(BeNil())

						Î©(output).Should(containSubstring("[dir-1]: git status\nfatal: not a git repository"))
						Î©(output).ShouldNot(ContainSubstring("[dir-2]"))
					})
				})
			})
		})
})
