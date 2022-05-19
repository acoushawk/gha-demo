// Copyright © 2021 Lucidworks
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func commitFilesToGithub(githubOrgName string, customerRepo string, commitBranch *string, commitMessage *string, filesToCommit []string, authorName *string, authorEmail *string, prSubject *string, prBranch *string, prDescription *string) {
	ref, err := getRef(githubOrgName, customerRepo, commitBranch)
	if err != nil {
		log.Fatalf("Unable to get/create the commit reference: %s\n", err)
	}
	if ref == nil {
		log.Fatalf("No error where returned but the reference is nil")
	}

	tree, err := getTree(ref, githubOrgName, customerRepo, filesToCommit)
	if err != nil {
		log.Fatalf("Unable to create the tree based on the provided files: %s\n", err)
	}

	if err := pushCommit(ref, tree, githubOrgName, customerRepo, commitMessage, authorName, authorEmail); err != nil {
		log.Fatalf("Unable to create the commit: %s\n", err)
	}

	if err := createPR(githubOrgName, customerRepo, prSubject, commitBranch, prBranch, prDescription); err != nil {
		log.Fatalf("Error while creating the pull request: %s", err)
	}
}

func githubClient() *github.Client {
	token := getSecretFromGCP("lw-cloud-infra-service-account_github_PAT")
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

// getRef returns the commit branch reference object if it exists or creates it
// from the base branch before returning it.
func getRef(sourceOwner string, sourceRepo string, commitBranch *string) (ref *github.Reference, err error) {
	if ref, _, err = githubClient().Git.GetRef(ctx, sourceOwner, sourceRepo, fmt.Sprintf("refs/heads/%s", *commitBranch)); err == nil {
		return ref, nil
	}

	if *commitBranch == "main" {
		return nil, errors.New("committing directly to the main branch is not allowed")
	}

	var baseRef *github.Reference
	if baseRef, _, err = githubClient().Git.GetRef(ctx, sourceOwner, sourceRepo, "refs/heads/main"); err != nil {
		return nil, err
	}
	newRef := &github.Reference{Ref: github.String(fmt.Sprintf("refs/heads/%s", *commitBranch)), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err = githubClient().Git.CreateRef(ctx, sourceOwner, sourceRepo, newRef)
	return ref, err
}

// getTree generates the tree to commit based on the given files and the commit
// of the ref you got in getRef.
func getTree(ref *github.Reference, sourceOwner string, sourceRepo string, sourceFiles []string) (tree *github.Tree, err error) {
	// Create a tree with what to commit.
	entries := []*github.TreeEntry{}

	// Load each file into the tree.
	for _, fileArg := range sourceFiles {
		file, content, err := getFileContent(fileArg)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &github.TreeEntry{Path: github.String(file), Type: github.String("blob"), Content: github.String(string(content)), Mode: github.String("100644")})
	}

	tree, _, err = githubClient().Git.CreateTree(ctx, sourceOwner, sourceRepo, *ref.Object.SHA, entries)
	return tree, err
}

// getFileContent loads the local content of a file and return the target name
// of the file in the target repository and its contents.
func getFileContent(fileArg string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty `files` parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = ioutil.ReadFile(localFile)
	return targetName, b, err
}

// pushCommit creates the commit in the given reference using the given tree.
func pushCommit(ref *github.Reference, tree *github.Tree, sourceOwner string, sourceRepo string, commitMessage *string, authorName *string, authorEmail *string) (err error) {
	// Get the parent commit to attach the commit to.
	parent, _, err := githubClient().Repositories.GetCommit(ctx, sourceOwner, sourceRepo, *ref.Object.SHA, nil)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: authorName, Email: authorEmail}
	commit := &github.Commit{Author: author, Message: commitMessage, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	newCommit, _, err := githubClient().Git.CreateCommit(ctx, sourceOwner, sourceRepo, commit)
	if err != nil {
		return err
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = githubClient().Git.UpdateRef(ctx, sourceOwner, sourceRepo, ref, false)
	return err
}

// createPR creates a pull request. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func createPR(sourceOwner string, sourceRepo string, prSubject *string, commitBranch *string, prBranch *string, prDescription *string) (err error) {
	// if prSubject == "" {
	// 	return errors.New("missing `-pr-title` flag; skipping PR creation")
	// }

	newPR := &github.NewPullRequest{
		Title:               prSubject,
		Head:                commitBranch,
		Base:                prBranch,
		Body:                prDescription,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := githubClient().PullRequests.Create(ctx, sourceOwner, sourceRepo, newPR)
	if err != nil {
		return err
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}
