# HLX

### Test application build with Go and Dgraph.

## Installation

Clone/download and run ```docker-compose up``` inside project folder.
This will start the Dgraph cluster and the Go application.


## API

The API server is listening on port ```:8080```

### Endpoints:

- ```/api/categories/load``` (POST): 

    Use this route to load the JSON for the categories passing it into the request body. This will delete all previous data stored on the database and save the new data.
    The structure of the JSON must follow the ```categories.json sample``` file.

- ```/api/categories/leaves``` (POST):

    Use this route to query for leaves (categories without children) of the categories tree. Results include categories matching the ```Name``` value of the request. Sending ```SubLayerName``` value will get results only for that sub-layer of the tree.
    Request must have the following structure:
    ```json
    {
        "name": "string",
        "subLayerName": "string (optional)"
    }
    ```

- ```/api/psychos/load``` (POST): 

    Use this route to load the JSON for the psychographics passing it into the request body. This will delete all previous data stored on the database and save the new data.
    The structure of the JSON must follow the ```psychographics.json sample``` file.

- ```/api/psychos/leaves``` (POST):

    Use this route to query for leaves (psychographics without nested values) of the psychographics tree. Results include psychographics matching the ```Label``` value of the request. Sending ```SubLayerLabel``` value will get results only for that sub-layer of the tree.
    Request must have the following structure:
    ```json
    {
        "label": "string",
        "subLayerLabel": "string (optional)"
    }
    ```

## Notes

- Load data first using one of the /load endpoints.
- Loading data will delete all previous data from the database.