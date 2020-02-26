![Docker Image CI](https://github.com/relizaio/relizaGoClient/workflows/Docker%20Image%20CI/badge.svg?branch=master)

# Reliza Go Client

This tool allows streaming metadata about instances, releases and artifacts to [Reliza Hub at relizahub.com](https://relizahub.com) (currently in public preview mode). Community forum and support available at [r/Reliza](https://reddit.com/r/Reliza).

Docker image is available at [relizaio/reliza-go-client](https://hub.docker.com/r/relizaio/reliza-go-client)

## Get Version Assignment Use Case

Note that project schema must be preset in any case.

Sample command for semver version schema:

```
docker run --rm relizaio/reliza-go-client getversion -i project_api_id -k project_api_key -b master --pin 1.2.patch
```

Flags stand for:
- getversion - command that denotes we are obtaning next available release version for the branch. Note that if the call succeeds version assignment will be recorded and not given again by Reliza Hub, even if not consumed.
- -i - flag for project api id (required).
- -k - flag for project api key (required).
- -b - flag to denote branch (required). If branch is not recorded yet, Reliza Hub will attempt to create it.
- --pin - flag to denote branch pin. Required for new branches, optional for existing branches. If supplied for an existing branch and pin is different from current, it will override current pin.