![Docker Image CI](https://github.com/relizaio/reliza-cli/actions/workflows/dockerimage.yml/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/relizaio/reliza-cli)](https://goreportcard.com/report/github.com/relizaio/reliza-cli)
# Reliza CLI

This tool allows for command-line interactions with [Reliza Hub at relizahub.com](https://relizahub.com) (currently in public preview mode). Particularly, Reliza CLI can stream metadata about instances, releases, artifacts, resolve bundles based on Reliza Hub data. Available as either a Docker image or binary.

Video tutorial about key functionality of Reliza Hub is available on [YouTube](https://www.youtube.com/watch?v=yDlf5fMBGuI).

Argo CD GitOps Integration using Kustomize [tutorial](https://itnext.io/building-kubernetes-cicd-pipeline-with-github-actions-argocd-and-reliza-hub-e7120b9be870).

Community forum and support is available at [r/Reliza](https://reddit.com/r/Reliza).

Docker image is available at [relizaio/reliza-cli](https://hub.docker.com/r/relizaio/reliza-cli)

## Download Reliza CLI

Below are the available downloads for the latest version of the Reliza CLI (2024.07.6). Please download the proper package for your operating system and architecture.

The CLI is distributed as a single binary. Install by unzipping it and moving it to a directory included in your system's PATH.

[SHA256 checksums](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/sha256sums.txt)

macOS: [64-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/reliza-cli-2024.07.6-darwin-amd64.zip)

FreeBSD: [32-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/reliza-cli-2024.07.6-freebsd-386.zip) | [64-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/reliza-cli-2024.07.6-freebsd-amd64.zip) | [Arm](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/reliza-cli-2024.07.6-freebsd-arm.zip)

Linux: [32-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/reliza-cli-2024.07.6-linux-386.zip) | [64-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/reliza-cli-2024.07.6-linux-amd64.zip) | [Arm](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/reliza-cli-2024.07.6-linux-arm.zip) | [Arm64](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/reliza-cli-2024.07.6-linux-arm64.zip)

OpenBSD: [32-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/reliza-cli-2024.07.6-openbsd-386.zip) | [64-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/reliza-cli-2024.07.6-openbsd-amd64.zip)

Solaris: [64-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/reliza-cli-2024.07.6-solaris-amd64.zip)

Windows: [32-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/reliza-cli-2024.07.6-windows-386.zip) | [64-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2024.07.6/reliza-cli-2024.07.6-windows-amd64.zip)

It is possible to set authentication data via explicit flags, login command (see below) or following environment variables:

- APIKEYID - for API Key ID
- APIKEY - for API Key itself
- URI - for Reliza Hub Uri (if not set, default at https://app.relizahub.com is used)

# Table of Contents - Use Cases
1. [Get Version Assignment From Reliza Hub](#1-use-case-get-version-assignment-from-reliza-hub)
2. [Send Release Metadata to Reliza Hub](#2-use-case-send-release-metadata-to-reliza-hub)
3. [Check If Artifact Hash Already Present In Some Release](#3-use-case-check-if-artifact-hash-already-present-in-some-release)
4. [Send Deployment Metadata From Instance To Reliza Hub](#4-use-case-send-deployment-metadata-from-instance-to-reliza-hub)
5. *Deprecated* [Request What Releases Must Be Deployed On This Instance From Reliza Hub](#5-use-case-request-what-releases-must-be-deployed-on-this-instance-from-reliza-hub)
6. [Request Latest Release Per Project Or Product](#6-use-case-request-latest-release-per-project-or-product)
7. GitOps Operations:
    1. *Deprecated* [Parse Deployment Templates To Inject Correct Artifacts For GitOps](#71-use-case-parse-deployment-templates-to-inject-correct-artifacts-for-gitops)
    2. [Replace Tags On Deployment Templates To Inject Correct Artifacts For GitOps Using Instance And Revision](#72-use-case-replace-tags-on-deployment-templates-to-inject-correct-artifacts-for-gitops-using-instance-and-revision)
    3. [Replace Tags On Deployment Templates To Inject Correct Artifacts For GitOps Using Bundle](#73-use-case-replace-tags-on-deployment-templates-to-inject-correct-artifacts-for-gitops-using-bundle)
    4. [Replace Tags On Deployment Templates To Inject Correct Artifacts For GitOps Using Environment](#74-use-case-replace-tags-on-deployment-templates-to-inject-correct-artifacts-for-gitops-using-environment)
8. [Programmatic Approvals of Releases on Reliza Hub](#8-use-case-programmatic-approvals-of-releases-on-reliza-hub)
9. [Check if Specific Approval is Needed for a Release on Reliza Hub](#9-use-case-check-if-specific-approval-is-needed-for-a-release-on-reliza-hub)
10. [Persist Reliza Hub Credentials in a Config File](#10-use-case-persist-reliza-hub-credentials-in-a-config-file)
11. [Match list of images with digests to a bundle version on Reliza Hub](#11-use-case-match-list-of-images-with-digests-to-a-bundle-version-on-reliza-hub)
12. [Create New Project in Reliza Hub](#12-use-case-create-new-project-in-reliza-hub)
13. [Export Instance CycloneDX Spec](#13-use-case-export-instance-cyclonedx-spec)
14. [Add new artifacts to release in Reliza Hub](#14-use-case-add-new-artifacts-to-release-in-reliza-hub)
15. [Get changelog between releases in Reliza Hub](#15-use-case-get-changelog-between-releases-in-reliza-hub)
16. [Get specific properties and secrets defined for the instance in Reliza Hub](#16-use-case-get-specific-properties-and-secrets-defined-for-the-instance-in-reliza-hub)
17. [Export Bundle CycloneDX Spec](#17-use-case-export-bundle-cyclonedx-spec)
18. [Override and get merged helm chart values](#18-use-case-override-and-get-merged-helm-chart-values)
19. [Send Pull Request Data to Reliza Hub](#19-use-case-send-pull-request-data-to-reliza-hub)
20. [Attach a downloadable artifact to a Release on Reliza Hub](#20-use-case-attach-a-downloadable-artifact-to-a-release-on-reliza-hub)
## 1. Use Case: Get Version Assignment From Reliza Hub

This use case requests Version from Reliza Hub for our project. Note that project schema must be preset on Reliza Hub prior to using this API. API key must also be generated for the project from Reliza Hub.

Sample command for semver version schema:

```bash
docker run --rm relizaio/reliza-cli    \
    getversion    \
    -i project_or_organization_wide_rw_api_id    \
    -k project_or_organization_wide_rw_api_key    \
    -b master    \
    --pin 1.2.patch
```

Sample command with commit details for a git commit:

```bash
docker run --rm relizaio/reliza-cli    \
    getversion    \
    -i project_or_organization_wide_rw_api_id    \
    -k project_or_organization_wide_rw_api_key    \
    -b master    \
    --vcstype git \
    --commit $CI_COMMIT_SHA \
    --commitmessage $CI_COMMIT_MESSAGE \
    --vcsuri $CI_PROJECT_URL \
    --date $(git log -1 --date=iso-strict --pretty='%ad')
```

Sample command to obtain only version info and skip creating the release:

```bash
docker run --rm relizaio/reliza-cli    \
    getversion    \
    -i project_or_organization_wide_rw_api_id    \
    -k project_or_organization_wide_rw_api_key    \
    -b master    \
    --onlyversion
```

Flags stand for:

- **getversion** - command that denotes we are obtaining the next available release version for the branch. Note that if the call succeeds, the version assignment will be recorded and will not be given again by Reliza Hub, even if it is not consumed. It will create the release in pending status.
- **-i** - flag for project api id (required).
- **-k** - flag for project api key (required).
- **-b** - flag to denote branch (required). If the branch is not recorded yet, Reliza Hub will attempt to create it.
- **project** - flag to denote project uuid (optional). Required if organization-wide read-write key is used, ignored if project specific api key is used.
- **--pin** - flag to denote branch pin (optional for existing branches, required for new branches). If supplied for an existing branch and pin is different from current, it will override current pin.
- **--vcsuri** - flag to denote vcs uri (optional). This flag is needed if we want to set a commit for the release. However, soon it will be needed only if the vcs uri is not yet set for the project.
- **--vcstype** - flag to denote vcs type (optional). Supported values: git, svn, mercurial. As with vcsuri, this flag is needed if we want to set a commit for the release. However, soon it will be needed only if the vcs uri is not yet set for the project.
- **--commit** - flag to denote vcs commit id or hash (optional). This is needed to provide source code entry metadata into the release.
- **--commitmessage** - flag to denote vcs commit message (optional). Alongside *commit* flag this would be used to provide source code entry metadata into the release.
- **--commits** - flag to provide base64-encoded list of commits in the format *git log --date=iso-strict --pretty='%H|||%ad|||%s|||%an|||%ae' | base64 -w 0* (optional). If *commit* flag is not set, top commit will be used as commit bound to release.
- **--date** - flag to denote date time with timezone when commit was made, iso strict formatting with timezone is required, i.e. for git use git log --date=iso-strict (optional).
- **--vcstag** - flag to denote vcs tag (optional). This is needed to include vcs tag into commit, if present.
- **--metadata** - flag to set version metadata (optional). This may be semver metadata or custom version schema metadata.
- **--modifier** - flag to set version modifier (optional). This may be semver modifier or custom version schema metadata.
- **--manual** - flag to indicate a manual release (optional). Sets status as "draft", otherwise "pending" status is used.
- **--onlyversion** - boolean flag to skip creation of the release (optional). Default is false.

## 2. Use Case: Send Release Metadata to Reliza Hub

This use case is commonly used in the CI workflow to stream Release metadata to Reliza Hub. As in previous case, API key must be generated for the project on Reliza Hub prior to sending release details.

Sample command to send release details:

```bash
docker run --rm relizaio/reliza-cli    \
    addrelease    \
    -i project_or_organization_wide_rw_api_id    \
    -k project_or_organization_wide_rw_api_key    \
    -b master    \
    -v 20.02.3    \
    --vcsuri github.com/relizaio/reliza-cli    \
    --vcstype git    \
    --commit 7bfc5ce7b0da277d139f7993f90761223fa54442    \
    --vcstag 20.02.3    \
    --artid relizaio/reliza-cli    \
    --artbuildid 1    \
    --artcimeta Github Actions    \
    --arttype Docker    \
    --artdigests sha256:4e8b31b19ef16731a6f82410f9fb929da692aa97b71faeb1596c55fbf663dcdd    \
    --tagkey key1
    --tagval val1
```

Flags stand for:

- **addrelease** - command that denotes we are sending Release Metadata of a Project to Reliza Hub.
- **-i** - flag for project api id or organization-wide read-write api id (required).
- **-k** - flag for project api key or organization-wide read-write api key (required).
- **-b** - flag to denote branch (required). If branch is not recorded yet, Reliza Hub will attempt to create it.
- **-v** - version (required). Note that Reliza Hub will reject the call if a release with this exact version is already present for this project.
- **endpoint** - flag to denote test endpoint URI (optional). This would be useful for systems where every release gets test URI.
- **project** - flag to denote project uuid (optional). Required if organization-wide read-write key is used, ignored if project specific api key is used.
- **vcsuri** - flag to denote vcs uri (optional). Currently this flag is needed if we want to set a commit for the release. However, soon it will be needed only if the vcs uri is not yet set for the project.
- **vcstype** - flag to denote vcs type (optional). Supported values: git, svn, mercurial. As with vcsuri, this flag is needed if we want to set a commit for the release. However, soon it will be needed only if the vcs uri is not yet set for the project.
- **commit** - flag to denote vcs commit id or hash (optional). This is needed to provide source code entry metadata into the release.
- **commitmessage** - flag to denote vcs commit subject (optional). Alongside *commit* flag this would be used to provide source code entry metadata into the release.
- **commits** - flag to provide base64-encoded list of commits in the format *git log --date=iso-strict --pretty='%H|||%ad|||%s|||%an|||%ae' | base64 -w 0* (optional). If *commit* flag is not set, top commit will be used as commit bound to release.
- **date** - flag to denote date time with timezone when commit was made, iso strict formatting with timezone is required, i.e. for git use git log --date=iso-strict (optional).
- **vcstag** - flag to denote vcs tag (optional). This is needed to include vcs tag into commit, if present.
- **status** - flag to denote release status (optional). Supply "rejected" for failed releases, otherwise "complete" is used.
- **artid** - flag to denote artifact identifier (optional). This is required to add artifact metadata into release.
- **artbuildid** - flag to denote artifact build id (optional). This flag is optional and may be used to indicate build system id of the release (i.e., this could be circleci build number).
- **artbuilduri** - flag to denote artifact build uri (optional). This flag is optional and is used to denote the uri for where the build takes place.
- **artcimeta** - flag to denote artifact CI metadata (optional). This flag is optional and like artbuildid may be used to indicate build system metadata in free form.
- **arttype** - flag to denote artifact type (optional). This flag is used to denote artifact type. Types are based on [CycloneDX](https://cyclonedx.org/) spec. Supported values: Docker, File, Image, Font, Library, Application, Framework, OS, Device, Firmware.
- **datestart** - flag to denote artifact build start date and time, must conform to ISO strict date (in bash, use *date -Iseconds*, if used there must be one datestart flag entry per artifact, optional).
- **dateend** - flag to denote artifact build end date and time, must conform to ISO strict date (in bash, use *date -Iseconds*, if used there must be one datestart flag entry per artifact, optional).
- **artpublisher** - flag to denote artifact publisher (if used there must be one publisher flag entry per artifact, optional).
- **artversion** - flag to denote artifact version if different from release version (if used there must be one publisher flag entry per artifact, optional).
- **artpackage** - flag to denote artifact package type according to CycloneDX spec: MAVEN, NPM, NUGET, GEM, PYPI, DOCKER (if used there must be one publisher flag entry per artifact, optional).
- **artgroup** - flag to denote artifact group (if used there must be one group flag entry per artifact, optional).
- **artdigests** - flag to denote artifact digests (optional). This flag is used to indicate artifact digests. By convention, digests must be prefixed with type followed by colon and then actual digest hash, i.e. *sha256:4e8b31b19ef16731a6f82410f9fb929da692aa97b71faeb1596c55fbf663dcdd* - here type is *sha256* and digest is *4e8b31b19ef16731a6f82410f9fb929da692aa97b71faeb1596c55fbf663dcdd*. Multiple digests are supported and must be comma separated. I.e.:

```bash
--artdigests sha256:4e8b31b19ef16731a6f82410f9fb929da692aa97b71faeb1596c55fbf663dcdd,sha1:fe4165996a41501715ea0662b6a906b55e34a2a1
```

- **tagkey** - flag to denote keys of artifact tags (optional, but every tag key must have corresponding tag value). Multiple tag keys per artifact are supported and must be comma separated. I.e.:

```bash
--tagkey key1,key2
```

- **tagval** - flag to denote values of artifact tags (optional, but every tag value must have corresponding tag key). Multiple tag values per artifact are supported and must be comma separated. I.e.:

```bash
--tagval val1,val2
```

Note that multiple artifacts per release are supported. In which case artifact specific flags (artid, arbuildid, artbuilduri, artcimeta, arttype, artdigests, tagkey and tagval must be repeated for each artifact).

For sample of how to use workflow in CI, refer to the GitHub Actions build yaml of this project [here](https://github.com/relizaio/reliza-cli/blob/master/.github/workflows/dockerimage.yml).

## 3. Use Case: Check If Artifact Hash Already Present In Some Release

This is particularly useful for monorepos to see if there was a change in sub-project or not. We are using this technique in our sample [playground project](https://github.com/relizaio/reliza-hub-playground). We supply an artifact hash to Reliza Hub - and if it's present already, we get release details; if not - we get an empty json response {}. Search space is scoped to a single project which is defined by API Id and API Key.

Sample command:

```bash
docker run --rm relizaio/reliza-cli    \
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

```bash
docker run --rm relizaio/reliza-cli    \
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
- **--images** - flag which lists sha256 digests of images sent from the instances (optional, either images string or image file must be provided). Images must be white space separated. Note that sending full docker image URIs with digests is also accepted, i.e. it's ok to send images as relizaio/reliza-cli:latest@sha256:ebe68a0427bf88d748a4cad0a419392c75c867a216b70d4cd9ef68e8031fe7af
- **--imagefile** - flag which sets absolute path to the file with image string or image k8s json (optional, either images string or image file must be provided). Default value: */resources/images*. Use *kubectl get po -o json | jq "[.items[] | {namespace:.metadata.namespace, pod:.metadata.name, status:.status.containerStatuses[]}]"* to obtain k8s json.
- **--imagestyle** - flag which sets image format to k8s json if set to "k8s" (optional).
- **--namespace** - flag to denote namespace where we are sending images (optional, if not sent "default" namespace is used). Namespaces are useful to separate different products deployed on the same instance.
- **--sender** - flag to denote unique sender within a single namespace (optional). This is useful if say there are different nodes where each streams only part of application deployment data. In this case such nodes need to use same namespace but different senders so that their data does not stomp on each other.

## 5. Use Case: Request What Releases Must Be Deployed On This Instance From Reliza Hub

*DEPRECATED:* Note, this functionality is now deprecated and [13. Export Instance CycloneDX Spec](#13-use-case-export-instance-cyclonedx-spec) should be used instead where possible.

This use case is when your instance queries Reliza Hub to receive information about what release versions and specific artifacts it needs to deploy. This would usually be used by either simple deployment scripts or full-scale CD systems. A sample use is presented in our [playground project script](https://github.com/relizaio/reliza-hub-playground/blob/master/sample-instance-agent-scripts/request_instance_target.sh).

Sample command:

```bash
docker run --rm relizaio/reliza-cli    \
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

This use case is when Reliza Hub is queried either by CI or CD environment or by integration instance to check latest release version available per specific Project or Product. Only releases with *COMPLETE* status may be returned.

Sample command:

```bash
docker run --rm relizaio/reliza-cli    \
    getlatestrelease    \
    -i api_id    \
    -k api_key    \
    --project b4534a29-3309-4074-8a3a-34c92e1a168b    \
    --branch master    \
    --env TEST
```

Flags stand for:

- **getlatestrelease** - command that denotes we are requesting latest release data for Project or Product from Reliza Hub
- **-i** - flag for api id which can be either api id for this project or organization-wide read API (required).
- **-k** - flag for api key which can be either api key for this project or organization-wide read API (required).
- **--project** - flag to denote UUID of specific Project or Product, UUID must be obtained from [Reliza Hub](https://relizahub.com) (optional if project api key is used, otherwise required).
- **--product** - flag to denote UUID of Product which packages Project or Product for which we inquiry about its version via --project flag, UUID must be obtained from [Reliza Hub](https://relizahub.com) (optional).
- **--branch** - flag to denote required branch of chosen Project or Product (optional, if not supplied settings from Reliza Hub UI are used).
- **--env** - flag to denote environment to which release approvals should match. Environment can be one of: DEV, BUILD, TEST, SIT, UAT, PAT, STAGING, PRODUCTION. If not supplied, latest release will be returned regardless of approvals (optional).
- **--tagkey** - flag to denote tag key to use as a selector for artifact (optional, if provided tagval flag must also be supplied). Note that currently only single tag is supported.
- **--tagval** - flag to denote tag value to use as a selector for artifact (optional, if provided tagkey flag must also be supplied).
- **--instance** - flag to denote specific instance for which release should match (optional, if supplied namespace flag is also used and env flag gets overrided by instance's environment).
- **--namespace** - flag to denote specific namespace within instance, if instance is supplied (optional).
- **--status** - Status of the last known release to return, default is complete (optional, can be - [complete, pending or rejected]). If set to "pending", will return either pending or complete release. If set to "rejected", will return either pending or complete or rejected releases.

Here is a full example how we can use the getlatestrelease command leveraging jq to obtain the latest docker image with sha256 that we need to use for integration (don't forget to change api_id, api_key, project, branch and env to proper values as needed):

```bash
rlzclientout=$(docker run --rm relizaio/reliza-cli    \
    getlatestrelease    \
    -i api_id    \
    -k api_key    \
    --project b4534a29-3309-4074-8a3a-34c92e1a168b    \
    --branch master    \
    --env TEST);    \
    echo $(echo $rlzclientout | jq -r .artifactDetails[0].identifier)@$(echo $rlzclientout | jq -r .artifactDetails[0].digests[] | grep sha256)
```

## 7.1 Use Case: Parse Deployment Templates To Inject Correct Artifacts For GitOps

*DEPRECATED:* Note, this functionality is now deprecated and replacetags should be used instead where possible (section 7.2 and below).

This use case was designed specifically for GitOps. Imagine that you have GitOps deployment to different environments, i.e. TEST and PRODUCTION but they require different versions of artifacts. Reliza Hub would manage the versions but Reliza CLI can be leveraged to retrieve this information and create correct deployment files that can later be pushed to GitOps.

For a real-life use-case please refer to a working script in [deployment project for Classic Mafia Game Card Shuffle](https://github.com/taleodor/mafia-deployment/blob/master/pull_reliza_push_github_production.sh) while working templates can be found [here](https://github.com/taleodor/mafia-deployment/tree/master/k8s_templates).

Allowed template formatting types:

1. Basic project

    ```text
    image: <%PROJECT__9678805c-c8fd-4199-b682-1d5d2d73ad31%>
    ```

    where **9678805c-c8fd-4199-b682-1d5d2d73ad31** is a project UUID from [Reliza Hub](https://relizahub.com). In this format release branch would be resolved via settings in Reliza Hub UI in project settings -> **What branch to use for which environment?** setting.

2. Basic project with branch - template formatting may specify branch explicitly as following:

    ```text
    image: <%PROJECT__9678805c-c8fd-4199-b682-1d5d2d73ad31__master%>
    ```

    where **9678805c-c8fd-4199-b682-1d5d2d73ad31** is a project UUID from [Reliza Hub](https://relizahub.com) and **master** is our desired branch.

3. Project conditioned on Product

    ```text
    image: <%PROJECT__9678805c-c8fd-4199-b682-1d5d2d73ad31__PRODUCT__f407a320-8c3f-4658-be34-7635a69a8c05%>
    ```

    where **9678805c-c8fd-4199-b682-1d5d2d73ad31** is a project UUID from [Reliza Hub](https://relizahub.com), and **f407a320-8c3f-4658-be34-7635a69a8c05** is a product UUID from Reliza Hub which bundles this project we inquire about. In this format release feature set would be resolved via settings in Reliza Hub UI in product settings -> **What feature set to use for which environment?** setting.

4. Project conditioned on Product with explicit feature set

    ```text
    image: <%PROJECT__9678805c-c8fd-4199-b682-1d5d2d73ad31__PRODUCT__f407a320-8c3f-4658-be34-7635a69a8c05__Base Feature Set%>
    ```

    where **9678805c-c8fd-4199-b682-1d5d2d73ad31** is a project UUID from [Reliza Hub](https://relizahub.com), and **f407a320-8c3f-4658-be34-7635a69a8c05** is a product UUID from Reliza Hub which bundles this project we inquire about, and **Base Feature Set** is out desired feature set.

Sample command:

```bash
docker run --rm \
    -v /deployment/k8s_templates/:/indir
    -v /deployment/k8s_production/:/outdir
    relizaio/reliza-cli \
    parsetemplate \
    -i api_id \
    -k api_key \
    --env PRODUCTION
```

Note that selectors are generally identical to the **getlatestrelease** command.

Directory mapped to **/indir** (in this case **/deployment/k8s_templates/**) - is a directory containing parseable files with Reliza templates. Similarly, directory mapped to **/outdir** is a directory where output parsed files will be written. Both of those directories must exist.

Flags stand for:

- **parsetemplate** - command that denotes we are going to parse Reliza templates
- **-i** - flag for api id which can be either api id for this project or organization-wide read API (required).
- **-k** - flag for api key which can be either api key for this project or organization-wide read API (required).
- **--env** - flag to denote environment to which release approvals should match. Environment can be one of: DEV, BUILD, TEST, SIT, UAT, PAT, STAGING, PRODUCTION. If not supplied, latest release will be returned regardless of approvals (optional).
- **--tagkey** - flag to denote tag key to use as a selector for artifact (optional, if provided tagval flag must also be supplied). Note that currently only single tag is supported.
- **--tagval** - flag to denote tag value to use as a selector for artifact (optional, if provided tagkey flag must also be supplied).
- **--instance** - flag to denote specific instance for which releases should match (optional, if supplied namespace flag is also used and env flag gets overrided by instance's environment).
- **--namespace** - flag to denote specific namespace within instance, if instance is supplied (optional).
- **--indirectory** - input directory when using executable cli instead of docker, must use entire path (required if using executable)
- **--outdirectory** - output directory when using executable cli instead of docker, must use entire path (required if using executable)

## 7.2 Use Case: Replace Tags On Deployment Templates To Inject Correct Artifacts For GitOps Using Instance And Revision

This use case is designed for the case when we have to roll back our deployments to a specific version of artifacts. Reliza CLI can be leveraged to update deployments with the correct version of artifacts that can be pushed to GitOps.

Sample Command:

```text
docker run --rm \
    -v /local/path/to/values.yaml:/values.yaml \
    -v /local/path/to/output_dir:/output_dir \
    relizaio/reliza-cli \
    replacetags \
    --instanceuri <instance uri> \
    --revision <revision_number> \
    --infile /values.yaml \
    --outfile /output_dir/output_values.yaml
```

Flags stand for:

- **-i** - flag for api id which can be either api id for this project or organization-wide read API (required).
- **-k** - flag for api key which can be either api key for this project or organization-wide read API (required).
- **--instanceuri** - URI of the instance (optional, either instanceuri or instance or tagsource flag must be used).
- **--instance** - UUID of the instance (optional, either instanceuri or instance or tagsource flag must be used).
- **--revision** - Revision number for the instance to use as a source for tags (optional, if not specified tags will be resolved by environment to which the instance belongs).
- **--namespace** - Specific namespace of the instance to use to retrieve tag sources (optional).
- **--infile** - Input file to parse, such as helm values file or docker compose file.
- **--outfile** - Output file with parsed values (optional, if not supplied - outputs to stdout).
- **--indirectory** - Path to directory of input files to parse (either infile or indirectory is required)
- **--outdirectory** - Path to directory of output files (required if indirectory is used)
- **--tagsource** - Source file with tags (optional, either instanceuri or instance or tagsource flag must be used).
- **--defsource** - Source file for definitions. For helm, should be output of helm template command. (Optional, if not specified - *infile* will be parsed for definitions).
- **--type** - Type of source tags file: cyclonedx (default) or text.
- **--provenance** - Set --provenance=[true|false] flag to enable/disable adding provenance (metadata) to beginning of outfile. (optional, default true)
- **--parsemode** - Flag to set the parse mode. *Extended*: normal operation. *Simple*: Only replace 'image' tags. *Strict*: Exit process if an artifact is not found upstream.(optional) (default extended)
- **--resolveprops** - Set --resolveprops=[true|false] flag to true to enable resolution of properties and secrets from instance - see below for details. (optional, default false)
- **--fordiff** - Set --fordiff=[true|false] flag to true to resolve templated secrets to their timestamps instead of sealed secret value. This can be used to establish where an update happened, since the sealed value otherwise may be changed every time. If true, this will also disable provenance regardless of its flag. (optional, default false)

To resolve secrets and properties from instances, the resolveprops flag must be set to true. Other than that, in the templated file the properties should be defined as following:

`$RELIZA{PROPERTY.property_key}` - where `property_key` part must be set on the corresponding instance on the Reliza Hub, or `$RELIZA{PROPERTY.property_key:default_value}` to also set a default value in case the property key is not found on Reliza Hub.

While secrets should be defined as:

`$RELIZA{SECRET.secret_key}` - where `secret_key` part must be set on the Reliza Hub. More so, the secret must be allowed for usage by particular instance. Finally, instance must have a property set for the sealed certificates, since we are only sending sealed certificates and not in plain text.

`$RELIZA{PLAINSECRET.secret_key}` is same as SECRET but resolves to plain value. This only works in the reliza-cd context.

## 7.3 Use Case: Replace Tags On Deployment Templates To Inject Correct Artifacts For GitOps Using Bundle

This use case is designed for the case when we have to deploy a specific version of a bundle or approved bundle by environment. Reliza CLI can be leveraged to update deployments with the correct version of artifacts that can be pushed to GitOps.

Sample Command:

```text
docker run --rm \
    -v /local/path/to/values.yaml:/values.yaml \
    -v /local/path/to/output_dir:/output_dir \
    relizaio/reliza-cli \
    replacetags \
    --bundle <bundle name> \
    --version <bundle version> \
    --infile /values.yaml \
    --outfile /output_dir/output_values.yaml
```

Flags stand for:

- **-i** - flag for api id which can be a organization-wide read API (required).
- **-k** - flag for api key which can be a organization-wide read API (required).
- **--bundle** - Name of the bundle (optional, either bundle name & version or tagsource flag must be used).
- **--version** - Version number for the bundle to use as a source for tags (optional, either version or environment must be used with the bundle flag).
- **--environment** - Environment for which latest approved bundle should be used as a source for tags (optional, either version or environment must be used bundle flag).
- **--infile** - Input file to parse, such as helm values file or docker compose file.
- **--outfile** - Output file with parsed values (optional, if not supplied - outputs to stdout).
- **--indirectory** - Path to directory of input files to parse (either infile or indirectory is required)
- **--outdirectory** - Path to directory of output files (required if indirectory is used)
- **--tagsource** - Source file with tags (optional, either bundle name & version or tagsource flag must be used).
- **--defsource** - Source file for definitions. For helm, should be output of helm template command. (Optional, if not specified - *infile* will be parsed for definitions).
- **--type** - Type of source tags file: cyclonedx (default) or text.
- **--provenance** - Set --provenance=[true|false] flag to enable/disable adding provenance (metadata) to beginning of outfile. (optional) (default true)
- **--parsemode** - Flag to set the parse mode. *Extended*: normal operation. *Simple*: Only replace 'image' tags. *Strict*: Exit process if an artifact is not found upstream.(optional) (default extended)

## 7.4 Use Case: Replace Tags On Deployment Templates To Inject Correct Artifacts For GitOps Using Environment

This use case is designed for the case when we have to deploy to a specific environment. Reliza CLI can be leveraged to update deployments with the correct version of artifacts that can be pushed to GitOps.

Sample Command:

```text
docker run --rm \
    -v /local/path/to/values.yaml:/values.yaml \
    -v /local/path/to/output_dir:/output_dir \
    relizaio/reliza-cli \
    replacetags \
    --env <environment name> \
    --infile /values.yaml \
    --outfile /output_dir/output_values.yaml
```

Flags stand for:

- **-i** - flag for api id which can be a organization-wide read API (required).
- **-k** - flag for api key which can be a organization-wide read API (required).
- **--env** - flag to denote the environment to which we wish to deploy. Environment can be one of: DEV, BUILD, TEST, SIT, UAT, PAT, STAGING, PRODUCTION.
- **--infile** - Input file to parse, such as helm values file or docker compose file.
- **--outfile** - Output file with parsed values (optional, if not supplied - outputs to stdout).
- **--indirectory** - Path to directory of input files to parse (either infile or indirectory is required)
- **--outdirectory** - Path to directory of output files (required if indirectory is used)
- **--defsource** - Source file for definitions. For helm, should be output of helm template command. (Optional, if not specified - *infile* will be parsed for definitions).
- **--provenance** - Set --provenance=[true|false] flag to enable/disable adding provenance (metadata) to beginning of outfile. (optional) (default true)
- **--parsemode** - Flag to set the parse mode. *Extended*: normal operation. *Simple*: Only replace 'image' tags. *Strict*: Exit process if an artifact is not found upstream.(optional) (default extended)

## 8. Use Case: Programmatic Approvals of Releases on Reliza Hub

This use case is for the case when we have configured an API key in Org settings which is allowed to perform programmatic approvals in releases.

Sample command:

```bash
docker run --rm relizaio/reliza-cli    \
    approverelease    \
    -i api_id    \
    -k api_key    \
    --release release_uuid    \
    --approval approval_type
```

Flags stand for:

- **approverelease** - command that denotes that we are approving release programmatically for the specific type
- **-i** - flag for api id (required).
- **-k** - flag for api key (required).
- **--releaseid** - flag to specify release uuid, which can be obtained from the release view or programmatically (either this flag or project id and release version or project id and instance are required).
- **--project** - flag to specify project uuid, which can be obtained from the project settings on Reliza Hub UI (either this flag and release version or releaseid must be provided).
- **--instance** - flag to specify instance uuid or URI for which release must be approved (either this flag and project or project and release version or releaseid must be provided).
- **--namespace** - flag to specify namespace of the instance for which release must be approved (optional, only taken in consideration if instance is provided).
- **--releaseversion** - flag to specify release string version with the project flag above (either this flag and project or releaseid must be provided).
- **--approval** - approval type as per approval matrix on the Organization Settings page in Reliza Hub (required).
- **--disapprove** - flag to indicate disapproval event instead of approval (optional).

## 9. Use Case: Check if Specific Approval is Needed for a Release on Reliza Hub

This use case is auxiliary to the previous use case with programmatic approvals. It checks Reliza Hub if a specific approval type is still pending for a release. For example, some approval might have already been given previously, or the release may have already been rejected - in both of these cases, an approval is not needed any more. Such check may be useful for example, to decide whether to perform a set of automated tests for a release.

Sample command:

```bash
docker run --rm relizaio/reliza-cli    \
    isapprovalneeded    \
    -i api_id    \
    -k api_key    \
    --release release_uuid    \
    --approval approval_type
```

Flags stand for:

- **isapprovalneeded** - command that denotes that we are programmatically checking if the approval is needed for a particular release
- **-i** - flag for api id (required).
- **-k** - flag for api key (required).
- **--releaseid** - flag to specify release uuid, which can be obtained from the release view or programmatically (either this flag or project id and release version or project id and instance are required).
- **--project** - flag to specify project uuid, which can be obtained from the project settings on Reliza Hub UI (either this flag and release version or releaseid must be provided).
- **--instance** - flag to specify instance uuid or URI for which release must be approved (either this flag and project or project and release version or releaseid must be provided).
- **--namespace** - flag to specify namespace of the instance for which release must be approved (optional, only taken in consideration if instance is provided).
- **--releaseversion** - flag to specify release string version with the project flag above (either this flag and project or releaseid must be provided).
- **--approval** - approval type as per approval matrix on the Organization Settings page in Reliza Hub (required).

## 10. Use Case: Persist Reliza Hub Credentials in a Config File

This use case is for the case when we want to persist Reliza Hub API Credentials and URL in a config file.

The `login` command saves `API ID`, `API KEY` and `URI` as specified by flags in a config file `.reliza.env` in the home directory for the executing user.

Sample Command:

```bash
docker run --rm \
    -v ~:/home/apprunner \
    relizaio/reliza-cli \
    login \
    -i api_id \
    -k api_key \
    -u reliza_hub_uri
```

Flags stand for:

- **-i** - flag for api id.
- **-k** - flag for api key.
- **-u** - flag for reliza hub uri.


## 11. Use Case: Match list of images with digests to a bundle version on Reliza Hub

This use case is to match a list of images with digests, in example on local Docker enviornment to a bundle version on Reliza Hub. Works with User or Organization-Wide API-keys.

Sample dockerized command:

```bash
docker run --rm relizaio/reliza-cli    \
    matchbundle    \
    -i api_id    \
    -k api_key    \
    --images "sha256:c10779b369c6f2638e4c7483a3ab06f13b3f57497154b092c87e1b15088027a5 sha256:e6c2bcd817beeb94f05eaca2ca2fce5c9a24dc29bde89fbf839b652824304703"
```

Sample flow to use for matching local docker images to a bundle release:

```bash
images=$(docker ps --no-trunc | awk 'NR>2 {print $2}' | tr "\n" " ")
reliza-cli matchbundle --images "$images"
```

Flags stand for:

- **matchbundle** - command that denotes we are trying to match list of images to a bundle release.
- **-i** - flag for api id (either User, or Organization, or Organization Read-Write, can be obtained via Reliza Hub, required).
- **-k** - flag for api key (either User, or Organization, or Organization Read-Write, can be obtained via Reliza Hub, required).
- **--images** - flag which lists images with sha256 digests or only digests of images sent from the instances (optional, either images string or image file must be provided). Images must be white space separated. Note that sending full docker image URIs with digests is also accepted, i.e. it's ok to send images as relizaio/reliza-cli:latest@sha256:ebe68a0427bf88d748a4cad0a419392c75c867a216b70d4cd9ef68e8031fe7af
- **--imagefile** - flag which sets absolute path to the file with image string or image k8s json (optional, either images string or image file must be provided). Default value: */resources/images*.
- **--namespace** - flag to denote namespace where we are sending images (optional, unused, present for compatibility with instance data command, which uses simialr underlying logic).

## 12. Use Case: Create New Project in Reliza Hub

This use case creates a new project for our organization. API key must be generated prior to using.

Sample command to create project:

```bash
docker run --rm relizaio/reliza-cli    \
    createproject    \
    -i org_api_id    \
    -k org_api_key    \
    --name projectname
    --type project
    --versionschema semver
    --featurebranchversioning Branch.Micro
    --vcsuri github.com/relizaio/reliza-cli
    --includeapi
```

Flags stand for:

- **createproject** - command that denotes we are creating a new project for our organization. Note that a vcs repository must either already exist or be created during this call.
- **-i** - flag for org api id (required).
- **-k** - flag for org api key (required).
- **name** - flag to denote project name (required).
- **type** - flag to denote project type (required). Supported values: project, bundle.
- **defaultbranch** - flag to denote default branch name (optional, if not set "main" will be used). Available names are either main or master.
- **versionschema** - flag to denote version schema (optional, if not set "semver" will be used). [Available version schemas](https://github.com/relizaio/versioning).
- **featurebranchversioning** - flag to denote feature branch version schema (optional, if not set "Branch.Micro will be used).
- **vcsuuid** - flag to denote uuid of vcs repository for the project (for existing repositories, either this flag or vcsuri are required).
- **vcsuri** - flag to denote uri of vcs repository for the project, if existing repository with uri does not exist and vcsname and vcstype are not set, Reliza Hub will attempt to autoparse github, gitlab, and bitbucket uri's.
- **vcsname** - flag to denote name of vcs repository to create for project (required if Reliza Hub cannot parse uri).
- **vcstype** - flag to denote type of vcs to create for project. Supported values: git, svn, mercurial (required if Reliza Hub cannot parse uri).
- **includeapi** - boolean flag to return project api key and id of newly created project (optional). Default is false.

## 13. Use Case: Export Instance CycloneDX Spec

This use case exports the present, past or expected state of the instance in [CycloneDX](https://cyclonedx.org) format. API key must be generated prior to using.

The **--revision** flag is what defines the type of the state (present, past, expected). It behaves is following:
- Default value is *-1*, which means *expected* state - this will output all project releases that *are approved* for the specific instance.
- The value set to *-2* means *present* state - this would output project releases currently deployed on the specific instance.
- The value set to an actual revision obtained from Reliza Hub would output project releases deployed on that specific revision.

Sample command:

```bash
docker run --rm relizaio/reliza-cli    \
    exportinst    \
    -i api_id    \
    -k api_key    \
    --instance instanceUuid     \
    --revision -2
```

Flags stand for:

- **exportinst** - command that denotes we are exporting the CycloneDX spec from our instance.
- **-i** - flag for api id (required).
- **-k** - flag for api key (required).
- **instance** - flag to denote instance UUID (either instance api, instance, or instanceuri field or Instance API Key must be supplied).
- **instanceuri** - flag to denote instance URI (either instance api, instance, or instanceuri or Instance API Key field must be supplied).
- **revision** - Revision number for the instance (optional, default value is -1).
- **--namespace** - Specific namespace of the instance - if provided, only deployed releases on this particular namespace will be exported (optional).

## 14. Use Case: Add new artifacts to release in Reliza Hub

This use case adds 1 or more artifacts to an existing release. API key must be generated prior to using.

Sample command to add artifact:

```bash
docker run --rm relizaio/reliza-cli    \
    addartifact    \
    -i project_or_organization_wide_rw_api_id    \
    -k project_or_organization_wide_rw_api_key    \
    -v 20.02.3    \
    --artid relizaio/reliza-cli    \
    --artbuildid 1    \
    --artcimeta Github Actions    \
    --arttype Docker    \
    --artdigests sha256:4e8b31b19ef16731a6f82410f9fb929da692aa97b71faeb1596c55fbf663dcdd    \
    --tagkey key1
    --tagval val1
```

Flags stand for:

- **addartifact** - command that denotes we are adding artifact(s) to a release.
- **-i** - flag for project api id or organization-wide read-write api id (required).
- **-k** - flag for project api key or organization-wide read-write api key (required).
- **releaseid** - flag to specify release uuid, which can be obtained from the release view or programmatically (either this flag or project and version are required).
- **project** - flag to denote project uuid (optional). Required if organization-wide read-write key is used and releaseid isn't, ignored if project specific api key is used.
- **version** - version (either this flag and project or releaseid are required)
- **artid** - flag to denote artifact identifier (optional). This is required to add artifact metadata into release.
- **artbuildid** - flag to denote artifact build id (optional). This flag is optional and may be used to indicate build system id of the release (i.e., this could be circleci build number).
- **artbuilduri** - flag to denote artifact build uri (optional). This flag is optional and is used to denote the uri for where the build takes place.
- **artcimeta** - flag to denote artifact CI metadata (optional). This flag is optional and like artbuildid may be used to indicate build system metadata in free form.
- **arttype** - flag to denote artifact type (optional). This flag is used to denote artifact type. Types are based on [CycloneDX](https://cyclonedx.org/) spec. Supported values: Docker, File, Image, Font, Library, Application, Framework, OS, Device, Firmware.
- **datestart** - flag to denote artifact build start date and time, must conform to ISO strict date (in bash, use *date -Iseconds*, if used there must be one datestart flag entry per artifact, optional).
- **dateend** - flag to denote artifact build end date and time, must conform to ISO strict date (in bash, use *date -Iseconds*, if used there must be one datestart flag entry per artifact, optional).
- **artpublisher** - flag to denote artifact publisher (if used there must be one publisher flag entry per artifact, optional).
- **artversion** - flag to denote artifact version if different from release version (if used there must be one publisher flag entry per artifact, optional).
- **artpackage** - flag to denote artifact package type according to CycloneDX spec: MAVEN, NPM, NUGET, GEM, PYPI, DOCKER (if used there must be one publisher flag entry per artifact, optional).
- **artgroup** - flag to denote artifact group (if used there must be one group flag entry per artifact, optional).
- **artdigests** - flag to denote artifact digests (optional). This flag is used to indicate artifact digests. By convention, digests must be prefixed with type followed by colon and then actual digest hash, i.e. *sha256:4e8b31b19ef16731a6f82410f9fb929da692aa97b71faeb1596c55fbf663dcdd* - here type is *sha256* and digest is *4e8b31b19ef16731a6f82410f9fb929da692aa97b71faeb1596c55fbf663dcdd*. Multiple digests are supported and must be comma separated. I.e.:

```bash
--artdigests sha256:4e8b31b19ef16731a6f82410f9fb929da692aa97b71faeb1596c55fbf663dcdd,sha1:fe4165996a41501715ea0662b6a906b55e34a2a1
```

- **tagkey** - flag to denote keys of artifact tags (optional, but every tag key must have corresponding tag value). Multiple tag keys per artifact are supported and must be comma separated. I.e.:

```bash
--tagkey key1,key2
```

- **tagval** - flag to denote values of artifact tags (optional, but every tag value must have corresponding tag key). Multiple tag values per artifact are supported and must be comma separated. I.e.:

```bash
--tagval val1,val2
```

Notes:
1. Multiple artifacts per release are supported. In which case artifact specific flags (artid, arbuildid, artbuilduri, artcimeta, arttype, artdigests, tagkey and tagval must be repeated for each artifact).
2. Artifacts may be added to Complete or Rejected releases (this can be used for adding for example test reports), however a special tag would be placed on those artifacts by Reliza Hub.

## 15. Use Case: Get changelog between releases in Reliza Hub

This use case constructs a changelog using 2 different reference points from your project. API key must be generated prior to using.

Sample command:

```bash
docker run --rm relizaio/reliza-cli    \
    getchangelog    \
    -i project_or_organization_wide_api_id    \
    -k project_or_organization_wide_api_key    \
    --version 1.3.36     \
    --version2 1.3.33
    --aggregated
```

Flags stand for:

- **getchangelog** - command that denotes we are constructing a changelog.
- **-i** - flag for api id which can be either api id for this project or organization-wide read API (required).
- **-k** - flag for api key which can be either api key for this project or organization-wide read API (required).
- **project** - flag to denote project UUID (required only if using org-wide key and attaining changelog using versions).
- **version** - Release version (either this and version2 or commit and commit2 must be supplied).
- **version2** - Second release version to construct changelog from.
- **commit** - Commit id (either this and commit2 or version and version2 must be supplied).
- **commit2** - Second commit id to construct changelog from.
- **revision** - Boolean flag to create aggregated changelog (optional). Default is false.

## 16. Use Case: Get specific properties and secrets defined for the instance in Reliza Hub

This use case retrieves properties and secrets set for the instance on Reliza Hub. Note that secrets are only retrieved in sealed form and require (Bitnami Sealed Secret)[https://github.com/bitnami-labs/sealed-secrets] certificate property to be set on the instance - the key for that property is `SEALED_SECRETS_CERT`.

Note that secrets must be allowed to be read by the particular instance.

Also note that secrets are retrieved as sealed with the namespace-scope.

Sample command:

```bash
docker run --rm relizaio/reliza-cli    \
    instprops    \
    -i instance_or_organization_wide_api_id    \
    -k instance_or_organization_wide_api_key    \
    --instanceuri test.com     \
    --property FQDN
    --property my_property
    --secret test_secret
    --secret test_secret2

Flags stand for:

- **instprops** - command that denotes we are retrieving properties and secrets for the instance.
- **-i** - flag for api id which can be either api id for this instance or organization-wide read API (required).
- **-k** - flag for api key which can be either api key for this instance or organization-wide read API (required).
- **--instanceuri** - URI of the instance (optional, either instanceuri or instance flag must be used).
- **--instance** - UUID of the instance (optional, either instanceuri or instance flag must be used).
- **--revision** - Revision number for the instance to use as a source for properties (optional, defaults to latest).
- **--namespace** - Specific namespace of the instance to use to retrieve sealed secrets - as secrets are returned sealed with namespace scope (optional, default to "default").
- **--property** - Specifies name of the property to retrieve. For multiple properties, use multiple --property flags.
- **--secret** - Specifies name of the secret to retrieve. For multiple secrets, use multiple --secret flags.
```

## 17. Use Case: Export Bundle CycloneDX Spec

This use case exports a specific version of a bundle in [CycloneDX](https://cyclonedx.org) format. API key must be generated prior to using.

Sample command:

```bash
docker run --rm relizaio/reliza-cli    \
    exportbundle    \
    -i api_id    \
    -k api_key    \
    --bundle <bundle name>    \
    --version <bundle version>
```

Flags stand for:

- **exportbundle** - command that denotes we are exporting the CycloneDX spec from our bundle.
- **-i** - flag for api id (required).
- **-k** - flag for api key (required).
- **bundle** - flag to denote bundle name (required).
- **version** - flag to denote bundle version (either version or environment must be set).
- **env** - flag to denote environment, for which to export latest approved bundle (either version or environment must be set).


## 18. Use Case: Override and get merged helm chart values

This use case lets you do a helm style override of the default helm chart values and outputs merged helm values.

Sample command:
```bash
docker run --rm relizaio/reliza-cli    \
    helmvalues <Absolute or Relative Path to the Chart>   \
    -f <values-override-1.yaml>    \
    -f <values-override-2.yaml>    \
    -o <output-values.yaml>
```

Flags stand for:

- **--outfile | -o** - Output file with merge values (optional, if not supplied - outputs to stdout).
- **--values | -f** - Specify override values YAML file. Indicate file name only here, path would be resolved according to path to the chart in the command. Can specify multiple value file - in that case and if different values files define same properties, properties in the files that appear later in the command will take precedence - just like helm works.

## 19. Use Case: Send Pull Request Data to Reliza Hub

This use case is used in the CI workflow to stream Pull Request metadata to Reliza Hub.

Sample command to send Pull Request details:

Sample command:
```bash
docker run --rm relizaio/reliza-cli    \
    prdata \
    -i project_or_organization_wide_api_id    \
    -k project_or_organization_wide_api_key    \
    -b <base branch name> \
    -s <pull request state - OPEN | CLOSED | MERGED> \
    -t <target branch name> \
    --endpoint <pull request endpoint> \
    --title <title> \
    --createdDate <ISO 8601 date > \
    --number <pull request number> \
    --commits <comma separated list of commit shas>
```

Flags stand for:

- **--branch | -b** - Name of the base branch for the pull request.
- **--state** - State of the pull request
- **--targetBranch | t** - Name of the target branch for the pull request.
- **--endpoint** - HTML endpoint of the pull request.
- **--title** - Title of the pull request.
- **--number** - Number of the pull request.
- **--commits** - Comma seprated commit shas on this pull request.
- **--commits** - SHA of current commit on the Pull Request (will be merged with existing list)
- **--createdDate** - Datetime when the pull request was created.
- **--closedDate** - Datetime when the pull request was closed.
- **--mergedDate** - Datetime when the pull request was merged.
- **--endpoint** - Title of the pull request.
- **--project** - Project UUID if org-wide key is used.

## 20. Use Case: Attach a downloadable artifact to a Release on Reliza Hub

This use case is to attach a downloadable artifact to a Release on Reliza Hub. For example, to add a report obtained by automated tests for a release.

Sample command:

```bash
docker run --rm relizaio/reliza-cli    \
    addDownloadableArtifact \
    -i api_id \ 
    -k api_key    \
    --releaseid release_uuid    \
    --artifactType TEST_REPORT \
    --file <path_to_the_report>
```

Flags stand for:
- **--file | -f** - flag to specify path to the artifact file.
- **--releaseid** - flag to specify release uuid, which can be obtained from the release view or programmatically (either this flag or project id and release version or project id and instance are required).
- **--project** - flag to specify project uuid, which can be obtained from the project settings on Reliza Hub UI (either this flag and release version or releaseid must be provided).
- **--instance** - flag to specify instance uuid or URI for which release must be approved (either this flag and project or project and release version or releaseid must be provided).
- **--namespace** - flag to specify namespace of the instance for which release must be approved (optional, only taken in consideration if instance is provided).
- **--releaseversion** - flag to specify release string version with the project flag above (either this flag and project or releaseid must be provided).
- **--artifactType** - flag to specify type of the artifact - can be (TEST_REPORT, SECURITY_SCAN, DOCUMENTATION, GENERIC) or some user defined value .

# Development of Reliza-CLI

## Adding dependencies to Reliza-CLI

Dependencies are handled using go modules and imports file is automatically generated. If importing a github repository use this command first:

```bash
go get github.com/xxxxxx/xxxxxx
```

You then should be able to add what you need as an import to your files. Once they've been imported call this command to generate the imports file:

```bash
go generate ./internal/imports
```
