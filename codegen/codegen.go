//go:build codegen

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

//go:generate echo "Removing old go source files..."
//go:generate rm -rf ./clients
//go:generate rm -rf ../clients

//go:generate echo "Generating source files..."
//go:generate go run codegen.go

//go:generate echo "Running gofmt on source files..."
//go:generate gofmt -s -w .

//go:generate echo "Moving source files to repository root"
//go:generate mv -f clients/ ../

// add list of AWS service names to include in the code generation
// if empty, include everything
var include_services = []string{}

// add list of AWS servic names to exclude in the code generation
// if empty, exclude nothing
var exclude_services = []string{"internal"}

// set simple datastructure over
type set map[string]struct{}

func (s set) add(member string) {
	s[member] = struct{}{}
}

func (s set) contains(member string) bool {
	_, ok := s[member]
	return ok
}

func (s set) remove(member string) {
	delete(s, member)
}

func (s set) empty() bool {
	return len(s) == 0
}

// GitResponse
type GitResponse struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// fetchServicesList uses github API to GET the'repos/OWNER/REPO/contents/PATH' response to
// build a list of sub-directories under aws/aws-sdk-go-v2/service path.
// //github API reference: https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28
func fetchServicesList(owner, repo, path string) ([]string, error) {

	// build include map
	// service names in this set to be included only...
	include := make(set)
	for _, i := range include_services {
		include.add(i)
	}

	// build exclude map
	// service names in this list must be excluded..
	exclude := make(set)
	for _, i := range exclude_services {
		exclude.add(i)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/contents/%v", owner, repo, path)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)

	var resp []GitResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	contents := make([]string, 0)
	for _, r := range resp {
		// if the includes list is empty OR this service is part of the list OR
		if include.empty() || include.contains(r.Name) {
			if exclude.empty() || !exclude.contains(r.Name) {
				contents = append(contents, r.Name)
			}
		}
	}

	return contents, nil
}

// generateClients_go generates the clients.go
func generateClients_go(prefix, owner, repo, path string, services []string) error {

	// create the path /clients/_<service>
	relativePath := filepath.Join(".", "clients", fmt.Sprintf("_%v", prefix))
	err := createIfNotExists(relativePath)
	if err != nil {
		fmt.Println(fmt.Sprintf("error checking or creating sub-dir %v", relativePath))
		return err
	}

	// create files /clients/_<service>/client.go
	f, err := os.Create(filepath.Join(relativePath, "client.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	clients_go_Tmpl.Execute(f, TemplateData{
		GeneratedAt: timeStr(),
		Package:     prefix,
		Owner:       owner,
		Repo:        repo,
		Path:        path,
		Services:    services,
	})

	return nil
}

// timeStr returns the current timestamp as string
func timeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// createIfNotExists checks if the dir exists, else create
func createIfNotExists(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(fmt.Sprintf("Creating directory %v", dir))
			err := os.Mkdir(dir, 0777)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		fmt.Println(fmt.Sprintf("Directory exists: %v", dir))
	}
	return nil
}

type TemplateData struct {
	GeneratedAt string
	Package     string
	Owner       string
	Repo        string
	Path        string
	Services    []string
}

type formatPackageNameFunc func(string) string

var formatPackageName formatPackageNameFunc = func(name string) string {
	return fmt.Sprintf("_%v", name)
}

var clients_go_Tmpl = template.Must(template.New("client").Funcs(template.FuncMap{
	"formatPackageName": formatPackageName,
}).Parse(`
// AUTO-GENERATED CODE - DO NOT EDIT
// See instructions under /codegen/README.md
// GENERATED ON {{ .GeneratedAt }}

// Package {{ formatPackageName .Package }} provides AWS client management functions for the {{ .Package }} 
// AWS service. 
//
// The Client() is a wrapper on {{ .Package }}.NewFromConfig(), which creates & caches
// the client.
//
// The Delete() clears the cached client.
//
package {{ formatPackageName .Package }}

import (
	"sync"

	"github.com/TouchBistro/aws-ccp-go/providers"
	{{- range .Services}}
	"github.com/aws/aws-sdk-go-v2/service/{{ printf "%v" . }}"
	{{- end}}
)

var cmap sync.Map

{{ range $idx, $svc := .Services}}
// Client builds or returns the singleton {{ $svc }} client for the supplied provider
// If functional options are supplied, they are passed as-is to the underlying NewFromConfig(...)
// for the corresponding client
func Client(provider providers.CredsProvider, optFns ...func(*{{ $svc }}.Options)) (*{{ $svc }}.Client, error) {

	if provider == nil {
		return nil, providers.ErrNilProvider
	}
	if _, ok := cmap.Load(provider.Key()); !ok {
		client := {{ $svc }}.NewFromConfig(provider.Config(), optFns...)
		cmap.Store(provider.Key(),client)
	}
	client, _ := cmap.Load(provider.Key())
	return client.(*{{ $svc }}.Client), nil
}

// Must wraps the _{{ $svc }}.Client( ) function & panics if a non-nil error is returned. 
func Must(provider providers.CredsProvider, optFns ...func(*{{ $svc }}.Options)) *{{ $svc }}.Client {

	client, err := Client(provider,optFns...)
	if err != nil {
		panic(err)
	}
	return client
}

// Delete removes the cached {{ $svc }} client for the supplied provider; This foreces the subsequent
// calls to Client() for the same provider to recreate & return a new instnce.
func Delete(provider providers.CredsProvider) error {

	if provider == nil {
		return providers.ErrNilProvider
	}
	if _, ok := cmap.Load(provider.Key()); ok {
		cmap.Delete(provider.Key())
	}
	return nil
}


// Refresh discards the cached {{ $svc }} client if it exists, builds & returns a new singleton instance
func Refresh(provider providers.CredsProvider, optFns ...func(*{{ $svc }}.Options)) (*{{ $svc }}.Client, error) {

	err := Delete(provider)
	if err != nil {
		return nil, err
	}
	return Client(provider,optFns...)
}


{{ end }}
`))

func main() {
	owner := "aws"
	repo := "aws-sdk-go-v2"
	path := "service"

	contents, err := fetchServicesList(owner, repo, path)
	if err != nil {
		fmt.Println("Error fetching list of services from aws-go-sdk-v2 SDK")
		os.Exit(1)
	}

	//create clients sub-directory
	err = createIfNotExists(filepath.Join(".", "clients"))
	if err != nil {
		fmt.Println(fmt.Sprintf("Error checking or creating sub-dir clients"))
		os.Exit(1)
	}

	for _, content := range contents {
		contentsSubset := []string{content}
		if len(contentsSubset) > 0 {
			err := generateClients_go(content, owner, repo, path, contentsSubset)
			if err != nil {
				fmt.Println(fmt.Sprintf("Error generating client.go for %v", path))
			}
		}
	}
}
