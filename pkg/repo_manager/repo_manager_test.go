package repo_manager
import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/DJ-66/multi-git/pkg/helpers"
	"os"
	"path"
	"strings"
)

const baseDir = "tmp/test-multi-git"
var repoList = []string{}

var _ = Describe("Repo manager tests", func(){
	var err error

	removeAll := func() {
		err = os.RemoveAll(baseDir)
		Î©(err).Should(BeNil())
	}
	BeforeEach(func() {
		removeAll()
		err = CreateDir(baseDir, "dir-1", true)
		Î©(err).Should(BeNil())
		repoList = []string{"dir-1"}
	})
	AfterEach(removeAll)

	Context("Tests for failure cases", func() {
		It("Should fail with invalid base dir", func() {
			_, err := NewRepoManager("/no-such-dir", repoList, true)
			Î©(err).ShouldNot(BeNil())
		})

		It("Should fail with empty repo list", func() {
			_, err := NewRepoManager(baseDir, []string{}, true)
			Î©(err).ShouldNot(BeNil())
		})
	})
		Context("Test for success cases", func(){
			It("Should get repo list successfully", func(){
				rm, err := NewRepoManager(baseDir, repoList, true)
				Î©(err).Should(BeNil())

				repos := rm.GetRepos()
				Î©(repos).Should(HaveLen(1))
				actual := path.Join(baseDir, repoList[0])
				expected := repos[0]
				Î©(actual).Should(Equal(expected))
			})
			It("Should get repo list successfully with non-git directories", func(){
				repoList = append(repoList, "dir-2")
				CreateDir(baseDir, repoList[1], true)
				CreateDir(baseDir, "not-a-git-repo", false)
				rm, err := NewRepoManager(baseDir, repoList, true)
				Î©(err).Should(BeNil())

				repos := rm.GetRepos()
				Î©(repos).Should(haveLen(2))
				Î©(repos[0] == path.Join(baseDir, repoList[0])).Should(BeTrue())
				Î©(repos[1] == path.Join(baseDir, repoList[1])).Should(BeTrue())
			})
			It("Should get repo list successfully with non-git directories", func() {
				repoList = append(repoList, "dir-2")
				CreateDir(baseDir, repoList[1], true)
				CreateDir(baseDir, "Not-a-git-repo", false)
				rm, err := NewRepoManager(baseDir, repoList, true)
				Î©(err).Should(BeNil())

				repos := rm.GetRepos()
				Î©(repos).Should(HaveLen(2))
				Î©(repos[0] == path.Join(baseDir, repoList[0])).Should(BeTrue())
				Î©(repos[1] == path.Join(baseDir, repoList[1])).Should(BeTrue())
			})
			
			It("Should create branches successfully", func() {
				repoList = append(repoList, "dir-2")
				CreateDir (baseDir, repoList[1], true)
				rm, err := NewRepoManager(baseDir, repoList, true)
				Î©(err).Should(BeNil())

				output, err := rm.Exec("checkout -b test-branch")
				Î©(err).Should(BeNil())

				for _, out := range output {
					Î©(out).Should(Equal("Switched to a new branch 'test-branch'\n"))
				}
			})
			It("Should commit files successfully", func() {
				rm, err := NewRepoManager(baseDir, repoList, true)
				Î©(err).Should(BeNil())

				output, err := rm.Exec("checkout -b test-branch")
				Î©(err).Should(BeNil())

				for _, out := range output {
					Î©(out).Should(Equal("Switched to a new branch 'test-branch'\n"))
				}
				err = AddFiles(baseDir, repoList[0], true, "file_1.txt", "file_2.txt")
				Î©(err).Should(BeNil())

				wd, _ := os.Getwd()
				defer os.Chdir(wd)

				dir := path.Join(baseDir, repoList[0])
				err = os.Chdir(dir)
				Î©(err).Should(BeNil())

				Output, err = rm.Exec("log --oneline")
				Î©(err).Should(BeNil())

				ok := strings.HasSuffix(output[dir], "added some files...\n")
				Î©(ok).Should(BeTrue())
			})
		})
	})

