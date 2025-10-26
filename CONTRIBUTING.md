# Contributing to BLEF

Thank you for your interest in contributing to BLEF! 🎉

BLEF is a community-driven open standard, and we welcome contributions from everyone. Whether you're fixing a typo, suggesting a feature, or implementing a converter tool, your input is valuable.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Process](#development-process)
- [Specification Changes](#specification-changes)
- [Submitting Changes](#submitting-changes)
- [Style Guidelines](#style-guidelines)
- [Community](#community)

## 📜 Code of Conduct

This project adheres to a [Code of Conduct](./CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## 🤝 How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates.

When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the expected behavior** and what actually happened
- **Provide specific examples** (BLEF files, code snippets)
- **Describe your environment** (OS, tool version, etc.)

Use the [Bug Report template](.github/ISSUE_TEMPLATE/bug_report.md).

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description** of the proposed feature
- **Explain why this enhancement would be useful** to most BLEF users
- **Provide examples** of how it would work

Use the [Feature Request template](.github/ISSUE_TEMPLATE/feature_request.md).

### Improving Documentation

Documentation improvements are always welcome! This includes:

- Fixing typos or clarifying wording
- Adding examples or use cases
- Translating documentation to other languages
- Improving the specification clarity

### Building Tools & Converters

We encourage the development of:

- Platform-specific converters (Goodreads, Babelio, etc.)
- Validation tools
- Libraries in various programming languages
- Plugins for book management software

Share your implementations via PR or discussion!

### Implementing BLEF Support

If you're adding BLEF support to your application:

1. Follow the [specification](./blef_specification.md)
2. Validate against the [JSON Schema](./schema/blef-schema-v0.1.0.json)
3. Test with the provided examples
4. Share your implementation for community feedback
5. Add your tool to the README implementations list

## 🔧 Development Process

### 1. Fork & Clone

```bash
git clone https://github.com/yoanbernabeu/BLEF.git
cd BLEF
```

### 2. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

Branch naming conventions:
- `feature/` — New features or enhancements
- `fix/` — Bug fixes
- `docs/` — Documentation changes
- `tool/` — New tools or converters

### 3. Make Your Changes

- Follow the style guidelines
- Add or update tests if applicable
- Update documentation as needed
- Validate JSON files against the schema

### 4. Commit Your Changes

```bash
git add .
git commit -m "feat: add support for custom metadata fields"
```

Use conventional commits:
- `feat:` — New feature
- `fix:` — Bug fix
- `docs:` — Documentation changes
- `style:` — Formatting, missing semicolons, etc.
- `refactor:` — Code restructuring
- `test:` — Adding tests
- `chore:` — Maintenance tasks

### 5. Push & Create a Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub using the PR template.

## 📝 Specification Changes

Changes to the BLEF specification require special consideration:

### Minor Changes (Patch)
- Typo fixes
- Clarifications
- Example improvements
- **No PR required** for trivial fixes

### New Optional Fields (Minor)
- Adding optional fields
- Extending enums
- New identifier types
- **Requires discussion** and community feedback

### Breaking Changes (Major)
- Changing required fields
- Removing fields
- Changing data types
- **Requires RFC** (Request for Comments) and broader consensus

### Process for Significant Changes

1. **Open a Discussion** — Describe the problem and proposed solution
2. **Gather Feedback** — Allow time for community input
3. **Create RFC** — Write a formal proposal if needed
4. **Implementation** — Code the changes after consensus
5. **Update Version** — Follow semantic versioning

## 🎨 Style Guidelines

### JSON Files

- Use **2 spaces** for indentation
- Use **double quotes** for strings
- Include trailing commas where allowed
- Sort fields alphabetically when possible
- Validate all JSON against the schema

Example:
```json
{
  "format": "BLEF",
  "version": "0.1.0",
  "exported_at": "2025-10-26T14:00:00Z"
}
```

### Markdown Files

- Use ATX-style headers (`#`)
- Use reference-style links for readability
- Include code blocks with language tags
- Keep lines under 100 characters when possible

### Code (if contributing tools)

- Follow the conventions of the language you're using
- Include comments for complex logic
- Write tests for your code
- Document public APIs

## 📨 Submitting Changes

### Pull Request Process

1. **Fill in the PR template** completely
2. **Link related issues** using keywords (fixes #123)
3. **Ensure all checks pass** (validation, linting)
4. **Request review** from maintainers
5. **Address feedback** promptly
6. **Keep commits clean** (squash if needed)

### What to Include

- **Clear description** of changes
- **Motivation** for the change
- **Breaking changes** (if any)
- **Examples** demonstrating the change
- **Updated documentation**

### Review Criteria

PRs will be evaluated on:

- ✅ **Adherence to specification** and goals
- ✅ **Code quality** and style
- ✅ **Documentation** completeness
- ✅ **Backward compatibility** (when possible)
- ✅ **Test coverage** (for tools)
- ✅ **Community impact** and usefulness

## 💬 Community

### Where to Ask Questions

- **GitHub Discussions** — General questions, ideas, showcase
- **GitHub Issues** — Bug reports, feature requests
- **Email** — For private concerns

### Recognition

Contributors will be:
- Listed in CHANGELOG for significant contributions
- Credited in release notes
- Mentioned in documentation where applicable

## 🙏 Thank You!

Your contributions make BLEF better for everyone. Whether you're:

- 🐛 Reporting a bug
- 💡 Suggesting a feature
- 📝 Improving docs
- 💻 Writing code
- 🌍 Translating content
- 📢 Spreading the word

**You're helping build an open, interoperable future for personal book libraries!**

---

*Happy contributing! 📚✨*

