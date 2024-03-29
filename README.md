# InternLinter
If you keep opening up pull requests with TODO comments and print statements, use this linter to double check your work.
## Installation
1. Clone this repository: `git clone github.com/bdeleonardis1/InternLinter`
2. cd into the repo: `cd InternLinter`
3. Build intern linter: `go build` (requires go to be installed)
4. Add the generated executable to your path
## Setup
InternLinter requires a config file to work successfully. By default, we will look for this file in `~/InternLinter/config.yaml`, but if want to store it somewhere else you can specify the path with the `--config` flag.
 
The config file must be yaml, however! The required fields in the config file are:
- `github.defaultBase`: the branch you will be opening the PR against
- `github.organization`: the organization the repo exists in
- `github.repository`: the name of the repository

If you are opening a pull request from a forked repo to the repo it was forked from, you will also need to specify:
- `github.isFork`: whether or not the repo is forked
- `github.username`: your username

Other optional parameters are:
- `checkForPrints`: set to false if you don't want the linter to look for new print statements
- `checkForTODOs`: set to false if you don't want the linter to look for new TODO comments

Here is an example of a `config.yaml` file:
```yaml
checkForPrints: true
checkForTODOs: false
github:
	defaultBase: master
	defaultMaintainerCanModify: true
	organization: codebase-berkeley-mentored-project-fa17
	repository: LinterTester
	username: "bdeleonardis1"
	isFork: true
```
Finally, you need to create a Github oauth token ([directions](https://help.github.com/en/articles/git-automation-with-oauth-tokens)) with admin permissions and set the `GITHUBOAUTH` environment variable to this token's value.

## Usage
Inside the repository that you want to open the PR from, while checked out to the branch you want to use to open the PR, run `InternLinter --title "Your PR title"`. `--title` is the only required command line argument. As mentioned above, if your config file is not located at `~/InternLinter/config.yaml` you will need to specify the `--config` argument as well. An example with a different config file is as follows: ```InternLinter --config ../config.yaml --title "Your PR title"```
