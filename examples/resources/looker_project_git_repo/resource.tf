resource "looker_project_git_repo" "myrepo-project" {
  project_name     = "project-name"
  git_remote_url   = "git@github.com:workspace/repo.git"
  git_service_name = "github"
}
