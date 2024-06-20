### Commune Server
This server allows for querying public room data from Matrix/Synapse.

It currently implements the following API endpoints:
 - [x] Query public rooms - `/publicRooms`
 - [x] Query room public state - `/rooms/{room_id/public`
 - [x] Query room state - `/rooms/{room_id/state`
 - [x] Query room current state events - `/room/{room_id/state_events`
 - [x] Query space hierarchy - `/room/{room_id}/hierarchy`
 - [x] Query room messages - `/room/{room_id}/messages`
 - [x] Search room messages - `/search`
 - [ ] Sync room events - `/sync`

These are all public endpoints and don't require authorization.

#### Public Room
A public room in Commune is a matrix room that has the following state events:
- `history_visibility` set to `world_readable`
- `guest_access` set to `can_join`
- `join_rules` set to public`
- `canonical_alias` exists
- A state event of the type `commune.room.public` with content `{"public": true}`
- Room must be local to the homeserver, not federated

Rooms without these state events cannot be queried. Further limits will likely be added to ensure sensitive room data cannot be publicly accessed.
