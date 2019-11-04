package config

import (
	"errors"
	"fmt"
	"github.com/buildtool/scaffold/pkg/config/ci"
	"github.com/buildtool/scaffold/pkg/config/vcs"
	"github.com/buildtool/scaffold/pkg/file"
	"github.com/buildtool/scaffold/pkg/stack"
	"github.com/buildtool/scaffold/pkg/templating"
	"github.com/caarlos0/env"
	"github.com/imdario/mergo"
	"github.com/liamg/tml"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	VCS          *VCSConfig `yaml:"vcs"`
	CI           *CIConfig  `yaml:"ci"`
	RegistryUrl  string     `yaml:"registry" env:"REGISTRY"`
	Organisation string     `yaml:"organisation"`
	CurrentCI    ci.CI
	CurrentVCS   vcs.VCS
}

type VCSConfig struct {
	Github *vcs.Github `yaml:"github"`
	Gitlab *vcs.Gitlab `yaml:"gitlab"`
}

type CIConfig struct {
	Buildkite *ci.Buildkite `yaml:"buildkite"`
	Gitlab    *ci.Gitlab    `yaml:"gitlab"`
}

func (c *Config) Configure() error {
	c.CurrentVCS.Configure()
	return c.CurrentCI.Configure()
}

func (c *Config) ValidateConfig() error {
	if c.CurrentVCS == nil {
		return errors.New("no VCS configured")
	}
	if c.CurrentCI == nil {
		return errors.New("no CI configured")
	}
	return nil
}

func Load(dir string, out io.Writer) (*Config, error) {
	cfg := InitEmptyConfig()

	err := parseConfigFiles(dir, out, func(dir string) error {
		return parseConfigFile(dir, cfg)
	})
	if err != nil {
		return nil, err
	}

	err = env.Parse(cfg)

	return cfg, err
}

func (c *Config) Validate(name string) error {
	if err := c.CurrentVCS.Validate(name); err != nil {
		return err
	}
	return c.CurrentCI.Validate(name)
}

func (c *Config) Scaffold(dir, name string, stack stack.Stack, out io.Writer) int {
	_, _ = fmt.Fprint(out, tml.Sprintf("<lightblue>Creating new service </lightblue><white><bold>'%s'</bold></white> <lightblue>using stack </lightblue><white><bold>'%s'</bold></white>\n", name, stack.Name()))
	_, _ = fmt.Fprint(out, tml.Sprintf("<lightblue>Creating repository at </lightblue><white><bold>'%s'</bold></white>\n", c.CurrentVCS.Name()))
	repository, err := c.CurrentVCS.Scaffold(name)
	if err != nil {
		_, _ = fmt.Fprintln(out, tml.Sprintf("<red>%s</red>", err.Error()))
		return -7
	}
	_, _ = fmt.Fprint(out, tml.Sprintf("<green>Created repository </green><white><bold>'%s'</bold></white>\n", repository.SSHURL))
	if err := c.CurrentVCS.Clone(dir, name, repository.SSHURL, out); err != nil {
		_, _ = fmt.Fprintln(out, tml.Sprintf("<red>%s</red>", err.Error()))
		return -8
	}
	projectDir := filepath.Join(dir, name)
	_, _ = fmt.Fprint(out, tml.Sprintf("<lightblue>Creating build pipeline for </lightblue><white><bold>'%s'</bold></white>\n", name))
	badges, err := c.CurrentCI.Badges(name)
	if err != nil {
		_, _ = fmt.Fprintln(out, tml.Sprintf("<red>%s</red>", err.Error()))
		return -9
	}
	parsedUrl, err := url.Parse(repository.HTTPURL)
	if err != nil {
		_, _ = fmt.Fprintln(out, tml.Sprintf("<red>%s</red>", err.Error()))
		return -10
	}
	data := templating.TemplateData{
		ProjectName:    name,
		Badges:         badges,
		Organisation:   c.Organisation,
		RegistryUrl:    c.RegistryUrl,
		RepositoryUrl:  repository.SSHURL,
		RepositoryHost: parsedUrl.Host,
		RepositoryPath: strings.Replace(parsedUrl.Path, ".git", "", 1),
	}
	webhook, err := c.CurrentCI.Scaffold(projectDir, data)
	if err != nil {
		_, _ = fmt.Fprintln(out, tml.Sprintf("<red>%s</red>", err.Error()))
		return -11
	}
	if err := addWebhook(name, webhook, c.CurrentVCS); err != nil {
		_, _ = fmt.Fprintln(out, tml.Sprintf("<red>%s</red>", err.Error()))
		return -12
	}
	if err := createDotfiles(projectDir); err != nil {
		_, _ = fmt.Fprintln(out, tml.Sprintf("<red>%s</red>", err.Error()))
		return -13
	}
	if err := createReadme(projectDir, data); err != nil {
		_, _ = fmt.Fprintln(out, tml.Sprintf("<red>%s</red>", err.Error()))
		return -14
	}
	if err := createDeployment(projectDir, data); err != nil {
		_, _ = fmt.Fprintln(out, tml.Sprintf("<red>%s</red>", err.Error()))
		return -15
	}
	if err := stack.Scaffold(projectDir, data); err != nil {
		_, _ = fmt.Fprintln(out, tml.Sprintf("<red>%s</red>", err.Error()))
		return -16
	}
	return 0
}

func InitEmptyConfig() *Config {
	return &Config{
		VCS: &VCSConfig{
			Github: &vcs.Github{},
			Gitlab: &vcs.Gitlab{},
		},
		CI: &CIConfig{
			Buildkite: &ci.Buildkite{},
			Gitlab:    &ci.Gitlab{},
		},
	}
}

func addWebhook(name string, url *string, vcs vcs.VCS) error {
	if url != nil {
		return vcs.Webhook(name, *url)
	}
	return nil
}

func createDotfiles(dir string) error {
	if err := file.Write(dir, ".gitignore", ""); err != nil {
		return err
	}
	editorconfig := `
root = true

[*]
end_of_line = lf
insert_final_newline = true
charset = utf-8
trim_trailing_whitespace = true
`
	if err := file.Write(dir, ".editorconfig", editorconfig); err != nil {
		return err
	}
	dockerignore := `
.git
.editorconfig
Dockerfile
README.md
`
	if err := file.Write(dir, ".dockerignore", dockerignore); err != nil {
		return err
	}
	return nil
}

func createReadme(dir string, data templating.TemplateData) error {
	content := `
| README.md
# {{.ProjectName}}
{{range .Badges}}[![{{.Title}}]({{.ImageUrl}})]({{.LinkUrl}}){{end}}
`
	return file.WriteTemplated(dir, "README.md", content, data)
}

func createDeployment(dir string, data templating.TemplateData) error {
	return file.WriteTemplated(dir, filepath.Join("k8s", "deploy.yaml"), deployment, data)
}

var deployment = `
apiVersion: apps/v1
kind: Deployment
metadata:
 labels:
   app: {{ .ProjectName }}
 name: {{ .ProjectName }}
 annotations:
   kubernetes.io/change-cause: "${TIMESTAMP} Deployed commit id: ${COMMIT}"
spec:
 replicas: 2
 selector:
   matchLabels:
     app: {{ .ProjectName }}
 strategy:
   rollingUpdate:
     maxSurge: 1
     maxUnavailable: 1
   type: RollingUpdate
 template:
   metadata:
     labels:
       app: {{ .ProjectName }}
   spec:
     affinity:
       podAntiAffinity:
         preferredDuringSchedulingIgnoredDuringExecution:
         - weight: 100
           podAffinityTerm:
             labelSelector:
               matchExpressions:
               - key: "app"
                 operator: In
                 values:
                 - {{ .ProjectName }}
             topologyKey: kubernetes.io/hostname
     containers:
     - name: {{ .ProjectName }}
       readinessProbe:
         httpGet:
           path: /
           port: 80
         initialDelaySeconds: 5
         periodSeconds: 5
         timeoutSeconds: 5
       imagePullPolicy: Always
       image: {{ .RegistryUrl }}/{{ .ProjectName }}:${COMMIT}
       ports:
       - containerPort: 80
     restartPolicy: Always
---

apiVersion: v1
kind: Service
metadata:
 name: {{ .ProjectName }}
spec:
 ports:
 - port: 80
   protocol: TCP
   targetPort: 80
 selector:
   app: {{ .ProjectName }}
 type: ClusterIP
`

var abs = filepath.Abs

func parseConfigFiles(dir string, out io.Writer, fn func(string) error) error {
	parent, err := abs(dir)
	if err != nil {
		return err
	}
	var files []string
	for parent != "/" {
		filename := filepath.Join(parent, ".scaffold.yaml")
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			files = append(files, filename)
		}

		parent = filepath.Dir(parent)
	}
	for i, file := range files {
		if i == 0 {
			_, _ = fmt.Fprintln(out, tml.Sprintf("Parsing config from file: <green>'%s'</green>", file))
		} else {
			_, _ = fmt.Fprintln(out, tml.Sprintf("Merging with config from file: <green>'%s'</green>", file))
		}
		if err := fn(file); err != nil {
			return err
		}
	}

	return nil
}

func parseConfigFile(filename string, cfg *Config) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return parseConfig(data, cfg)
}

func parseConfig(content []byte, config *Config) error {
	temp := &Config{}
	if err := yaml.UnmarshalStrict(content, temp); err != nil {
		return err
	} else {
		if err := mergo.Merge(config, temp); err != nil {
			return err
		}
		return validate(config)
	}
}

func validate(config *Config) error {
	elem := reflect.ValueOf(config.CI).Elem()
	for i := 0; i < elem.NumField(); i++ {
		currentCI := elem.Field(i).Interface().(ci.CI)
		if currentCI.ValidateConfig() == nil {
			if config.CurrentCI != nil && config.CurrentCI != currentCI {
				return fmt.Errorf("scaffold CI already defined, please check configuration")
			}
			config.CurrentCI = currentCI
		}
	}

	elem = reflect.ValueOf(config.VCS).Elem()
	for i := 0; i < elem.NumField(); i++ {
		currentVCS := elem.Field(i).Interface().(vcs.VCS)
		if currentVCS.ValidateConfig() == nil {
			if config.CurrentVCS != nil && config.CurrentVCS != currentVCS {
				return fmt.Errorf("scaffold VCS already defined, please check configuration")
			}
			config.CurrentVCS = currentVCS
		}
	}
	return nil
}
