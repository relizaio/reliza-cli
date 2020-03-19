![Docker Image CI](https://github.com/relizaio/relizaGoClient/workflows/Docker%20Image%20CI/badge.svg?branch=master)

# Reliza Go Client

This tool allows streaming metadata about instances, releases and artifacts to [Reliza Hub at relizahub.com](https://relizahub.com) (currently in public preview mode).

Playground instance is operational at [playground.relizahub.com](https://playground.relizahub.com). Video tutorial about key functionality available on [YouTube](https://www.youtube.com/watch?v=yDlf5fMBGuI).

Community forum and support available at [r/Reliza](https://reddit.com/r/Reliza).

Docker image is available at [relizaio/reliza-go-client](https://hub.docker.com/r/relizaio/reliza-go-client)

## 1. Use Case: Get Version Assignment From Reliza Hub

This use case requests Version from Reliza Hub for our project. Note that project schema must be preset on Reliza Hub prior to using this API. API key must also be generated for the project from Reliza Hub.

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
- **--pin** - flag to denote branch pin (optional for existing branches, required for new branches). If supplied for an existing branch and pin is different from current, it will override current pin.


## 2. Use Case: Send Release Metadata to Reliza Hub

This use case is commonly used in the CI workflow to stream Release metadata to Reliza Hub. As in previous case, API key must be generated for the project on Reliza Hub prior to sending release details.

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

## 3. Use Case: Check If Artifact Hash Already Present In Some Release

This is particularly useful for monorepos to see if there was a change in sub-project or not. We are using this technique in our sample [playground project](https://github.com/relizaio/reliza-hub-playground). We supply artifact hash to Reliza Hub - and if it's present already, we get release details; if not - we get empty json response {}. Search space is scoped to single project which is defined by API Id and API Key.

Sample command:

```
docker run --rm relizaio/reliza-go-client    \
    checkhash    \
    -i project_api_id    \
    -k project_api_key    \
    --hash sha256:hash
```

Flags stand for:
- **checkhash** - command that denotes we are checking artifact hash.
- **-i** - flag for project api id (required).
- **-k** - flag for project api key (required).
- **--hash** - flag to denote actual hash (required). By convention, hash must include hashing algorithm as its first part, i.e. sha256: or sha512:


## 4. Use Case: Send Deployment Metadata From Instance To Reliza Hub

This use case is for sending digests of active deployments from instance to Reliza Hub. API key must also be generated for the instance from Reliza Hub. Sample script is also provided in our [playground project](https://github.com/relizaio/reliza-hub-playground/blob/master/sample-instance-agent-scripts/send_instance_data.sh).

Sample command:

```
docker run --rm relizaio/reliza-go-client    \
    instdata    \
    -i instance_api_id    \
    -k instance_api_key    \
    --images "sha256:c10779b369c6f2638e4c7483a3ab06f13b3f57497154b092c87e1b15088027a5 sha256:e6c2bcd817beeb94f05eaca2ca2fce5c9a24dc29bde89fbf839b652824304703"   \
    --namespace default    \
    --sender sender1
```

Flags stand for:
- **instdata** - command that denotes we sending digest data from instance.
- **-i** - flag for instance api id (required).
- **-k** - flag for instance api key (required).
- **--images** - flag which lists sha256 digests of images sent from the instances (required). Images must be white space separated. Note that sending full docker image URIs with digests is also accepted, i.e. it's ok to send images as relizaio/reliza-go-client:latest@sha256:ebe68a0427bf88d748a4cad0a419392c75c867a216b70d4cd9ef68e8031fe7af
- **--namespace** - flag to denote namespace where we are sending images (optional, if not sent "default" namespace is used). Namespaces are useful to separate different products deployed on the same instance.
- **--sender** - flag to denote unique sender within a single namespace (optional). This is useful if say there are different nodes where each streams only part of application deployment data. In this case such nodes need to use same namespace but different senders so that their data does not stomp on each other.


## 5. Use Case: Request What Releases Must Be Deployed On This Instance From Reliza Hub

This use case is when instance queries Reliza Hub to receive infromation about what release versions and specific artifacts it needs to deploy. This would usually be used by either simple deployment scripts or full-scale CD systems. A sample use is presented in our [playground project script](https://github.com/relizaio/reliza-hub-playground/blob/master/sample-instance-agent-scripts/request_instance_target.sh).

Sample command:

```
docker run --rm relizaio/reliza-go-client    \
    getmyrelease    \
    -i instance_api_id    \
    -k instance_api_key    \
    --namespace default
```

Flags stand for:
- **getmyrelease** - command that denotes we are requesting release data for instance from Reliza Hub.
- **-i** - flag for instance api id (required).
- **-k** - flag for instance api key (required).
- **--namespace** - flag to denote namespace for which we are requesting release data (optional, if not sent "default" namespace is used). Namespaces are useful to separate different products deployed on the same instance.


## 6. Use Case: Request Latest Release Per Project Or Product

This use case is when Reliza Hub is queried either by CI or CD environment or by integration instance to check latest release version available per specific Project or Product.

Sample command:

```
docker run --rm relizaio/reliza-go-client    \
    getlatestrelease    \
    -i instance_api_id    \
    -k instance_api_key    \
    --project b4534a29-3309-4074-8a3a-34c92e1a168b    \
    --branch master
```

Flags stand for:
- **getlatestrelease** - command that denotes we are requesting latest release data for Project or Product from Reliza Hub
- **-i** - flag for api id which can be either api id for this project or organization-wide read API (required).
- **-k** - flag for api key which can be either api key for this project or organization-wide read API (required).
- **--project** - flag to denote UUID of specific Project or Product, UUID must be obtained from (Reliza Hub)[https://relizahub.com] (required).
- **--branch** - flag to denote required branch of chosen Project or Product (required).