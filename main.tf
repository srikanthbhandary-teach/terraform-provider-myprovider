terraform {
  required_providers {
    myprovider = {
      source = "github.com/srikanthbhandary-teach/myprovider"
    }
  }
}

provider "myprovider" {
  apikey = "myAppSecret12254"
}

data "myprovider_users" "example" {
  filter = {
    id = 1
  }
}
output "myprovider_users" {
  value = data.myprovider_users.example.users
}
