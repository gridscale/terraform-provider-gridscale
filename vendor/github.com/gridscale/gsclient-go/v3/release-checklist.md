# Release Checklist

This repo uses [git flow](https://nvie.com/posts/a-successful-git-branching-model/) (I'm sorry). Not obligatory but you might want to install [gitflow tools](https://github.com/nvie/gitflow) to help with the commands (It's complicated). Install like `brew install git-flow` or `sudo dnf install gitflow` or whatever is suitable on your system. Make sure to tun `git flow init` once and answer the questions.

On develop branch:

- [ ] Create a release branch from develop: `git flow release start 3.1.0`

This will create and checkout a new release branch. On that release branch:

- [ ] test everything is working as expected
- [ ] update changelog: add additions and fixes and set a release date
- [ ] `git flow release publish 3.1.0`
- [ ] `git flow release finish 3.1.0`
- [ ] create and push the tag `git push origin --tags v3.1.0` (Note: its origin here but it needs to be your remote's name here)

Make sure the new release branch is pushed to the right remote. Then go to GitHub and

- [ ] create two (yes that's two) PRs on GitHub: one from release branch → develop
- [ ] another from release branch → master (make sure you do not accidentally remove the release branch when merging those PRs)

Back on develop branch:

- [ ] create new CHANGELOG.md stub for next release
