package github

import (
	"context"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/vendor4git"
	"golang.org/x/oauth2"
)

type gitHubVendor struct {
	client *github.Client
}

func NewGitHubVendor(ctx context.Context, accessToken string) vendor4git.Vendor {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return NewGithubVendorWithClient(ctx, client)
}

func NewGithubVendorWithClient(ctx context.Context, client *github.Client) vendor4git.Vendor {
	return &gitHubVendor{client: client}
}

func (x *gitHubVendor) CreateRepo(ctx context.Context, name string, owner string, private bool) (*vendor4git.Repository, error) {

	repo := &github.Repository{
		Name:    github.String(name),
		Private: github.Bool(private),
	}
	repo, _, err := x.client.Repositories.Create(ctx, owner, repo)
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

	_, err := x.client.Repositories.Delete(ctx, owner, name)
	ghErr, ok := err.(*github.ErrorResponse)
	if ghErr != nil && ok && ghErr.Response.StatusCode == 404 {
		return vendor4git.ErrRepoNotFound
	}

	return err
}
