# govh-mail-redirection-manager
GO Client to manage mail redirection from OVH

## What is this ?

This project is a go client to manage email redirection for OVH Provider.

## Prerequisites

Create a configuration file in one of these path :
  - /etc/ovh/ovh.yaml
  - $HOME/.config/ovh/ovh.yaml
  - same direction as binary

The configuration file must have the format below :
```
Endpoint: your_endpoint (ie. ovh-eu)
ApplicationKey: your_application_key
ApplicationSecret: your_application_secret
ConsumerKey: your_consumer_key
Domain: 
  - yourdomain.com
  - yourotherdomain.com
```

To create your APIKeys, follow the link : https://eu.api.ovh.com/createToken/

You have to give the rights below :
  - GET ``/email/domain/*/redirection``
  - POST ``/email/domain/*/redirection``
  - GET ``/email/domain/*/redirection/*``
  - DELETE ``/email/domain/*/redirection/*``

## Build application
```
    go build -o govh-mrm main.go
```

## Commands

To list current redirection
```
    govh-mrm list 
```

You can filter the list by source and destination mail
```
    govh-mrm list --from <mail source> --to <mail destination>
```

To add a new redirection
```
    govh-mrm add --from <mail source> --to <mail destination>
```

To remove a redirection
```
    govh-mrm remove --from <mail source>
```

## Author

Julien Vinet <contact@julienvinet.dev>

## License

The source and documentation in this project are released under the [GNU general public license](https://github.com/julienvinet/govh-mail-redirection-manager/blob/main/LICENSE)