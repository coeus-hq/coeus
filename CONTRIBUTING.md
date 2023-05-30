# Contributing Guide

Thank you for considering contributing to our project! This guide outlines how to contribute in a way that is efficient and consistent with the project's standards.

## Contributing

While you're free to add features that aren't listed in the project's issues, your chances of having a pull request accepted will be higher if you consult a repository manager beforehand.

### New to the Project?

If you're interested in contributing but don't know where to start, feel free to schedule a time with repository manager Collin https://calendly.com/collin-w to learn more about the codebase. He will provide you with a helpful overview and answer any questions you might have.

## Naming Conventions

### Branching

When creating a branch off of 'dev', please follow the naming convention: `[name]-[action]`. For example, `collin-button-fix`.

### JavaScript

For functions and variables, use camelCase. For example, `myFunction`.

### HTML

CSS classes and IDs should follow kebab-case. For example, `my-class` or `my-id`.

### Golang

For variable and function names, please use camelCase. Exported functions should use PascalCase as per Go's convention or the function will not export.

### SQLite3

In SQLite3 syntax, use underscores between words. For example, `my_table`.

### Image Files

Image file names should follow the format: `[what]-[type]-[subtype].[extension]`. For example, `icon-arrow-up.svg`. SVG files are preferred when possible.

## Before Submitting a Pull Request

Before making a pull request into 'dev', please make sure to check the following:

- Test your code to ensure it behaves as expected and contains no bugs.
- Ensure that model tests pass with the default seed data. From the root run ``` cd models && go test -v && cd .. ```
- Confirm that there are no browser console errors.
- Check that there are no errors in the server.
- Verify that your code follows the proper naming conventions outlined above.

By following these guidelines, you can help make the review process go smoothly, and increase the likelihood that your contribution will be accepted. Thank you for your time and effort!


