resource "looker_alert" "default" {
  cron = "0 1 * * *"
  description = "A description"
  comparison_type = "GREATER_THAN_OR_EQUAL_TO"
  dashboard_element_id = 218
  destinations {
    destination_type = "EMAIL"
    email_address = "test@test.com"
  }
  field {
    title = "confirmed cases"
    name = "population_data.sum_of_population"
  }
  owner_id = 60
  treshold = 10
  is_disabled = false
  is_public = true
}