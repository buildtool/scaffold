package vcs

import (
	"context"
	"fmt"
	"github.com/buildtool/scaffold/pkg/config/vcs/mocks"
	"github.com/buildtool/scaffold/pkg/wrappers"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v28/github"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGithub_Validate(t *testing.T) {
	vcs := &Github{}

	err := vcs.Validate("project")

	assert.NoError(t, err)
}

func TestGithub_Configure(t *testing.T) {
	vcs := &Github{}

	vcs.Configure()

	assert.NotNil(t, vcs.repositories)
}

func TestGithubVCS_Scaffold(t *testing.T) {
	repoName := "reponame"
	repoSSHUrl := "cloneurl"
	repoCloneUrl := "https://github.com/example/repo"
	orgName := "org"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepositoriesService(ctrl)
	git := Github{
		Organisation: orgName,
		repositories: m,
	}

	repository := github.Repository{
		Name:     wrappers.String(repoName),
		AutoInit: wrappers.Bool(true),
		Private:  wrappers.Bool(false),
	}

	repositoryResponse := repository
	repositoryResponse.SSHURL = wrappers.String(repoSSHUrl)
	repositoryResponse.CloneURL = wrappers.String(repoCloneUrl)

	m.EXPECT().
		Create(context.Background(), orgName, &repository).Return(
		&repositoryResponse, githubCreatedResponse, nil).
		Times(1)

	m.EXPECT().
		UpdateBranchProtection(context.Background(), orgName, repoName, "master", &github.ProtectionRequest{
			EnforceAdmins: true,
			RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
				DismissStaleReviews:          true,
				RequiredApprovingReviewCount: 1,
			},
		}).Return(nil, githubOkResponse, nil).
		Times(1)

	res, err := git.Scaffold(repoName)
	assert.NoError(t, err)
	assert.Equal(t, &RepositoryInfo{repoSSHUrl, repoCloneUrl}, res)
}

func TestGithubVCS_ScaffoldWithoutOrganisation(t *testing.T) {
	repoName := "reponame"
	repoSSHUrl := "cloneurl"
	repoCloneUrl := "https://github.com/example/repo"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepositoriesService(ctrl)
	git := Github{repositories: m}

	repository := github.Repository{
		Name:     wrappers.String(repoName),
		AutoInit: wrappers.Bool(true),
		Private:  wrappers.Bool(false),
	}

	repositoryResponse := repository
	repositoryResponse.SSHURL = wrappers.String(repoSSHUrl)
	repositoryResponse.CloneURL = wrappers.String(repoCloneUrl)
	repositoryResponse.Owner = &github.User{
		Login: wrappers.String("user-login"),
	}

	m.EXPECT().
		Create(context.Background(), "", &repository).Return(
		&repositoryResponse, githubCreatedResponse, nil).
		Times(1)

	m.EXPECT().
		UpdateBranchProtection(context.Background(), "user-login", repoName, "master", &github.ProtectionRequest{
			EnforceAdmins: true,
			RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
				DismissStaleReviews:          true,
				RequiredApprovingReviewCount: 1,
			},
		}).Return(nil, githubOkResponse, nil).
		Times(1)

	res, err := git.Scaffold(repoName)
	assert.NoError(t, err)
	assert.Equal(t, &RepositoryInfo{repoSSHUrl, repoCloneUrl}, res)

}

func TestGithubVCS_Scaffold_RepositoryAlreadyExist(t *testing.T) {
	// TODO In reality we get this error for any response that is not http.StatusCreated
	repoName := "ALREADY_EXISTS"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepositoriesService(ctrl)
	git := Github{repositories: m}

	repository := github.Repository{
		Name:     wrappers.String(repoName),
		AutoInit: wrappers.Bool(true),
		Private:  wrappers.Bool(false),
	}

	repositoryResponse := repository
	repositoryResponse.Owner = &github.User{
		Login: wrappers.String("user-login"),
	}

	m.EXPECT().
		Create(context.Background(), "", &repository).Return(&repositoryResponse,
		&github.Response{
			Response: &http.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Status:     "already exists",
			},
		}, nil).
		Times(1)

	_, err := git.Scaffold(repoName)
	assert.EqualError(t, err, "failed to create repository ALREADY_EXISTS, already exists")
}

func TestGithubVCS_Scaffold_CreateError(t *testing.T) {
	repoName := "ALREADY_EXISTS"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepositoriesService(ctrl)
	git := Github{
		Organisation: "org",
		repositories: m,
	}

	repository := &github.Repository{
		Name:     wrappers.String(repoName),
		AutoInit: wrappers.Bool(true),
		Private:  wrappers.Bool(false),
	}

	m.EXPECT().
		Create(context.Background(), "org", repository).Return(
		repository, nil, fmt.Errorf("failed to create repo")).
		Times(1)

	_, err := git.Scaffold(repoName)
	assert.EqualError(t, err, "failed to create repo")
}

func TestGithubVCS_ScaffoldProtectBranchError(t *testing.T) {
	repoName := "reponame"
	repoCloneUrl := "cloneurl"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepositoriesService(ctrl)
	git := Github{repositories: m}

	repository := github.Repository{
		Name:     wrappers.String(repoName),
		AutoInit: wrappers.Bool(true),
		Private:  wrappers.Bool(false),
	}

	repositoryResponse := repository
	repositoryResponse.SSHURL = wrappers.String(repoCloneUrl)
	repositoryResponse.Owner = &github.User{
		Login: wrappers.String("user-login"),
	}

	m.EXPECT().
		Create(context.Background(), "", &repository).Return(
		&repositoryResponse, githubCreatedResponse, nil).
		Times(1)

	m.EXPECT().
		UpdateBranchProtection(context.Background(), "user-login", repoName, "master", &github.ProtectionRequest{
			EnforceAdmins: true,
			RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
				DismissStaleReviews:          true,
				RequiredApprovingReviewCount: 1,
			},
		}).Return(
		nil, githubBadRequestResponse, nil).
		Times(1)

	_, err := git.Scaffold(repoName)
	assert.EqualError(t, err, "failed to set repository branch protection something went wrong")
}

func TestGithubVCS_SillyTests(t *testing.T) {
	githubVCS := Github{}
	assert.EqualErrorf(t, githubVCS.ValidateConfig(), "token is required", "")
	githubVCS.Token = ""
	assert.EqualErrorf(t, githubVCS.ValidateConfig(), "token is required", "")

	githubVCS.Token = "token"
	assert.NoError(t, githubVCS.ValidateConfig())

	assert.Equal(t, githubVCS.Name(), "Github")
}

func TestGithubVCS_Webhook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepositoriesService(ctrl)
	githubVCS := Github{
		repoOwner:    "test",
		repositories: m,
	}
	m.EXPECT().CreateHook(context.Background(), "test", "repo", &github.Hook{
		Events: []string{
			"push",
			"pull_request",
			"deployment",
		},
		Config: map[string]interface{}{
			"url":          "https://ab.cd",
			"content_type": "json",
		},
		Active: wrappers.Bool(true),
	}).Return(nil, githubCreatedResponse, nil).
		Times(1)

	err := githubVCS.Webhook("repo", "https://ab.cd")
	assert.NoError(t, err)
}

func TestGithubVCS_WebhookError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepositoriesService(ctrl)
	githubVCS := Github{
		repoOwner:    "test",
		repositories: m,
	}
	m.EXPECT().CreateHook(context.Background(), "test", "repo", &github.Hook{
		Events: []string{
			"push",
			"pull_request",
			"deployment",
		},
		Config: map[string]interface{}{
			"url":          "https://ab.cd",
			"content_type": "json",
		},
		Active: wrappers.Bool(true),
	}).Return(nil, githubBadRequestResponse, nil).
		Times(1)

	err := githubVCS.Webhook("repo", "https://ab.cd")
	assert.EqualError(t, err, "failed to create webhook something went wrong")
}

var githubOkResponse = &github.Response{
	Response: &http.Response{
		StatusCode: http.StatusOK,
	},
}

var githubCreatedResponse = &github.Response{
	Response: &http.Response{
		StatusCode: http.StatusCreated,
	},
}

var githubBadRequestResponse = &github.Response{
	Response: &http.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "something went wrong",
	},
}
