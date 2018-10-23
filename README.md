# go-unattended

A simple update package for Go

*Note: This is a Go implementation of the original [.NET Unattended package](https://github.com/ProjectLimitless/Unattended)*

**Limitations compared to the .NET version**

The .NET version checks and updates DLLs in addition to the primary executable.
Go builds to a single binary and does not need additional DLLs to be present.
Once Go's plugins are supported on all platforms the additional functionality
to check their versions and update independently might be implemented.

go-unattended is still able to update the binary together with any other included
files and folders as with the .NET version.

## Documentation and Getting Started Guide

Available at the [Project Limitless Documentation site](https://docs.projectlimitless.io/unattended). The documentation is for
the .NET version of the package, but the same tutorial can be followed and used
with the sample code in the example folder.
