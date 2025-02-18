## Clone submodules recursively by manual since the paths in gitmodules are not defined relatively
## See More: https://docs.gitlab.com/ee/ci/git_submodules.html
make update-submodules

## Clear used containers
docker container rm -f build-env