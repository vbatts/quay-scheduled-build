# quay-scheduled-build

[![Docker Repository on Quay](https://quay.io/repository/vbatts/quay-scheduled-build/status "Docker Repository on Quay")](https://quay.io/repository/vbatts/quay-scheduled-build)

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

- https://quay.io/ to sign up for your account (there is a free plan)
- https://docs.quay.io/api/ has some information on setting up Oauth2 tokens
- https://docs.quay.io/glossary/robot-accounts.html on robot accounts

## container available

There is a container build of this project available at https://quay.io/repository/vbatts/quay-scheduled-build

You can `docker|podman pull quay.io/vbatts/quay-scheduled-build`.

Use this to generate your config, run oneshot or as a server.


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

It accepts both flags or environment variables:
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

## oneshot

You may just want to kick the build(s) imediately.
So `quay-scheduled-build oneshot` uses the same config file, but does not need the `schedule` field.

```shell
NAME:
   main oneshot - trigger the builds right meow

USAGE:
   main oneshot [command options] [arguments...]

DESCRIPTION:
   trigger your builds on quay.io right meow.
  If there are multiple builds in your config, they are triggered in serial.

OPTIONS:
   --config value  build config to manage (default: "quay-build.json")
```

```shell
INFO[0000] reading config from "quay-build.json"        
INFO[0000] requesting imediate build of "ohman/buildah" 
INFO[0000] {"archive_url":"https://github.com/knative/build-templates/archive/master.tar.gz","context":"/build-templates-master/buildah","display_name":"b778aca","dockerfile_path":"/build-templates-master/buildah/Dockerfile","error":null,"id":"075e22ce-d6fa-4b19-9d5a-10c5b1f88a6e","is_writer":true,"manual_user":"vbatts","phase":"waiting","pull_robot":{"is_robot":true,"kind":"user","name":"ohman+buildahbot"},"repository":{"name":"buildah","namespace":"ohman"},"resource_key":null,"started":"Tue, 16 Apr 2019 17:51:26 -0000","status":{},"subdirectory":"/build-templates-master/buildah/Dockerfile","tags":["latest","master"],"trigger":null,"trigger_metadata":{}} 
```

## serve

Run as a daemon to trigger the container builds on quay based on their own schedule.
See https://godoc.org/github.com/robfig/cron for particulars on the `schedule` syntax.

```shell
NAME:
   main serve - the build scheduler

USAGE:
   main serve [command options] [arguments...]

DESCRIPTION:
   run the scheduler for your builds on quay.io

OPTIONS:
   --config value  build config to manage (default: "quay-build.json")
```

```shell
INFO[0000] readying the scheduler ...                   
INFO[0000] queuing build of "ohman/buildah" for "@weekly" 
INFO[0000] running the build scheduler ...              
INFO[0299] {"archive_url":"https://github.com/knative/build-templates/archive/master.tar.gz","context":"/build-templates-master/buildah","display_name":"b778aca","dockerfile_path":"/build-templates-master/buildah/Dockerfile","error":null,"id":"c5d94e76-dc0f-4d3a-b8e7-0123c3ddc31d","is_writer":true,"manual_user":"vbatts","phase":"waiting","pull_robot":{"is_robot":true,"kind":"user","name":"ohman+buildahbot"},"repository":{"name":"buildah","namespace":"ohman"},"resource_key":null,"started":"Tue, 16 Apr 2019 17:48:09 -0000","status":{},"subdirectory":"/build-templates-master/buildah/Dockerfile","tags":["latest","master"],"trigger":null,"trigger_metadata":{}}
```

## environment variables

Setting `BUILD_COMMAND=` environment variable is useful for running the container image (`quay.io/vbatts/quay-scheduled-build`), since it has a fixed entrypoint.
The values supported here:
- `serve`
- `oneshot`

Further, each of the sub commands have environment variables for their flag values.
