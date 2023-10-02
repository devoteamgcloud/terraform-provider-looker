resource "looker_user_attribute" "attribute" {
  name = "attribute"
  label = "attribute label"
  type = "string"
  default_value = "sandbox"
  user_can_edit = false
  user_can_view = false
  value_is_hidden = true
  hidden_value_domain_whitelist = "https://example.com*,localhost:9932,https://www.*.my-destination.com/*"
}