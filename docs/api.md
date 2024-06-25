This page lists the API endpoints available to clients querying the Commune server. All requests are made at a public server at `https://public.commune.sh`.

### GET /publicRooms

This API returns top level public rooms.

##### Example Request:
`GET https://public.commune.sh/publicRooms`
##### Example Response.
```json
{
  "chunk": [
    {
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh",
      "room_type": "m.space",
      "canonical_alias": "#commune:commune.sh",
      "join_rule": "public",
      "world_readable": true,
      "name": "Commune",
      "num_joined_members": 26,
      "guest_can_join": true
    }
  ],
  "total_room_count_estimate": 1
}
```

### GET /rooms/{room_id}/public

This API returns returns the public state of a room. It only returns `true` for public rooms, and `false` for everything else. A `false` is not confirmation that the room actually exists.

##### Example Request:
`https://public.commune.sh/rooms/!rtrpurCAEOKmDDKrPs:commune.sh/public`
##### Example Response.
```json
{
  "public": true
}
```

### GET /rooms/{room_id}/state

This API returns returns the state summary of a public room.

##### Example Request:
`https://public.commune.sh/rooms/!rtrpurCAEOKmDDKrPs:commune.sh/state`
##### Example Response.
```json
{
  "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh",
  "name": "Commune",
  "canonical_alias": "#commune:commune.sh",
  "topic": null,
  "avatar_url": "",
  "join_rule": "public",
  "room_type": "m.space",
  "guest_can_join": true,
  "world_readable": true,
  "num_joined_members": 26
}
```

### GET /rooms/{room_id}/state_events

This API returns returns the current state events of a public room.

##### Example Request:
`https://public.commune.sh/rooms/!rtrpurCAEOKmDDKrPs:commune.sh/state_events`
##### Example Response.
```json
{
  "current_state_events": [
    {
      "type": "m.room.create",
      "sender": "@commune:commune.sh",
      "content": {
        "type": "m.space",
        "room_version": "10",
        "creator": "@commune:commune.sh"
      },
      "state_key": "",
      "origin_server_ts": 1689387184301,
      "unsigned": {
        "age_ts": 1689387184301
      },
      "event_id": "$JqicaDASyzav2GkXPHBqidYfbFCAMWteGbfkGAeOKxA",
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh"
    },
    {
      "type": "m.room.canonical_alias",
      "sender": "@commune:commune.sh",
      "content": {
        "alias": "#commune:commune.sh"
      },
      "state_key": "",
      "origin_server_ts": 1689387184738,
      "unsigned": {
        "age_ts": 1689387184738
      },
      "event_id": "$uFUWXYvttR0Imzs9l6LC7CmYSUM_bN1u6bcmQY8cOdo",
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh"
    },
    {
      "type": "m.room.guest_access",
      "sender": "@commune:commune.sh",
      "content": {
        "guest_access": "can_join"
      },
      "state_key": "",
      "origin_server_ts": 1689387184742,
      "unsigned": {
        "age_ts": 1689387184742
      },
      "event_id": "$RXSLSWe8s8ugRwFWsnljiAo10In-5tQFXGbulEPnBV4",
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh"
    },
    {
      "type": "m.room.history_visibility",
      "sender": "@commune:commune.sh",
      "content": {
        "history_visibility": "world_readable"
      },
      "state_key": "",
      "origin_server_ts": 1689387184741,
      "unsigned": {
        "age_ts": 1689387184741
      },
      "event_id": "$T4gaIn7dQKxjm62LxK7thCVF03_vuTnp2PR59pNNKUM",
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh"
    },
    {
      "type": "m.room.join_rules",
      "sender": "@commune:commune.sh",
      "content": {
        "join_rule": "public"
      },
      "state_key": "",
      "origin_server_ts": 1689387184739,
      "unsigned": {
        "age_ts": 1689387184739
      },
      "event_id": "$MQZLDKDjOeIy8fS4wPLyzo-FLdJwBZVH8mQzmDcbXQs",
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh"
    },
    {
      "type": "m.room.avatar",
      "sender": "@ah:commune.sh",
      "content": {
        "url": ""
      },
      "state_key": "",
      "origin_server_ts": 1689388887530,
      "unsigned": {
        "age_ts": 1689388887530,
        "replaces_state": "$QvHZIUOzNXDjMEQkliofnbeCKnci_V1h229K4sNtrpM"
      },
      "event_id": "$Utl4xW0tW3DcLFSxKWL_8tzEzZX1DgVycN-OHQJR7No",
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh"
    },
    {
      "type": "m.room.power_levels",
      "sender": "@commune:commune.sh",
      "content": {
        "ban": 60,
        "events": {
          "m.room.create": 10,
          "m.room.name": 60,
          "m.room.power_levels": 100,
          "m.space.child": 10,
          "m.space.parent": 10
        },
        "events_default": 10,
        "invite": 0,
        "kick": 60,
        "notifications": {
          "room": 20
        },
        "redact": 10,
        "state_default": 10,
        "users": {
          "@commune:commune.sh": 100,
          "@ah:commune.sh": 100
        },
        "users_default": 10
      },
      "state_key": "",
      "origin_server_ts": 1689736928682,
      "unsigned": {
        "age_ts": 1689736928682,
        "replaces_state": "$nX7G8RZj2cbY0vpAuOvAqQdQ0MgPY5juKyy_RJu30lQ"
      },
      "event_id": "$HBxrIHVo9jqLkyDHKlcGwhMkmVt_hE-nc6XU161RZjs",
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh"
    },
    {
      "type": "m.space.type",
      "sender": "@ah:commune.sh",
      "content": {
        "type": "chat"
      },
      "state_key": "",
      "origin_server_ts": 1690631145454,
      "unsigned": {
        "age_ts": 1690631145454,
        "replaces_state": "$5UbUUsnDE1yN4QpBE4lr2dXg6EhUXBnWFUCPvyMsiAc"
      },
      "event_id": "$I_C2AX02LWwK2cbyU425iNk0Ji6pbt4PVcTrTnuEBDY",
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh"
    },
    {
      "type": "commune.room.public",
      "sender": "@ah:commune.sh",
      "content": {
        "public": true
      },
      "state_key": "",
      "origin_server_ts": 1718617329550,
      "unsigned": {
        "age_ts": 1718617329550
      },
      "event_id": "$Q5Ec4kXPd21BIrETTefvV7yrk_qmAmcEVs6yUgJIdAw",
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh"
    },
    {
      "type": "m.room.name",
      "sender": "@ah:commune.sh",
      "content": {
        "name": "Commune"
      },
      "state_key": "",
      "origin_server_ts": 1718617548259,
      "unsigned": {
        "age_ts": 1718617548259,
        "replaces_state": "$kSCV1vLZ3qBX-nt4a-jLU4fahmW-uN3Ft8UYVBLNSbQ"
      },
      "event_id": "$OaDmAn4mJDFCPYZ-bbs-8UXnGAJ7RxpYnP-PCBH5XaA",
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh"
    },
  ]
}
```


### GET /rooms/{room_id}/hierarchy

This API returns returns the hierachy of a public space, listing all it's
descendants.

##### Example Request:
`https://public.commune.sh/rooms/!rtrpurCAEOKmDDKrPs:commune.sh/hierarchy`
##### Example Response.
```json
{
  "rooms": [
    {
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh",
      "room_type": "m.space",
      "canonical_alias": "#commune:commune.sh",
      "join_rule": "public",
      "world_readable": true,
      "name": "Commune",
      "num_joined_members": 26,
      "guest_can_join": true,
      "children_state": [
        {
          "type": "m.space.child",
          "state_key": "!zVxyHsWDXPQCUvfOdA:commune.sh",
          "content": {
            "suggested": true,
            "via": [
              "commune.sh"
            ]
          },
          "sender": "@spaenny:boehm.sh",
          "origin_server_ts": 1712267509679
        }
      ]
    },
    {
      "room_id": "!zVxyHsWDXPQCUvfOdA:commune.sh",
      "room_type": "",
      "canonical_alias": "#yeeg6Too8Reegea4rae5oibiequeephe:commune.sh",
      "join_rule": "public",
      "world_readable": true,
      "name": "development",
      "num_joined_members": 12,
      "guest_can_join": false
    }
  ]
}
```

### GET /rooms/{room_id}/members

This API returns returns the members of a public room.

##### Example Request:
`https://public.commune.sh/rooms/!rtrpurCAEOKmDDKrPs:commune.sh/members`
##### Example Response.
```json
{
  "chunk": [
    {
      "type": "m.room.member",
      "sender": "@commune:commune.sh",
      "content": {
        "membership": "join",
        "displayname": "commune"
      },
      "state_key": "@commune:commune.sh",
      "origin_server_ts": 1689387184518,
      "unsigned": {
        "age_ts": 1689387184518
      },
      "event_id": "$xwN1k5w_HN52ruF3i5sxzDO7yGb9Cq3ymYP52-lMeNo",
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh"
    },
    {
      "type": "m.room.member",
      "sender": "@ah:commune.sh",
      "content": {
        "membership": "join",
        "displayname": "ah"
      },
      "state_key": "@ah:commune.sh",
      "origin_server_ts": 1689387470815,
      "unsigned": {
        "age_ts": 1689387470815
      },
      "event_id": "$F6WWUSy8gX0-FXHci4wcW3XFgGvHkB10FnP7jT7vuV4",
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh"
    },
    {
      "type": "m.room.member",
      "sender": "@kfjdkfjkdfkjf:shpong.com",
      "content": {
        "avatar_url": null,
        "displayname": "kfjdkfjkdfkjf",
        "membership": "join"
      },
      "state_key": "@kfjdkfjkdfkjf:shpong.com",
      "origin_server_ts": 1689712219241,
      "unsigned": {
        "age": 6
      },
      "event_id": "$PoVPcRWkJ_-OA8aYhEokE9eBsT54Ov5hTBkNPpAnl5o",
      "room_id": "!rtrpurCAEOKmDDKrPs:commune.sh"
    }
  ]
}
```

