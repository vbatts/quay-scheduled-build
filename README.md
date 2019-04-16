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

There is a config helper:
```shell
NAME:
   quay-scheduled-build generate config - helper to get a new buildref configuration

USAGE:
   quay-scheduled-build generate config [command options] [arguments...]

OPTIONS:
   --token value            quay.io oauth token for the build [$BUILD_TOKEN]
   --repo value             quay.io container repo for the build [$BUILD_REPO]
   --schedule value         cron style schedule to trigger this containter build (more info https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format) (default: "* * */1 * *") [$BUILD_SCHEDULE]
   --tags value             container name tags to apply to this build
   --robot value            quay.io robot account username for the build [$BUILD_ROBOT]
   --archive-url value      URL to the source of the build (which includes the Dockerfile) [$BUILD_ARCHIVE_URL]
   --dockerfile-path value  path (within the source archive) to the Dockerfile [$BUILD_DOCKERFILE_PATH]
   --subdirectory value     path (within the source archive) to the root of the build directory [$BUILD_SUBDIRECTORY]
```

So it accepts both flags or environment variables:
```shell
BUILD_ARCHIVE_URL="https://github.com/knative/build-templates/archive/master.tar.gz" \
  BUILD_DOCKERFILE_PATH="/build-templates-master/buildah/Dockerfile" \
  quay-scheduled-build \
  generate config \
  --token=asdfasdfasdf \
  --repo=ohman/buildah \
  --tags="latest" \
  --tags=master \
  --robot="ohman+dat_robot_tho"
```

```json
{
  "builds": [
    {
      "quay_repo": "ohman/buildah",
      "schedule": "* * */1 * *",
      "token": "asdfasdfasdf",
      "pull_robot": "ohman+dat_robot_tho",
      "docker_tags": [
        "latest",
        "master"
      ],
      "archive_url": "https://github.com/knative/build-templates/archive/master.tar.gz",
      "dockerfile_path": "/build-templates-master/buildah/Dockerfile"
    }
  ]
}
```

