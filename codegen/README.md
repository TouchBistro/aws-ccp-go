## codegen

`codegen.go` is a code generation utilty to create the `client.go` files for all service packages, i.e 
`github.com/aws-creds/provier/clients/_<service_name>`. These files contain the `Client()` builder functions which
construct, cache & return  the respective AWS SDK clients. The `<service_name>` here represents the AWS SDK service 
name as it is used by the [aws-sdk-go-v2](github.com/aws/aws-sdk-go-v2) module.

The main objective for having auto-generation is to simplify the maintenance of AWS Client builder functions. Since 
there are `244` services supported by the AWS SDK for Go v2 as of today, implementing & maintaining these builder 
functions by hand is cumbersome. Also if a logic change is required, it will be a waste of time & error prone to do 
that in all files.

<br>

## Design 
The `codegen` utility follows the 2 simple steps:

1. Use GitAPI to fetch the contents under the `github.com/aws/aws-sdk-go-v2/service` path. Each service supported by
AWS SDK has a sub-directory (package) under this location. The `internal` directory is ignored as it is not a client. 
Build a list of all services, by including and/or excluding the supplied list of services. see `include_services` 
and `exclude_services` variables to include or exclude specific clients if desired.
2. Create `/codegen/clients` sub-directory. For all services in the list from point 1, generate a `client.go` source 
file using a text template `client_go_Tmpl` for each service under `codegen/clients/_<service_name>/` sub-direcotries.

.: `codegen.go` has `go:generate` directives to perform the following steps:

1. Clean any previously generated code in the `/codegen/clients` directory
2. Run `go:generate` on `/codegen/codegen.go`
3. Run `gofmt` on all generated files under `/codegent/clients/`
4. Overwrite the existing files in the `clients` package with these new files.

<br>

## How to run 
The `codegen` utility is not run automatically during `build` or `ci` steps. So we can control when these auto-generated source 
files are to be updated.

To manually run the `codegen` utility, follow the steps below:

1. Pull latest code from the `master` branch
2. Create a new feature branch 
```
% git checkout -b feat/update-clients
```
3. Change to the `codegen` sub-directory & run `make` or `go:generate`
```
% cd codegen
% go generate codegen.go
```
4. If there are no errors reported, Three source files mentioned above will be updated.
5. Switch to the repository's root directory & build the code.
```bash
% go mod tidy
% make clean build
```
6. If successful, add & comit the changes & make a PR.
```
% git commit -am "updating aws service clients & builder methods"
% git push -u origin feat/update-clients
```
