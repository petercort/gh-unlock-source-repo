# Unlock Repository

This github cli extension is meant to unlock a repository after a migration.


## Notes

[List migrations API](https://docs.github.com/en/rest/migrations/orgs?apiVersion=2022-11-28#list-organization-migrations)  
[Unlock an organization repository](https://docs.github.com/en/rest/migrations/orgs?apiVersion=2022-11-28#unlock-an-organization-repository)  

https://docs.github.com/en/migrations/overview/about-locked-repositories  
https://github.github.com/enterprise-migrations/#/3.1.2-import-using-graphql-api?id=unlock-imported-repositories  


https://docs.github.com/en/rest/migrations/orgs?apiVersion=2022-11-28  

```json
{
  "data": {
    "unlockImportedRepositories": {
      "migration": {
        "guid": "805e4f0e-325a-49e9-9abd-2ac90a615732",
        "state": "UNLOCKED"
      },
      "unlockedRepositories": [
        {
          "nameWithOwner": "import-test/widgets"
        }
      ]
    }
  }
}
```