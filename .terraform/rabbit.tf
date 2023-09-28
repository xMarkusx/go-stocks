terraform {
  required_providers {
    rabbitmq = {
      source = "cyrilgdn/rabbitmq"
      version = ">=1.8.0"
    }
  }
}

provider "rabbitmq" {
  endpoint = "http://rabbit:15672"
  username = "guest"
  password = "guest"
}

resource "rabbitmq_user" "test" {
  name     = "test"
  password = "foobar"
  tags     = ["management"]
}

resource "rabbitmq_permissions" "test" {
  user  = rabbitmq_user.test.name

  permissions {
    write     = ""
    read      = ""
    configure = ""
  }
}

resource "rabbitmq_queue" "portfolio" {
  name     = "portfolio"
  settings {
    durable     = true
    auto_delete = false
    arguments = {
      x-queue-type = "stream"
    }
  }
}