# dogu-additional-mounts-init Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.1.2] - 2025-06-12
### Fixed
- [#5] Remove the permissions from the binary because the container will be executed with different uids and gids defined from the kubernetes security context. Otherwise, the container can not start.

## [v0.1.1] - 2025-06-02

### Changed

- Naming from `dogu-data-seeder` to `dogu-additional-mounts-init` for consistent naming with the dogu-operator.
- **Attention**: This also changes the artifact name from `dogu-data-seeder` to `dogu-additional-mounts-init`. Keep in mind to update your image references.

## [v0.1.0] - 2025-05-22

### Added

- [#1] Add first version of the dogu data seeder with the copy function for volume mounts.
