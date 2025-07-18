# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v0.0.4

- Fix stdin seeking issue in upload command by buffering stdin data before upload

## v0.0.3

- Modify upload command to read content from stdin instead of static string

## v0.0.2

- Add typed S3 credentials (URL, AccessKey, SecretKey) with String() methods
- Change package name from `s3` to `s3utils` to avoid conflicts with AWS SDK
- Update all commands to use new typed credentials

## v0.0.1

- Initial Version
