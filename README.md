# Timewarrior Sync Server
This repository contains the Server of the Timewarrior Sync project.

# Setup

## Building from source

First you need to build the server:
```sh
go build -o timew-server
```

If you haven't already, create a folder named `authorized_keys` in the same folder as your executable:
```sh
mkdir authorized_keys
```

Now you can start using your server:
```sh
./timew-server start
```

## Using Nix

To install `timew-sync-server` into your current environment, use:

```sh
nix-env -f default.nix -i
```

Then, follow the instructions above (create a directory for the keys
and start the server).

## Using docker

You can build a docker image using the provided `Dockerfile`:
```sh
# Build a docker image tagged timew-server
docker build -t timew-server .
```

To start the server, use:
```sh
# Running from the docker image tagged timew-server
docker run -p 8080:8080 timew-server

# Start an existing docker container
docker start <container-id>
```

Subcommands can be used via `docker exec`:
```sh
docker exec -it <container-id> server <subcommand>
```

# Usage

## Starting the server

Start the server using the `start` subcommand:
```sh
./timew-server start
```

The `start` subcommand supports the following (optional) flags:
- `--config-file`: Specifies the path to the configuration file
- `--port`: Specifies the port. Default: 8080
- `--keys-location`: Specifies the folder holding the authorized keys. Default: `authorized_keys`
- `--no-auth`: Deactivates client authentication. Only for testing purposes.

## Adding users

New users can be registered using the `add-user` subcommand:
```sh
./timew-server add-user
```

If no additional flags are specified, this command will return the new user id.

The `add-user` subcommand supports the following flags:
- `--path`: Specifies the path to a public key and associates it with the user.
- `--keys-loaction`: Specifies the folder holding the authorized keys. Default: `authorized_keys`

**Note:** If you are running our provided Docker image, see the note under `Adding keys to users`.

## Adding keys to users

New keys can be added using the `add-key` subcommand:
```sh
./timew-server add-key --path public-key.pem --id <user-id>
```

The `add-key` subcommand supports the following flags:
- `--path` (**required**): Specifies the path to a public key and associates it with the user.
- `--id` (**required**): Specifies the user id.
- `--keys-location`: Specifies the folder holding the authorized keys. Default: `authorized_keys`

**Note:** If you are running the server inside a docker container, you have to copy the key into the container first:

```sh
# Copy the key into /public-key.pem
docker cp public-key.pem <container-id>:/public-key.pem

# Add the key
docker exec -it <container-id> server add-key --path /public-key.pem --id <user-id>
```

# Development

The code has to be formatted using `go fmt` before commits. To enforce this, we provide a pre-commit-hook. It can be setup by copying it into to your `.git/hooks/pre-commit`:

```sh
cp git/hooks/pre-commit .git/hooks/pre-commit
```
