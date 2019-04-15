# quay-scheduled-build

deploy an app to enforce scheduled builds on quay.io

## overview

There is not a straight forward way to have a container built on quay based on a schedule.

Only by:
- UI clicking to trigger a build
- git trigger (ideal for source, but not ideal if the repo has a fairly-static Dockerfile)
- the swagger API

This simple service uses the swagger API, and attempts to make this process easier.

## requirements

- a quay account
- a robot account for the org (with 'write' permissions for the repo you want to be able to build)
- an 'oauth app' token, with at least the permissions:
  - Administer Organization (`org:admin`)
  - Administer Repositories (`repo:admin`)
  - Create Repositories (`repo:create`)
  - View all visible repositories (`repo:read`)
  - Read/Write to any accessible repositories (`repo:write`)

IMO these permissions are way too excessive, but it's the minimum set I could get working (as of 2019-04).

## config

```json
{
  "builds": [
    {
      "quay_repo": "ohman/buildah",
      "schedule": "* */1 * * *",
      "token": "asdfasdfasdfasdfasdfasdf",
      "pull_robot": "ohman+dat_robot_tho",
      "tags": [
        "master",
	"latest"
      ],
      "ref": {
        "archive_url": "https://github.com/knative/build-templates/archive/master.tar.gz",
	"dockerfile_path": "/build-templates-master/buildah/Dockerfile",
	"subdirectory": "/build-templates-master/buildah/",
        "context": "/build-templates-master/buildah/"
      }
    }
  ]
}
```
