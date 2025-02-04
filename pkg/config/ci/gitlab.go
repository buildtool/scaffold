package ci

import (
	"errors"
	"fmt"
	"github.com/buildtool/scaffold/pkg/file"
	"github.com/buildtool/scaffold/pkg/templating"
	"github.com/xanzy/go-gitlab"
	"path/filepath"
	"strings"
)

type Gitlab struct {
	Group           string `yaml:"group" env:"GITLAB_GROUP"`
	Token           string `yaml:"token" env:"GITLAB_TOKEN"`
	badgesService   badgesService
	usersService    usersService
	groupsService   groupsService
	projectsService projectsService
}

type badgesService interface {
	ListProjectBadges(pid interface{}, opt *gitlab.ListProjectBadgesOptions, options ...gitlab.OptionFunc) ([]*gitlab.ProjectBadge, *gitlab.Response, error)
}

type usersService interface {
	CurrentUser(options ...gitlab.OptionFunc) (*gitlab.User, *gitlab.Response, error)
}

type projectsService interface {
	GetProject(pid interface{}, opt *gitlab.GetProjectOptions, options ...gitlab.OptionFunc) (*gitlab.Project, *gitlab.Response, error)
}

type groupsService interface {
	GetGroup(gid interface{}, options ...gitlab.OptionFunc) (*gitlab.Group, *gitlab.Response, error)
}

var _ CI = &Gitlab{}

func (c *Gitlab) Name() string {
	return "Gitlab"
}

func (c *Gitlab) ValidateConfig() error {
	if len(c.Token) == 0 {
		return errors.New("token for Gitlab not configured")
	}
	return nil
}

func (c *Gitlab) Validate(name string) error {
	_, _, err := c.usersService.CurrentUser()
	if err != nil {
		return err
	}
	_, _, err = c.groupsService.GetGroup(c.Group)
	if err != nil {
		return err
	}
	path := filepath.Join(c.Group, name)
	project, response, err := c.projectsService.GetProject(path, nil)
	if err != nil {
		if response == nil || response.StatusCode != 404 {
			return err
		}
	}
	if project != nil {
		return fmt.Errorf("project named '%s/%s' already exists at Gitlab", c.Group, name)
	}
	return nil
}

func (c *Gitlab) Scaffold(dir string, data templating.TemplateData) (*string, error) {
	if err := file.WriteTemplated(dir, ".gitlab-ci.yml", gitlabCiYml, data); err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *Gitlab) Badges(name string) ([]templating.Badge, error) {
	path := filepath.Join(c.Group, name)

	badges, _, err := c.badgesService.ListProjectBadges(path, nil)
	if err != nil {
		return nil, err
	}
	result := make([]templating.Badge, len(badges))
	for i, b := range badges {
		title := ""
		if strings.Contains(b.ImageURL, "build") {
			title = "Build status"
		} else if strings.Contains(b.ImageURL, "coverage") {
			title = "Coverage report"
		}
		result[i] = templating.Badge{
			Title:    title,
			ImageUrl: b.RenderedImageURL,
			LinkUrl:  b.RenderedLinkURL,
		}
	}
	return result, nil
}

func (c *Gitlab) Configure() error {
	git := gitlab.NewClient(nil, c.Token)
	c.badgesService = git.ProjectBadges
	c.usersService = git.Users
	c.groupsService = git.Groups
	c.projectsService = git.Projects
	return nil
}

var gitlabCiYml = `
stages:
  - build
  - deploy-staging
  - deploy-prod

variables:
  DOCKER_HOST: tcp://docker:2375/

image: buildtool/build-tools

build:
  stage: build
  services:
    - docker:dind
  script:
  - build
  - push

deploy-to-staging:
  stage: deploy-staging
  when: on_success
  script:
    - echo Deploy {{ .ProjectName }} to staging.
    - deploy staging
  environment:
    name: staging

deploy-to-prod:
  stage: deploy-prod
  when: on_success
  script:
    - echo Deploy {{ .ProjectName }} to prod.
    - deploy prod
  environment:
    name: prod
  only:
    - master
`
