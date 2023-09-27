data "looker_public_ip_addresses" "ips" {
  
}

resource "local_file" "file" {
  filename = "test.txt"
  content = data.looker_public_ip_addresses.ips.public_ips[0]
}