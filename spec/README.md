# MyCIlium YAML spec

The YAML spec is quite loosely defined in order to give freedom to implementations.

A basic yaml file looks like this:

```yaml
platform-name:
  steps:
    - sh install-dependencies.sh
    - sh build.sh
    - sh test.sh

platform-name2:
  steps:
    - ...
```

The first required field that all implementations must have is the a platform identifier (first line).
The host machine will look for its own platform.

The second required field is the `steps` which specifies what the host machine should execute. Any number
of other fields can be specified, which is implementation-specific.

For example, this might be valid in an implementation which uses Docker:

```yaml
linux:
  dockerfile: ci/linux.dockerfile
  steps:
    - sh build.sh
    - sh test.sh
```
