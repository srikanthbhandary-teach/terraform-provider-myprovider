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

}
output "myprovider_users" {
  value = data.myprovider_users.example.users
}

resource "myprovider_user" "user1" {
  name = "user1"
  age  = 50
  id   = 400
}
