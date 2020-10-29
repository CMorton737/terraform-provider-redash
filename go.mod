module github.com/snowplow-devops/terraform-provider-redash

go 1.15

replace github.com/snowplow-devops/redash-client-go => ../redash-client-go

require (
	github.com/Jeffail/gabs v1.4.0 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.1
	github.com/mitchellh/gox v1.0.1 // indirect
	github.com/snowplow-devops/redash-client-go v0.4.2
	golang.org/x/tools v0.0.0-20201030204249-4fc0492b8eca // indirect
)
