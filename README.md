# ecr-image-sync

Downloads docker images from provided list, then tag with AWS ECR tags and uploads them to AWS ECR. This helps to keep
all images on AWS network (no need to download from public Internet), list can also be audited and reviewed.

ECR repository has to exist before image sync is run.

## build and run

If required, project can be built (`make build`) and run (`./ecr-image-sync`) locally:
```
ecr-image-sync [flags]
```
```
+--------------+----------------------------------------------------------------+
| flags                                                                         |
+--------------+----------------------------------------------------------------+
| -aws-account | AWS account where to sync images to                            |
+--------------+----------------------------------------------------------------+
| -aws-region  | AWS region where to sync images to                             |
+--------------+----------------------------------------------------------------+
| -dry-run     | whether to do sync - tag and push image                        |
+--------------+----------------------------------------------------------------+
| -images-file | file containing list of images to sync (default "images-list") |
+--------------+----------------------------------------------------------------+
| -version     | ecr-image-sync version                                         |
+--------------+----------------------------------------------------------------+
```

## build

`go build` or `go install`

## download

- [binary](https://github.com/pete911/ipcalc/releases)

## build/install

### brew

- add tap `brew tap pete911/tap`
- install `brew install pete911/tap/ipcalc`

### go

[go](https://golang.org/dl/) has to be installed.
- build `make build`
- install `make install`

### docker

- run (dry-run) `docker run --rm -v $HOME/.aws/credentials:/root/.aws/credentials:ro -v $(pwd)/images-list:/images-list:ro pete911/ecr-image-sync:latest -aws-account <account> -aws-region <region> -dry-run`
- exec inside container (debug) `docker run -it --rm --entrypoint /bin/sh pete911/ecr-image-sync:latest`

## release

Releases are published when the new tag is created e.g.
`git tag -m "add super cool feature" 0.1.0 && git push --follow-tags`
