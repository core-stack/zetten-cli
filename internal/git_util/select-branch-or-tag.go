package git_util

import (
	"errors"

	"github.com/core-stack/zetten-cli/internal/prompt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type SelectBranchOrTag struct {
	Branch string `help:"The branch to install" short:"b" long:"branch"`
	Tag    string `help:"The tag/version to install" short:"t" long:"tag"`
}

func SelectTag(repo *git.Repository) (string, error) {
	iterator, err := repo.Tags()
	if err != nil {
		return "", err
	}
	tags := ExtractTags(iterator)

	tag, err := prompt.PromptSelect("📝 Tag", tags, true)
	if err != nil {
		return "", err
	}
	return tag, nil
}

func SelectBranch(repo *git.Repository) (string, error) {
	branchs, err := repo.Branches()
	if err != nil {
		return "", err
	}
	branch, err := prompt.PromptSelect("📝 Branch", ExtractBranchs(branchs), true)
	if err != nil {
		return "", err
	}
	return branch, nil
}

func SelectTagOrBranch(repo *git.Repository) (string, error) {
	selected, err := prompt.PromptSelect("📝 Select tag or branch", []string{"tag", "branch"}, false)
	if err != nil {
		return "", err
	}
	var tagOrBranch string
	if selected == "tag" {
		tagOrBranch, err = SelectTag(repo)
	} else {
		tagOrBranch, err = SelectBranch(repo)
	}
	if err != nil && err != prompt.GoBack {
		return "", err
	}
	return tagOrBranch, nil
}

func LoadBranchOrTag(repo *git.Repository, branch, tag string) (string, error) {
	if tag == "" && branch == "" {
		return SelectTagOrBranch(repo)
	} else {
		if tag != "" {
			_, err := repo.Tag(tag)
			if err != nil {
				if err == plumbing.ErrReferenceNotFound {
					return SelectTagOrBranch(repo)
				}
				return "", err
			} else {
				return tag, nil
			}
		} else if branch != "" {
			_, err := repo.Branch(branch)
			if err != nil {
				if err == plumbing.ErrReferenceNotFound {
					SelectTagOrBranch(repo)
				}
				return "", err
			} else {
				return branch, nil
			}
		}
	}
	return "", errors.New("invalid tag or branch")
}
