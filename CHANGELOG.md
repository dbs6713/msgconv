Change Log
----------
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org)

## [0.9.0](https://github.com/donbstringham/httpmessageconverter/compare/0.8.0...0.9.0) - 2016-12-08
- Fixed request types to the Accept header
- Version bump

## [0.8.0](https://github.com/donbstringham/httpmessageconverter/compare/0.7.1...0.8.0) - 2016-12-05
- Fixed Content-Type bug in ProtobufHttpMessageConverter
- Version bump

## [0.7.1](https://github.com/donbstringham/httpmessageconverter/compare/0.7.0...0.7.1) - 2016-12-02
- Fixed Content-Type bug in ReadRequestCtx()
- Version bump

## [0.7.0](https://github.com/donbstringham/httpmessageconverter/compare/0.6.0...0.7.0) - 2016-12-02
- Passing in the RequestCtx to the Ctx() methods
- Version bump

## [0.6.0](https://github.com/donbstringham/httpmessageconverter/compare/0.5.0...0.6.0) - 2016-12-01
- Added WriteRequestCtx() method
- Version bump

## [0.5.0](https://github.com/donbstringham/httpmessageconverter/compare/0.4.0...0.5.0) - 2016-12-01
- Added ReadRequestCtx() method
- Version bump

## [0.4.0](https://github.com/donbstringham/httpmessageconverter/compare/0.3.1...0.4.0) - 2016-12-01
- Refactored Read() into ReadRequest() and ReadResponse()
- Refactored Write() into WriteRequest() and WriteResponse()
- Version bump

## [0.3.1](https://github.com/donbstringham/httpmessageconverter/compare/0.3.0...0.3.1) - 2016-12-01
- Refactored Read() to handle Response instead of RequestCtx
- Version bump

## [0.3.0](https://github.com/donbstringham/httpmessageconverter/compare/0.2.1...0.3.0) - 2016-12-01
- Added additional headers to WritePostRequest() and bumped version

## [0.2.1](https://github.com/donbstringham/httpmessageconverter/compare/0.2.0...0.2.1) - 2016-12-01
- Added POST method set and URI set to WritePostRequest()

## [0.2.0](https://github.com/donbstringham/httpmessageconverter/compare/0.1.0...0.2.0) - 2016-12-01
- Added WritePostRequest() method

## [0.1.0](https://github.com/donbstringham/httpmessageconverter/compare/0.1.0...HEAD) - 2016-11-30
- Added more unit-tests
- Initial release

## [Unreleased](https://github.com/donbstringham/go-emris/compare/HEAD) - 2016-11-29
- Initial commit of the code
- This `CHANGELOG.md` file
- `LICENSE` file
- `README.md` file
- `Makefile` for better project control
- `version.go` file in the `core` dir
