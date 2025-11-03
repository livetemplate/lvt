# Changelog

All notable changes to the LVT CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.0] - 2025-11-03

Initial release of LVT CLI as a standalone package extracted from the LiveTemplate monorepo.

### Features

- **Code Generation**
  - Generate CRUD resources with models, handlers, and views
  - Generate standalone views
  - Generate complete applications
  - Field parsing and validation

- **Kit System**
  - Built-in kits: Tailwind CSS, Bulma, Pico CSS, None
  - Kit cascade: Project → User → System
  - ~60 CSS helper methods per kit
  - Component templates
  - Generator templates
  - Kit creation and customization tools
  - Kit validation

- **Development Server**
  - Hot reload via WebSocket
  - File watching with fsnotify
  - Automatic browser refresh
  - Serves static assets
  - Live template rendering

- **Database Tools**
  - Migration creation and management
  - Seeder creation and execution
  - SQLite and modernc.org/sqlite support
  - Migration status tracking

- **Stack Generators**
  - Docker configurations
  - Systemd service files
  - Deployment scripts

- **Interactive UI**
  - Terminal UI for app creation
  - Resource generator wizard
  - View generator wizard
  - Built with Bubble Tea and Lipgloss

- **Testing Utilities**
  - E2E test helpers
  - Chromedp integration for browser testing
  - Test server utilities
  - Golden file testing

### Infrastructure

- **Release Automation**: Automated release script with version synchronization
- **CI/CD**: GitHub Actions workflows
- **Pre-commit Hooks**: Go formatting, linting, and testing
- **GoReleaser**: Multi-platform binary builds
- **Version Tracking**: VERSION file for release management

### Documentation

- Complete README with examples
- Contributing guidelines
- Version synchronization strategy with core library

### Related Versions

- Core Library: v0.1.0
- Client Library: v0.1.0
- Examples: v0.1.0

---

## Version Synchronization

LVT follows the LiveTemplate core library's major.minor version (X.Y):

- Patch versions (X.Y.Z) are independent
- Minor/major versions must match core library
- See README.md for details
