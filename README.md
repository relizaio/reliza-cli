![Docker Image CI](https://github.com/relizaio/reliza-cli/workflows/Docker%20Image%20CI/badge.svg?branch=master)

# Reliza CLI

This tool allows for command-line interactions with [Reliza Hub at relizahub.com](https://relizahub.com) (currently in public preview mode). Particularly, Reliza CLI can stream metadata about instances, releases, artifacts, resolve bundles based on Reliza Hub data. Available as either a Docker image or binary.

Video tutorial about key functionality of Reliza Hub is available on [YouTube](https://www.youtube.com/watch?v=yDlf5fMBGuI).

Community forum and support is available at [r/Reliza](https://reddit.com/r/Reliza).

Docker image is available at [relizaio/reliza-cli](https://hub.docker.com/r/relizaio/reliza-cli)

## Download Reliza CLI

Below are the available downloads for the latest version of the Reliza CLI (2021.03.3). Please download the proper package for your operating system and architecture.

The CLI is distributed as a single binary. Install by unzipping it and moving it to a directory included in your system's PATH.

[SHA256 checksums](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/sha256sums.txt)

macOS: [64-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/reliza-cli-2021.03.3-darwin-amd64.zip)

FreeBSD: [32-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/reliza-cli-2021.03.3-freebsd-386.zip) | [64-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/reliza-cli-2021.03.3-freebsd-amd64.zip) | [Arm](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/reliza-cli-2021.03.3-freebsd-arm.zip)

Linux: [32-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/reliza-cli-2021.03.3-linux-386.zip) | [64-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/reliza-cli-2021.03.3-linux-amd64.zip) | [Arm](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/reliza-cli-2021.03.3-linux-arm.zip) | [Arm64](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/reliza-cli-2021.03.3-linux-arm64.zip)

OpenBSD: [32-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/reliza-cli-2021.03.3-openbsd-386.zip) | [64-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/reliza-cli-2021.03.3-openbsd-amd64.zip)

Solaris: [64-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/reliza-cli-2021.03.3-solaris-amd64.zip)

Windows: [32-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/reliza-cli-2021.03.3-windows-386.zip) | [64-bit](https://d7ge14utcyki8.cloudfront.net/reliza-cli-download/2021.03.3/reliza-cli-2021.03.3-windows-amd64.zip)

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
- **date** - flag to denote date time with timezone when commit was made, iso strict formatting with timezone is required, i.e. for git use git log --date=iso-strict (optional).
- **vcstag** - flag to denote vcs tag (optional). This is needed to include vcs tag into commit, if present.
- **status** - flag to denote release status (optional). Supply "rejected" for failed releases, otherwise "completed" is used.
- **artid** - flag to denote artifact identifier (optional). This is required to add artifact metadata into release.
- **artbuildid** - flag to denote artifact build id (optional). This flag is optional and may be used to indicate build system id of the release (i.e., this could be circleci build number).
- **artbuilduri** - flag to denote artifact build uri (optional). This flag is optional and is used to denote the uri for where the build takes place.
- **artcimeta** - flag to denote artifact CI metadata (optional). This flag is optional and like artbuildid may be used to indicate build system metadata in free form.
- **arttype** - flag to denote artifact type (optional). This flag is used to denote artifact type. Types are based on [CycloneDX](https://cyclonedx.org/) spec. Supported values: Docker, File, Image, Font, Library, Application, Framework, OS, Device, Firmware.
- **datestart** - flag to denote artifact build start date and time, must conform to ISO strict date (in bash, use *date -Iseconds*, if used there must be one datestart flag entry per artifact, optional).
- **dateend** - flag to denote artifact build end date and time, must conform to ISO strict date (in bash, use *date -Iseconds*, if used there must be one datestart flag entry per artifact, optional).
- **publisher** - flag to denote artifact publisher (if used there must be one publisher flag entry per artifact, optional).
- **version** - flag to denote artifact version if different from release version (if used there must be one publisher flag entry per artifact, optional).
- **package** - flag to denote artifact package type according to CycloneDX spec: MAVEN, NPM, NUGET, GEM, PYPI, DOCKER (if used there must be one publisher flag entry per artifact, optional).
- **group** - flag to denote artifact group (if used there must be one group flag entry per artifact, optional).
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
- **--imagefile** - flag which sets absolute path to the file with image string or image k8s json (optional, either images string or image file must be provided). Default value: */resources/images*. Use *kubectl get po -o json | jq "[{namespace:.items[].metadata.namespace, pod:.items[].metadata.name, status:.items[].status.containerStatuses[]}]"* to obtain k8s json.
- **--imagestyle** - flag which sets image format to k8s json if set to "k8s" (optional).
- **--namespace** - flag to denote namespace where we are sending images (optional, if not sent "default" namespace is used). Namespaces are useful to separate different products deployed on the same instance.
- **--sender** - flag to denote unique sender within a single namespace (optional). This is useful if say there are different nodes where each streams only part of application deployment data. In this case such nodes need to use same namespace but different senders so that their data does not stomp on each other.

## 5. Use Case: Request What Releases Must Be Deployed On This Instance From Reliza Hub

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

This use case is when Reliza Hub is queried either by CI or CD environment or by integration instance to check latest release version available per specific Project or Product.

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
- **--project** - flag to denote UUID of specific Project or Product, UUID must be obtained from [Reliza Hub](https://relizahub.com) (required).
- **--product** - flag to denote UUID of Product which packages Project or Product for which we inquiry about its version via --project flag, UUID must be obtained from [Reliza Hub](https://relizahub.com) (optional).
- **--branch** - flag to denote required branch of chosen Project or Product (optional, if not supplied settings from Reliza Hub UI are used).
- **--env** - flag to denote environment to which release approvals should match. Environment can be one of: DEV, BUILD, TEST, SIT, UAT, PAT, STAGING, PRODUCTION. If not supplied, latest release will be returned regardless of approvals (optional).
- **--tagkey** - flag to denote tag key to use as a selector for artifact (optional, if provided tagval flag must also be supplied). Note that currently only single tag is supported.
- **--tagval** - flag to denote tag value to use as a selector for artifact (optional, if provided tagkey flag must also be supplied).
- **--instance** - flag to denote specific instance for which release should match (optional, if supplied namespace flag is also used and env flag gets overrided by instance's environment).
- **--namespace** - flag to denote specific namespace within instance, if instance is supplied (optional).

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
- **--instance** - UUID of the instanceo (optional, either instanceuri or instance or tagsource flag must be used).
- **--revision** - Revision number for the instance to use as a source for tags (optional, if not specified latest revision will be assumed).
- **--infile** - Input file to parse, such as helm values file or docker compose file.
- **--outfile** - Output file with parsed values.
- **--tagsource** - Source file with tags (optional, either instanceuri or instance or tagsource flag must be used).
- **--defsource** - Source file for definitions. For helm, should be output of helm template command. (Optional, if not specified - *infile* will be parsed for definitions).
- **type** - Type of source tags file: cyclonedx (default) or text.

## 7.3 Use Case: Replace Tags On Deployment Templates To Inject Correct Artifacts For GitOps Using Bundle And Version

This use case is designed for the case when we have to deploy a specific version of a bundle. Reliza CLI can be leveraged to update deployments with the correct version of artifacts that can be pushed to GitOps.

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
- **--version** - Version number for the bundle to use as a source for tags (optional, to be used with bundle flag).
- **--infile** - Input file to parse, such as helm values file or docker compose file.
- **--outfile** - Output file with parsed values.
- **--tagsource** - Source file with tags (optional, either bundle name & version or tagsource flag must be used).
- **--defsource** - Source file for definitions. For helm, should be output of helm template command. (Optional, if not specified - *infile* will be parsed for definitions).
- **--type** - Type of source tags file: cyclonedx (default) or text.

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

## 10. Use Case: Persist Reliza Hub Credentials In A Config File

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
