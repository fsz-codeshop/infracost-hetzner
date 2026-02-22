terraform {
  required_providers {
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = "~> 1.45"
    }
  }
}

provider "hcloud" {
  # dummy token as we just want a plan - must be 64 chars
  token = "1234567890123456789012345678901234567890123456789012345678901234"
}

resource "hcloud_server" "web" {
  name        = "web-server"
  server_type = "cx22"
  image       = "ubuntu-22.04"
  location    = "nbg1"
}

resource "hcloud_volume" "data" {
  name      = "database-data"
  size      = 50
  server_id = hcloud_server.web.id
  location  = "nbg1"
}
