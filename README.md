# go-uaac

[![wercker status](https://app.wercker.com/status/578d939c4223764a55e46b2b2079df34/s/master "wercker status")](https://app.wercker.com/project/byKey/578d939c4223764a55e46b2b2079df34)

A client library for the [UAA service](https://github.com/cloudfoundry/uaa) written in Go.

## Usage

This project makes use of the Command pattern. The various command implementations can be found in the relevant subpackage folders. For example:

```
config := &uaa.ClientConfig{
  ApiAddress: "http://localhost:8080/uaa",
  ClientID: "admin",
  ClientSecret: "adminsecret"
}

uaac, err := uaa.NewClient(clientConfig)
if err != nil {
  return nil, fmt.Errorf("Failed to initialize uaa client; %v", err)
}

var serverInfo ServerInfo
command := NewGetServerInfoCommand(uaac, &serverInfo)
if err := command.Execute(); err != nil {
  t.Errorf("Failed to get ServerInfo: %v", err)
}

fmt.Println("Your UAAC Server Version is ", serverInfo.Version())
```


## Develop

This project uses [Glide](https://github.com/Masterminds/glide)

To setup your local workspace, first clone this project, and then run `glide install`

To run the project test suite, run `go test $(glide novendor)`

## Resources

* https://github.com/cloudfoundry/uaa
* https://docs.cloudfoundry.org/api/uaa/
