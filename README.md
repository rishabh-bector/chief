# Chief

Chief is an elegant continuous integration & deployment server, for hobbyists.

## Installation

Run the following commands on a linux server:

1. `insert wget/curl command here`

2. `chief setup`

3. `chief start`

The `chief setup` command should only have to be run once, on installation. You'll be prompted to enter a username/password, to create your Chief user. This first user will have master clearance, so it can add/remove other users and manage the Chief server process. 

The `chief start` command starts the chief server process. To stop the chief server, run `chief kill`, and to check its status, `chief status`.

## Quickstart

Once you have the Chief server installed, create a file called `chief.pipeline` in the repo for which you wish to set up CI/CD. The pipeline file allows you to configure build & deploy steps which you'd like the Chief server to run once you push a release branch. Here's the minimum needed for a pipeline file:

```
- INFO -
repo: <insert git repo url here>

- BUILD PHASE -
echo "building..."

- DEPLOY PHASE -
echo "deploying"...
```

A `chief.pipeline` file has 3 sections: info, build phase, and deploy phase. Underneath the `info` section must be a `repo:` followed by the url of your git repo. 


## Access management

`chief access add <username>`

`chief access remove <username>`
