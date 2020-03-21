# Patches

## Environment Variables
* `PATCHES_DB_USERNAME`: username for accessing DB
* `PATCHES_DB_PASSWORD`: password for accessing DB
* `PATCHES_DB_HOST`: host where DB is located
* `PATCHES_DB_PORT`: port where DB is located
* `PATCHES_DB_DATABASE`: name of DB

## APIS

### `GET /patches/v1/patches/`
Gets patches from the database with filtering.
#### Query string parameters
```
convo_id      Required      int
user_id       Optional      int
type          Optional      ENUM(formatting)
end_time      Optional      Time (2020-01-02T00:00:05Z)
start_time    Optional      Time (2020-01-02T00:00:05Z)
```
#### Response format
`200 OK`
```
{
    "patches": [
        {
            "timestamp": "2019-10-01 20:00:00",
            "patch": "",
            "convo_id": 1,
            "user_id": 1,
            "type": ""
        },
        ...
    ]
}
```
