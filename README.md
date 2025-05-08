# Unlock Source Repository GitHub CLI Extension

This [GitHub CLI](https://cli.github.com/) extension is meant to unlock a migration source repository that was [locked by GitHub Enterpise Importer during a migration](https://docs.github.com/en/migrations/overview/about-locked-repositories#repositories-locked-by-github-enterprise-importer). 
## Installation
```bash
gh extension install robandpdx/gh-unlock-source-repo
```

## Usage
```bash
export GH_SOURCE_PAT="<token>"
gh unlock-source-repo --org <org> --repo <repository>
```

## Example
```bash
gh unlock-source-repo --org mindfulrob --repo java-private-library
```
