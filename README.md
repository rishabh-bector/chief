# Chief

Chief is an elegant continuous integration & deployment server, for minimalists.

## Installation

Run the following commands on a linux server:

1. `insert wget/curl command here`

2. `chief setup`

3. `chief start`

The `chief setup` command should only have to be run once, on installation. You'll be prompted to enter a username/password, to create your Chief user. This first user will have master clearance, so it can add/remove other users and manage the Chief server process. 

The `chief start` command starts the chief server process. To stop the chief server, run `chief kill`, and to check its status, `chief status`.

## Quickstart

### The Pipeline File

Once you have the Chief server installed, create a file called `chief.pipeline` _in the repo_ for which you wish to set up CI/CD. The pipeline file allows you to configure build & deploy steps which you'd like the Chief server to run once you push a release branch. Here's the minimum needed for a pipeline file:

```
- INFO -
repo: <insert git repo url here>

- BUILD PHASE -
echo "building..."

- TEST PHASE-
echo "testing..."

- DEPLOY PHASE -
echo "deploying"...
```

A `chief.pipeline` file has 4 sections: info, build phase, test phase, and deploy phase. Underneath the `info` section must be a `repo:` followed by the url of your git repo. Once the Chief server is running your new pipeline, this will be the repo that it pulls release branches from.

The build/test/deploy phase sections simply contain shell commands to execute once Chief has cloned your release branch. These commands will be executed in the top-level directory of the repo.

### Creating a new pipeline

Now that you've written a pipeline file, make sure to commit it to your repo. Now, back on your server, create the pipeline by running the following command in the repo:

`chief pipeline create`

This command will look for a `chief.pipeline` file, and the Chief server will begin polling your repo for updates. 

### Build branches

By default, the Chief server will only build branches with the following naming conventions:

`build/*`

`test/*`

`deploy/*`

These branch names correspond to the 3 pipeline phases, and only execute up until their respective phase. For example, maming a branch `test/feature-xxxx` will only run the `build phase` and the `test phase`. However, naming a branch `deploy/1.0.0` will run all 3 pipeline phases. It is recommended that _deploy_ branches are followed by a version number, but it is not required.


## Pipeline Management

`chief pipeline list` - Lists all pipelines

`chief pipeline status` - Lists all pipelines & their current status

`chief pipeline create` - Creates a new pipeline for a repo

`chief pipeline remove <url>` - Removes the pipeline from a repo

`chief pipeline stop <url>` - Stops a running pipeline

`chief pipeline start <url>` - Starts a stopped pipeline


## Access Management

Their are 3 security clearance levels in a Chief server: Master, Normal and Viewer. Master clearance is required to start/kill the Chief server and add/remove users. Normal clearance users can manage pipelines. Viewers only have access to the `status` commands.


`chief access add <username>` - Add a new user with Viewer clearance. Requires Master clearance to execute.

`chief access remove <username>` - Remove a user from the Chief server. Requires Master clearance to execute.

`chief access modify <username> <clearance>` - Upgrades a user to the given clearance. Example: `chief access modify jimbo master` will give the user `jimbo` master clearance.

## Github Releases

## Configuration Options
