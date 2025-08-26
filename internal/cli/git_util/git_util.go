package git_util

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func RemoteRefExists(repoURL, refTarget string, auth transport.AuthMethod) (bool, error) {
	remote := git.NewRemote(nil, &gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})

	refs, err := remote.List(&git.ListOptions{Auth: auth})
	if err != nil {
		return false, fmt.Errorf("failed to list remote references: %w", err)
	}

	for _, ref := range refs {
		if ref.Name().String() == refTarget {
			return true, nil
		}
	}

	return false, nil
}

func ExtractBranchs(refs storer.ReferenceIter) []string {
	var branches []string
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			branches = append(branches, strings.TrimPrefix(ref.Name().String(), "refs/heads/"))
		}
		return nil
	})
	return branches
}
func ExtractTags(refs storer.ReferenceIter) []string {
	var tags []string
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsTag() {
			tags = append(tags, strings.TrimPrefix(ref.Name().String(), "refs/tags/"))
		}
		return nil
	})
	return tags
}
