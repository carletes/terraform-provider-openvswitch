# Wire dumps of `ovs-vsctl`

## Creating a new bridge

Command-line:

```
# ovs-vsctl add-br pepe0
```

The `ovs-vsctl` client sends a message to fetch the schema for the
`_Server` database:

```
{"id":0, "method": "get_schema", "params": ["_Server"]}
```

The `ovsdb-server` process replies:

```
{
  "id": 0,
  "result": {
    "cksum": "3236486585 698",
    "name": "_Server",
    "version": "1.1.0",
    "tables": {
      "Database": {
        "columns": {
          "name": {
            "type": "string"
          },
          "model": {
            "type": {
              "key": {
                "type": "string",
                "enum": [
                  "set",
                  [
                    "clustered",
                    "standalone"
                  ]
                ]
              }
            }
          },
          "connected": {
            "type": "boolean"
          },
          "leader": {
            "type": "boolean"
          },
          "schema": {
            "type": {
              "min": 0,
              "key": "string"
            }
          },
          "sid": {
            "type": {
              "min": 0,
              "key": "uuid"
            }
          },
          "cid": {
            "type": {
              "min": 0,
              "key": "uuid"
            }
          },
          "index": {
            "type": {
              "min": 0,
              "key": "integer#"
            }
          }
        }
      }
    }
  },
  "error": null
}

```

The `ovs-vsctl` sends a [monitor
cond](https://docs.openvswitch.org/en/latest/ref/ovsdb-server.7/#monitor-cond)
request for several columns on table `Open_vSwitch`:

```
{
  "id": 1,
  "method": "monitor_cond",
  "params": [
    "_Server",
    [
      "monid",
      "_Server"
    ],
    {
      "Database": [
        {
          "where": [
            [
              "name",
              "==",
              "Open_vSwitch"
            ]
          ],
          "columns": [
            "cid",
            "connected",
            "index",
            "leader",
            "model",
            "name",
            "schema",
            "sid"
          ]
        }
      ]
    }
  ]
}

```

The `ovsdb-sever` replies with a full dump of the database schema (in
the `result.Database.<uuid>.initial.schema` field):


```
{
  "id": 1,
  "result": {
    "Database": {
      "ef19f98e-321f-4178-9d8a-5c4482e92694": {
        "initial": {
          "schema": {
            "cksum": "3682332033 23608",
            "name": "Open_vSwitch",
            "tables": {
              "Bridge": {
                # ..
              },
              "Interface": {
                # ..
              },
              "Open_vSwitch": {
                # ..
              },
              "Port": {
                # ..
              },
              # ..
            },
            "version": "7.15.1"
          },
          "leader": true,
          "connected": true,
          "model": "standalone",
          "name": "Open_vSwitch"
        }
      }
    }
  },
  "error": null
}

```

The `ovs-vsctl` client sends a monitor request for several tables
(`Port`, `Interface`, `Controller` and `Bridge`) in the `Open_vSwitch`
database:

```
{
  "id": 2,
  "method": "monitor_cond",
  "params": [
    "Open_vSwitch",
    [
      "monid",
      "Open_vSwitch"
    ],
    {
      "Port": [
        {
          "columns": [
            "fake_bridge",
            "interfaces",
            "name",
            "tag"
          ]
        }
      ],
      "Interface": [
        {
          "columns": [
            "error",
            "name",
            "ofport"
          ]
        }
      ],
      "Controller": [
        {
          "columns": []
        }
      ],
      "Bridge": [
        {
          "columns": [
            "controller",
            "fail_mode",
            "name",
            "ports"
          ]
        }
      ],
      "Open_vSwitch": [
        {
          "columns": [
            "bridges",
            "cur_cfg"
          ]
        }
      ]
    }
  ]
}
```

The `ovs-vsclt` client then calls `set_db_change_aware`, so that the
`ovsdb-server` [does not close the connection on database
changes](https://docs.openvswitch.org/en/latest/ref/ovsdb-server.7/#database-change-awareness):

```
{
  "id": 3,
  "method": "set_db_change_aware",
  "params": [
    true
  ]
}
```

The `ovsdb-server` replies to the `monitor_cond` request with a short message containing a UUID:

```
{
  "id": 2,
  "result": {
    "Open_vSwitch": {
      "731977d5-f606-4bb7-8778-ff2fa2aeb3a9": {
        "initial": {}
      }
    }
  },
  "error": null
}
```

The `ovs-vsctl` client sends a transaction request to create the
bridge, the local port and the local interface. The operations in the
transaction are the following:

1. Wait until there is *no* entry for bridge `<uuid>` in column
   `bridges` of table `Open_vSwitch`. The `<uuid>` value is the value
   returned by the `ovsdb-server` in its response to request 2.
2. Inserts the port object in the `Port` table, refrencing an
   interface that has not been inserted yet in the `Interface` table.
3. Updates the `Open_vSwitch` table by updating the `bridges` field
   with a reference to the yet-to-be-created bridge.
4. Inserts the interface object in the `Interface` table that was
   referenced in the port object in step 2 in this list.
5. Inserts the bridge object in the `Bridge` table, referencing the
   port object created in step 2 in this list.
6. Increments `next_cfg` in the `Open_vSwitch` table.
7. Selects the `next_cfg` field from the `Open_vSwitch` table.


The full transaction request follows:

```
{
  "id": 4,
  "method": "transact",
  "params": [
    "Open_vSwitch",
    {
      "rows": [
        {
          "bridges": [
            "set",
            []
          ]
        }
      ],
      "until": "==",
      "where": [
        [
          "_uuid",
          "==",
          [
            "uuid",
            "731977d5-f606-4bb7-8778-ff2fa2aeb3a9"
          ]
        ]
      ],
      "columns": [
        "bridges"
      ],
      "timeout": 0,
      "op": "wait",
      "table": "Open_vSwitch"
    },
    {
      "uuid-name": "rowf582b328_69bc_4c63_8754_a89a5db8c862",
      "row": {
        "name": "pepe0",
        "interfaces": [
          "named-uuid",
          "rowbf3547d5_0a4a_4751_acea_55c34fe00a5c"
        ]
      },
      "op": "insert",
      "table": "Port"
    },
    {
      "where": [
        [
          "_uuid",
          "==",
          [
            "uuid",
            "731977d5-f606-4bb7-8778-ff2fa2aeb3a9"
          ]
        ]
      ],
      "row": {
        "bridges": [
          "named-uuid",
          "row2ca16eaf_9c74_4b43_8cff_3d91c66c51c9"
        ]
      },
      "op": "update",
      "table": "Open_vSwitch"
    },
    {
      "uuid-name": "rowbf3547d5_0a4a_4751_acea_55c34fe00a5c",
      "row": {
        "name": "pepe0",
        "type": "internal"
      },
      "op": "insert",
      "table": "Interface"
    },
    {
      "uuid-name": "row2ca16eaf_9c74_4b43_8cff_3d91c66c51c9",
      "row": {
        "name": "pepe0",
        "ports": [
          "named-uuid",
          "rowf582b328_69bc_4c63_8754_a89a5db8c862"
        ]
      },
      "op": "insert",
      "table": "Bridge"
    },
    {
      "mutations": [
        [
          "next_cfg",
          "+=",
          1
        ]
      ],
      "where": [
        [
          "_uuid",
          "==",
          [
            "uuid",
            "731977d5-f606-4bb7-8778-ff2fa2aeb3a9"
          ]
        ]
      ],
      "op": "mutate",
      "table": "Open_vSwitch"
    },
    {
      "where": [
        [
          "_uuid",
          "==",
          [
            "uuid",
            "731977d5-f606-4bb7-8778-ff2fa2aeb3a9"
          ]
        ]
      ],
      "columns": [
        "next_cfg"
      ],
      "op": "select",
      "table": "Open_vSwitch"
    },
    {
      "comment": "ovs-vsctl (invoked by strace): ovs-vsctl add-br pepe0",
      "op": "comment"
    }
  ]
}
```

The `ovsdb-server` process replies with several messages to
acknowledge the creation of the bridge and its associated port and
interface.

The first message to arrive is the update for the monitor set up in
request 2, notifying that all relevant objects (interface, port, and
bridge) have been inserted, and that the bridge is known to the
system:

```
{
  "id": null,
  "method": "update2",
  "params": [
    [
      "monid",
      "Open_vSwitch"
    ],
    {
      "Interface": {
        "d1194f67-4c14-4e29-979a-cd0d87ec1448": {
          "insert": {
            "name": "pepe0"
          }
        }
      },
      "Port": {
        "63665e39-8601-4248-8258-9ac33ef822c4": {
          "insert": {
            "name": "pepe0",
            "interfaces": [
              "uuid",
              "d1194f67-4c14-4e29-979a-cd0d87ec1448"
            ]
          }
        }
      },
      "Bridge": {
        "7523cffb-1dcf-4b7c-9746-354c49dc9aa5": {
          "insert": {
            "name": "pepe0",
            "ports": [
              "uuid",
              "63665e39-8601-4248-8258-9ac33ef822c4"
            ]
          }
        }
      },
      "Open_vSwitch": {
        "731977d5-f606-4bb7-8778-ff2fa2aeb3a9": {
          "modify": {
            "bridges": [
              "uuid",
              "7523cffb-1dcf-4b7c-9746-354c49dc9aa5"
            ]
          }
        }
      }
    }
  ]
}
```

The next message to arrive is the response to the transaction request
itself:

```
{
  "id": 4,
  "result": [
    {},
    {
      "uuid": [
        "uuid",
        "63665e39-8601-4248-8258-9ac33ef822c4"
      ]
    },
    {
      "count": 1
    },
    {
      "uuid": [
        "uuid",
        "d1194f67-4c14-4e29-979a-cd0d87ec1448"
      ]
    },
    {
      "uuid": [
        "uuid",
        "7523cffb-1dcf-4b7c-9746-354c49dc9aa5"
      ]
    },
    {
      "count": 1
    },
    {
      "rows": [
        {
          "next_cfg": 1
        }
      ]
    },
    {}
  ],
  "error": null
}
```

Finally, the last message to arrive is a response on the
`monitor_cond` message sent in request 2, saying that:

* The interface object has a new `ofport` value, and
* The `cur_cfg` field in the `Open_vSwitch` table has been
  incremented.

```
{
  "id": null,
  "method": "update2",
  "params": [
    [
      "monid",
      "Open_vSwitch"
    ],
    {
      "Interface": {
        "d1194f67-4c14-4e29-979a-cd0d87ec1448": {
          "modify": {
            "ofport": 65534
          }
        }
      },
      "Open_vSwitch": {
        "731977d5-f606-4bb7-8778-ff2fa2aeb3a9": {
          "modify": {
            "cur_cfg": 1
          }
        }
      }
    }
  ]
}
```


## Creating a new bridge and setting a bridge parameter

Command line:

```
# ovs-vsctl add-br pepe0 -- set Bridge pepe0 fail-mode=secure
```

The first difference with the `ovs-vsctl add-br pepe0` command above
is the transaction request. The step to insert the bridje row in the
`Bridge` table now sets the field `fail_mode` to `secure`:


```
{
  "id": 4,
  "method": "transact",
  "params": [
    "Open_vSwitch",
    # ..
    {
      "uuid-name": "row2c6ced77_fdee_4cd6_a0da_ff68983663b0",
      "row": {
        "ports": [
          "named-uuid",
          "rowda9b4a74_8da5_4947_adc0_8b31fa63a4ec"
        ],
        "name": "pepe0",
        "fail_mode": "secure"
      },
      "op": "insert",
      "table": "Bridge"
    },
    # ..
  ]
}
```


## Deleting a bridge

Command-line:

```
# ovs-vsctl del-br pepe0
```

The initial messages are the same as in the cases above for `ovs-vsctl
add-br`. The transaction request, of course, is different:

1. Wait until there is an entry for `pepe0` in the column `bridges` of
   table `Open_vSwitch`.
2. Update field `bridges` in table `Open_vSwitch`, deleting the entry for `pepe0`.
3. Increase `next_cfg` by one in table `Open_vSwitch`
4. Select `next_cfg` from table `Open_vSwitch`.

Note that there are no actions for deleting the port nor the interface
objects from their respective tables.

```
{
  "id": 4,
  "method": "transact",
  "params": [
    "Open_vSwitch",
    {
      "rows": [
        {
          "bridges": [
            "uuid",
            "7523cffb-1dcf-4b7c-9746-354c49dc9aa5"
          ]
        }
      ],
      "until": "==",
      "where": [
        [
          "_uuid",
          "==",
          [
            "uuid",
            "731977d5-f606-4bb7-8778-ff2fa2aeb3a9"
          ]
        ]
      ],
      "columns": [
        "bridges"
      ],
      "timeout": 0,
      "op": "wait",
      "table": "Open_vSwitch"
    },
    {
      "where": [
        [
          "_uuid",
          "==",
          [
            "uuid",
            "731977d5-f606-4bb7-8778-ff2fa2aeb3a9"
          ]
        ]
      ],
      "row": {
        "bridges": [
          "set",
          []
        ]
      },
      "op": "update",
      "table": "Open_vSwitch"
    },
    {
      "mutations": [
        [
          "next_cfg",
          "+=",
          1
        ]
      ],
      "where": [
        [
          "_uuid",
          "==",
          [
            "uuid",
            "731977d5-f606-4bb7-8778-ff2fa2aeb3a9"
          ]
        ]
      ],
      "op": "mutate",
      "table": "Open_vSwitch"
    },
    {
      "where": [
        [
          "_uuid",
          "==",
          [
            "uuid",
            "731977d5-f606-4bb7-8778-ff2fa2aeb3a9"
          ]
        ]
      ],
      "columns": [
        "next_cfg"
      ],
      "op": "select",
      "table": "Open_vSwitch"
    },
    {
      "comment": "ovs-vsctl (invoked by strace): ovs-vsctl del-br pepe0",
      "op": "comment"
    }
  ]
}
```

The first response to arrive is the notification for the monitor,
ionforming us that the bridge has been deregistered, and its related
entries in the `Bridge`, `Port` and `Interface` tables have been
removed:

```
{
  "id": null,
  "method": "update2",
  "params": [
    [
      "monid",
      "Open_vSwitch"
    ],
    {
      "Interface": {
        "d1194f67-4c14-4e29-979a-cd0d87ec1448": {
          "delete": null
        }
      },
      "Port": {
        "63665e39-8601-4248-8258-9ac33ef822c4": {
          "delete": null
        }
      },
      "Bridge": {
        "7523cffb-1dcf-4b7c-9746-354c49dc9aa5": {
          "delete": null
        }
      },
      "Open_vSwitch": {
        "731977d5-f606-4bb7-8778-ff2fa2aeb3a9": {
          "modify": {
            "bridges": [
              "uuid",
              "7523cffb-1dcf-4b7c-9746-354c49dc9aa5"
            ]
          }
        }
      }
    }
  ]
}
```

Then the response to the transaction request itself is received:

```
{
  "id": 4,
  "result": [
    {},
    {
      "count": 1
    },
    {
      "count": 1
    },
    {
      "rows": [
        {
          "next_cfg": 2
        }
      ]
    },
    {}
  ],
  "error": null
}
```

And finally the notification that the `cur_cfg` field has been increased:

```
{
  "id": null,
  "method": "update2",
  "params": [
    [
      "monid",
      "Open_vSwitch"
    ],
    {
      "Open_vSwitch": {
        "731977d5-f606-4bb7-8778-ff2fa2aeb3a9": {
          "modify": {
            "cur_cfg": 2
          }
        }
      }
    }
  ]
}
```
