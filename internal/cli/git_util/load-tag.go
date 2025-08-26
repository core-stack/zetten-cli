package git_util

import (
	"fmt"

	"github.com/core-stack/zetten-cli/internal/cli/prompt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func SelectTag(repo *git.Repository, canGoBack bool) (string, error) {
	iterator, err := repo.Tags()
	if err != nil {
		return "", err
	}
	tags := ExtractTags(iterator)

	tag, err := prompt.PromptSelect("📝 Tag", tags, canGoBack)
	if err != nil {
		return "", err
	}
	return tag, nil
}

func LoadTag(repo *git.Repository, tag string, canGoBack bool) (string, error) {
	if tag == "" {
		return SelectTag(repo, canGoBack)
	} else {
		_, err := repo.Tag(tag)
		if err != nil {
			if err == plumbing.ErrReferenceNotFound {
				fmt.Println("tag not found, selecting tag...")
				return SelectTag(repo, canGoBack)
			}
			return "", err
		} else {
			return tag, nil
		}
	}
}
