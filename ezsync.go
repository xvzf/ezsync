package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
)

func getCloneURLandPath(repo *github.Repository, accesstoken string) (string, string) {
	githubCloneURL := repo.GetCloneURL()

	path := strings.Replace(githubCloneURL, "https://github.com/", "", -1)
	cloneURL := fmt.Sprintf("https://%s@github.com/%s", accesstoken, path)
	return cloneURL, path[:len(path)-4] // @TODO this is not perfect (just trim the trailing .git)
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: <tool> /path/to/backup")
		os.Exit(1)
	}
	backupPath := os.Args[1]

	accessToken := os.Getenv("GITHUB_ACCESSTOKEN")
	if accessToken == "" {
		fmt.Println("Please provide an access token for accessing github inside GITHUB_ACCESTOKEN")
		os.Exit(1)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	user, _, err := client.Users.Get(context.Background(), "")

	if err != nil {
		log.Fatal("Could not retrieve user information, check access token")
	}

	privateRepos, _, err := client.Repositories.List(context.Background(), "", &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: 1000, // This should be sufficient; @TODO Improve and actually make sure there are no repos left out
		},
	})

	if err != nil {
		log.Fatal("Could not retrieve private repositories, check security contexts of the access token")
	}

	publicRepos, _, err := client.Repositories.List(context.Background(), user.GetLogin(), &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: 1000, // This should be sufficient; @TODO Improve and actually make sure there are no repos left out
		},
	})
	if err != nil {
		log.Println(err)
		log.Fatal("Could not retrieve public repositories")
	}

	// All private  & public
	repos := append(privateRepos, publicRepos...)

	for _, repo := range repos {
		gitURL, path := getCloneURLandPath(repo, accessToken)
		log.Printf("[+] Syncing\t%s\n", path)

		// Try to clone the repo
		_, err := git.PlainClone(backupPath+path, false, &git.CloneOptions{
			URL: gitURL,
		})

		if err == git.ErrRepositoryAlreadyExists {
			// Try to pull changes
			repo, err := git.PlainOpen(backupPath + path)
			if err != nil {
				log.Printf("[!] Failed to open repo\t%s\n", path)
				continue
			}
			err = repo.Fetch(&git.FetchOptions{
				RemoteName: "origin",
			})

			// Merge master branch
			// workTree, err := repo.Worktree()
			// workTree.Pull(&git.PullOptions{
			// 	RemoteName: "origin",
			// })
			if err == git.NoErrAlreadyUpToDate {
				log.Printf("[+] Already up to date\t%s\n", path)
			} else if err != nil {
				log.Printf("[!] Failed to fetch repo\t%s\n", path)
				fmt.Println(err)
				continue
			} else {
				log.Printf("[+] Fetched remote\t%s\n", path)
				continue
			}
		} else if err != nil {
			log.Printf("[!] Failed to clone repo\t%s\n", path)
		} else {
			log.Printf("[+] Cloned\t%s\n", path)
		}
	}

}
