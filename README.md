# semantic-release

This project aims to automatically upgrade the release version based on git tags and commit messages.

## Why

It is a way of programming complexity abstraction. So the programmers do not have to update the release versions manually on CHANGELOG.md, setup.py files, generate and push git tags.
This project does this work automatically.

## How

Its Docker image must be called within a `.gitlab-ci.yml` Continuous Integration (CI) stage as follows.

```
TODO
```

## Required Environment Variables

- Set the GIT_HOST with your git host.

    I.e.:
    ```
    GIT_HOST=gilab.companyname.com
    ```
- Set the GIT_HOST with the gitlab project group name.
    
    I.e.:
    ```
    GIT_GROUP=developmentGroup
    ```
- Set the GIT_PROJECT with the gitlab project name.
    
    I.e.:
    ```
    GIT_PROJECT=projectName
    ```


## Testing

Run:

```
make check
```
Open generated coverage on a browser:

```
make coverage
```
To perform static analysis:

```
make analyze
```

## Releasing

Run:

```
make release version=<version>
```

It will create a git tag with the provided **<version>**
and build and publish a docker image.

## Git Hooks

To install the project githooks run:

```
make githooks
```
