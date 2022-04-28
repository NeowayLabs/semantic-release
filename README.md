# semantic-release

This project aims to automatically upgrade the release version based on git tags and commit messages.

![semantic](./docs/static/semantic.png)

<br></br>
# Why

It is a way of programming complexity abstraction. So the programmers do not have to update the release versions manually on CHANGELOG.md, setup.py files, generate and push git tags.
This project does this work automatically.

<br></br>
# How

### If you do not have semantic-release image in your registry follow the enumerated steps bellow.
1. Clone this project, generate and push its image:

    Generating the image.
    ```
    make image
    ```

2. Before pushing it to your registry, you must login to the docker registry.

    Note: Replace [] with the correct credentials and registry DNS.

    ```
    docker login -u [username] -p [password] [registry.server.com]
    ```

3. Push the image to the registry.
    ```
    docker push neowaylabs/semantic-release:[version]
    ```

### Once you have semantic-release image in your registry, you can use it on CI files.

The semantic release Docker image can be called within a `.gitlab-ci.yml` file, git hub actions, or travis ci stages as follows.

```yaml
semantic-release:
    stage: semantic-release
    only:
        refs:
            - master
    script:
        - docker run registry.com/group/semantic-release up -git-host $CI_SERVER_HOST -git-group gitGroupNameHere -git-project $CI_PROJECT_NAME -auth sshKeyHere
```

### If your project is a Python project you can add the flag `-setup-py true` to update the release version in this file too.

Note: The version must be placed in a variable called `__version__` as follows:

```py
from setuptools import setup, find_packages

__version__ = "1.0.0"
with open("README.md") as description_file:
    readme = description_file.read()

with open("requirements.txt") as requirements_file:
    requirements = [line for line in requirements_file]

setup(
    name="helloworld",
    version=__version__,
    author="dataplatform",
    python_requires=">=3.8.5",
    description="A sample project to test semantic-release automations.",
    long_description=readme,
    url="https://gitlab.com.br/group/py_project",
    install_requires=requirements,
    packages=find_packages(),
)
```

## CLI Help
If you need more information about the semantic release CLI usage you can run the following command.

```
docker run registry.com/group/semantic-release help
```

### If you want to check the commit tags you can use in your commit message run the following command.

```
docker run registry.com/group/semantic-release help-cmt
```

## Commit Messages Pattern
So the semantic release can find out the commit type to define the upgrade type (MAJOR, MINOR or PATCH), and the message to write to CHANGELOG.md file, one must follow the commit message pattern bellow:


```
type: [type here].
message: Commit message here.
```

I.e.
```
type: [feat]
message: Added new function to print the Fibonacci sequece.
```

### If you want to complete a Merge Request without triggering the versioning process then you can use one of the skip type tags as follows.

- type: [skip]
- type: [skip v]
- type: [skip versioning]

<br></br>
# Development Guides
## Environment Variables

Before running the integration tests you must set up the bellow environment variable with the SSH PRIVATE KEY.

```
SSH_INTEGRATION_SEMANTIC
```

You can find it on the CI engine such as, travis ci or github actions.


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
make static-analysis
```

## Integration Tests

Set up the required environment:
```
make env
```

Run the integration tests as soon as the gitlab container is available:
```
make check-integration
```

You can remove the environment running the following command:
```
make env-stop
```

## Local Running
You can also run the application locally by running the following commands:

Create a go binary file and run it:
```
make run-local
```

Run with docker:
```
make image
```

```
make run-docker-local
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