![Docker Image CI](https://github.com/relizaio/relizaGoClient/workflows/Docker%20Image%20CI/badge.svg?branch=master)

# Reliza Go Client

This tool allows streaming metadata about instances, releases and artifacts to [Reliza Hub at relizahub.com](https://relizahub.com) (currently in public preview mode). Community forum and support available at [r/Reliza](https://reddit.com/r/Reliza).

Docker image is available at [relizaio/reliza-go-client](https://hub.docker.com/r/relizaio/reliza-go-client)

## 1. Use Case: Get Version Assignment From Reliza Hub

Note that project schema must be preset on Reliza Hub prior to using this API. API key must also be generated for the project from Reliza Hub.

Sample command for semver version schema:

```
docker run --rm relizaio/reliza-go-client    \
    getversion    \
    -i project_api_id    \
    -k project_api_key    \
    -b master    \
    --pin 1.2.patch
```

Flags stand for:
- **getversion** - command that denotes we are obtaning next available release version for the branch. Note that if the call succeeds version assignment will be recorded and not given again by Reliza Hub, even if not consumed.
- **-i** - flag for project api id (required).
- **-k** - flag for project api key (required).
- **-b** - flag to denote branch (required). If branch is not recorded yet, Reliza Hub will attempt to create it.
- **--pin** - flag to denote branch pin. Required for new branches, optional for existing branches. If supplied for an existing branch and pin is different from current, it will override current pin.


## 2. Use Case: Send Release Metadata to Reliza Hub

As in previous case, API key must be generated for the project on Reliza Hub prior to sending release details.

Sample command to send release details:

```
docker run --rm relizaio/reliza-go-client    \
    addrelease    \
    -i project_api_id    \
    -k project_api_key    \
    -b master    \
    -v 20.02.3    \
    --vcsuri github.com/relizaio/relizaGoClient    \
    --vcstype git    \
    --commit 7bfc5ce7b0da277d139f7993f90761223fa54442    \
    --vcstag 20.02.3    \
    --artid relizaio/reliza-go-client    \
    --artbuildid 1    \
    --artcimeta Github Actions    \
    --arttype Docker    \
    --artdigests sha256:4e8b31b19ef16731a6f82410f9fb929da692aa97b71faeb1596c55fbf663dcdd
```

Flags stand for:
- **addrelease** - command that denotes we are sending Release Metadata of a Project to Reliza Hub.
- **-i** - flag for project api id (required).
- **-k** - flag for project api key (required).
- **-b** - flag to denote branch (required). If branch is not recorded yet, Reliza Hub will attempt to create it.
- **-v** - version (required). Note that Reliza Hub will reject the call if a release with this exact version is already present for this project.
- **vcsuri** - flag to denote vcs uri (optional). Currently this flag is needed if we want to set a commit for the release. However, soon it will be needed only if the vcs uri is not yet set for the project.
- **vcstype** - flag to denote vcs type (optional). Supported values: git, svn, mercurial. As with vcsuri, this flag is needed if we want to set a commit for the release. However, soon it will be needed only if the vcs uri is not yet set for the project.
- **commit** - flag to denote vcs commit id or hash (optional). This is needed to provide source code entry metadata into the release.
- **vcstag** - flag to denote vcs tag (optional). This is needed to include vcs tag into commit, if present.
- **artid** - flag to denote artifact identifier (optional). This is required to add artifact metadata into release.
- **artbuildid** - flag to denote artifact build id (optional). This flag is optional and may be used to indicate build system id of the release (i.e., this could be circleci build number).
- **artcimeta** - flag to denote artifact CI metadata (optional). This flag is optional and like artbuildid may be used to indicate build system metadata in free form.
- **arttype** - flag to denote artifact type (optional). This flag is used to denote artifact type. Currently supported values: Docker, JAR, WAR, Zip, Tar, Tar GZip.
- **artdigests** - flag to denote artifact digests (optional). This flag is used to indicate artifact digests. By convention, digests must be prefixed with type followed by colon and then actual digest hash, i.e. *sha256:4e8b31b19ef16731a6f82410f9fb929da692aa97b71faeb1596c55fbf663dcdd* - here type is *sha256* and digest is *4e8b31b19ef16731a6f82410f9fb929da692aa97b71faeb1596c55fbf663dcdd*. Multiple digests are supported and must be comma separated. I.e.: 
```
--artdigests sha256:4e8b31b19ef16731a6f82410f9fb929da692aa97b71faeb1596c55fbf663dcdd,sha1:fe4165996a41501715ea0662b6a906b55e34a2a1
```

Note that multiple artifacts per release are supported. In which case artifact specific flags (artid, arbuildid, artcimeta, arttype and artdigests must be repeated for each artifact).

For sample of how to use workflow in CI, refer to the GitHub Actions build yaml of this project [here](https://github.com/relizaio/relizaGoClient/blob/master/.github/workflows/dockerimage.yml).
