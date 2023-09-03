package github

import (
	"context"

	"github.com/google/go-github/v54/github"
	"github.com/gov4git/vendor4git"
	"golang.org/x/oauth2"
)

type gitHubVendor struct {
	accessToken string
}

func NewGitHubVendor(accessToken string) vendor4git.Vendor {
	return &gitHubVendor{accessToken: accessToken}
}

func (x *gitHubVendor) CreateRepo(ctx context.Context, name string, owner string, private bool) (*vendor4git.Repository, error) {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: x.accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repo := &github.Repository{
		Name:    github.String(name),
		Private: github.Bool(private),
	}
	repo, _, err := client.Repositories.Create(ctx, owner, repo)
	errResp, ok := err.(*github.ErrorResponse)
	if ok && errResp.Response.StatusCode == 422 {
		return nil, vendor4git.ErrRepoExists
	}
	if err != nil {
		return nil, err
	}

	return &vendor4git.Repository{
		HTTPSURL: repo.GetCloneURL(),
		SSHURL:   repo.GetSSHURL(),
	}, nil
}

func (x *gitHubVendor) RemoveRepo(ctx context.Context, name string, owner string) error {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: x.accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	_, err := client.Repositories.Delete(ctx, owner, name)
	ghErr, ok := err.(*github.ErrorResponse)
	if ghErr != nil && ok && ghErr.Response.StatusCode == 404 {
		return vendor4git.ErrRepoNotFound
	}

	return err
}
