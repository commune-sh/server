### What is it? 

Commune makes matrix servers publicly accessbile in the open web.

This server runs next to a Synapse Matrix server as a proxy, and provides API endpoints to read data from public rooms. This enables [matrix clients](https://github.com/commune-os/client) to implement publicly accessible chat rooms that can be viewed by non-authenticated users, allowing them to be indexed by search engines and ensuring communities and their information can exist in the public.

This server aims to provide a subset of read-only endpoints from the Matrix [client-server API](https://spec.matrix.org/v1.10/client-server-api/). 

A live version of this server is running on [public.commune.sh](https://public.commune.sh/publicRooms), querying our existing homeserver on `commune.sh`. You may try out the API on the live instance to test how it works and view what data it returns.

### Requirements
Commune only supports Synapse right now - but support for Conduit is on our roadmap. In the future, this server may run as an appservice, so that it's no longer tied to a particular homeserver implementation.

### How to run it
It is best to build this project from source. The only requirement is an existing [Synapse](https://element-hq.github.io/synapse/latest/setup/installation.html) instance, and [Golang](https://go.dev/doc/install).

Follow these steps to run this server:

1. `git clone` this repository.
2. Run `make deps` to install dependencies.
3. Run `make` to build the binary in `./bin/commune`.
4. Copy `config-sample.toml` to `config.toml` and update the configuration.
5. Run the binary as a [systemd service](https://github.com/commune-os/server/blob/main/docs/commune-server.service) and put it behind a [reverse proxy](https://github.com/commune-os/server/blob/main/docs/nginx-reverse-proxy).

Example configurations for the systemd unit and an nginx reverse proxy are included in the [/docs](https://github.com/commune-os/server/blob/main/docs) directory.

Additionally, you can run the server in a container with `docker compose up -v`. Note that this has only been minimally tested, and will change as soon as we work out a proper deployment strategy.

Run this server on the same VPS hosting your Synapse server. You can run it on a different host, but that would require opening up Postgres access to the public network. It is not recommended at this stage. If you do run it on a different host, consider settings up a private network. 


#### Configuration
A `config.toml` file is required to run Commune. A sample configuration file `config-sample.toml` is included. 

```toml
[app]
# The domain pointing to this server
domain = "public.commune.sh"
# The port the server will listen on
port = 8989

[matrix]
# Local domain of the Synapse server
homeserver = "localhost:8008"
# The server_name part of your Synapse configuration
server_name = "commune.sh"
# DB connection string for the Synapse database
db = "postgres://commune:password@host.docker.internal:5432/synapse?sslmode=disable"

[security]
# This should include the domain you'll run this server on, and any other
# domains you want to allow for local development
allowed_origins = ["http://public.commune.sh"]

[log]
max_size = 100
max_backups = 7
max_age = 30
compress = true

[capabilities.public_rooms]
# List public server capabilities. Clients can query this endpoint to see what the server supports.
list_rooms = true
view_hierarchy = true
read_messages = true

[cache]
# Cache configuration
public_rooms = true
```

#### Development
Developing this server requires a locally running Synapse server. Running `make deps` sets up the `modd` command. You can then run `modd` to run the server. It will watch for changes, rebuild and rerun the binary.

#### API Endpoints

We currently implement the following API endpoints:
 - [x] Query public rooms - `/publicRooms`
 - [x] Query room public state - `/rooms/{room_id/public`
 - [x] Query room state - `/rooms/{room_id/state`
 - [x] Query room current state events - `/room/{room_id/state_events`
 - [x] Query space hierarchy - `/room/{room_id}/hierarchy`
 - [x] Query room messages - `/room/{room_id}/messages`
 - [x] Search room messages - `/search`
 - [ ] Sync room events - `/sync`

These are all public endpoints and don't require authorization.

#### What is a public Room
A public room in Commune is a matrix room that has the following [state events](https://spec.matrix.org/v1.10/client-server-api/#types-of-room-events):
- `history_visibility` set to `world_readable`
- `guest_access` set to `can_join`
- `join_rules` set to public`
- `canonical_alias` exists
- A state event of the type `commune.room.public` with content `{"public": true}`
- Room must be local to the homeserver, not federated

Rooms without these state events cannot be queried. This means that simply running this server will not automatically allow access to rooms in a matrix homeserver. The room owners must explicitly include state events to make them publicly accessible. This provides a good balance between the need to have public rooms, while ensuring that every room doesn't become public by default.

This server does not read data from private or encrypted rooms, including DM rooms. Further limits will likely be added to ensure sensitive room data cannot be publicly accessed.


### Security Considerations
This server connects directly to the Synapse database, which gives it access to all the data belonging to the homeserver. We recommend you not run this on your primary matrix instance just yet. Consider setting up a secondary matrix server with the intention of it being a fully open server.