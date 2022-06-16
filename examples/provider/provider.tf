terraform {
  required_providers {
    looker = {
      source  = "devoteamgcloud/looker"
      version = "0.0.0-dev"
    }
  }

}

provider "looker" {
  base_url      = "https://org.cloud.looker.com:19999/api/" # Optionally use env var LOOKER_BASE_URL
  client_id     = ""                                        # Optionally use env var LOOKER_API_CLIENT_ID
  client_secret = ""                                        # Optionally use env var LOOKER_API_CLIENT_SECRET
}
