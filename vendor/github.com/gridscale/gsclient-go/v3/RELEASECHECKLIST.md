* Branch release branch of develop
* Finalise changelog
* Bump version number (if it exists)
* Push release branch (in case we want to do hotfixes later)
* Merge --no-ff release branch onto master
 git merge --no-ff release/release-version_number
 git push upstream new-release-version_number
* create release on master with tag following the naming scheme for the project
* Merge --no-ff release branch onto develop
git merge --no-ff release/release-version_number
git push upstream new-release-version_number
* Unlock changelog in develop