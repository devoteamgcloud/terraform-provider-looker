resource "looker_project_git_repo" "myrepo-project" {
  project_id       = "project-id"
  git_remote_url   = "git@github.com:workspace/repo.git"
  git_service_name = "github"
}
