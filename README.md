# AlgoSeas Hackathon

This project indexes every single AlgoSeas pirate which has been minted & had its stats claimed, and then continuously polls for new assets / updates to the existing metadata via on-chain activity extracted from AlgoIndexer. The software also polls for the set of active listings via the Algoseas Marketplace API.

The application maintains an in-memory KdTree which is a data structure that allows us to add/remove n-dimensional points (in this case [combat, constitution, luck, plunder], and to find the K nearest neighbors (K most similar via euclidean distance)).

The application also exposes a RESTFUL HTTP server on port 8080 with the following routes:

```
GET /similar?assetId=x&amount=y
GET /assets?assetId=x
```

The `similar` route returns the y most similar assets (calculated via euclidean distance), alongside any listings associated with the y closest related assets which have active listings.

**Request:**
`curl http://localhost:8080/similar?assetId=815577765&amount=2`

**Response:**

```
{
    "SimilarAssets": [
        {
            "Id": 695314315,
            "UpdatedAt": 20274435,
            "Collection": "AlgoSeas Pirates",
            "ImageUrl": "https://test.cdn.algoseas.io/pirates/1174-full.png",
            "Combat": 57,
            "Constitution": 53,
            "Luck": 49,
            "Plunder": 32,
            "Scenery": "Majesty",
            "LeftArm": "Empty",
            "Body": "Coconut",
            "BackItem": "Saber",
            "Pants": "Atlas",
            "Footwear": "Sandals",
            "RightArm": "Empty",
            "Shirts": "Atlas",
            "Hat": "Flow",
            "Face": "Serious",
            "BackgroundAccent": "Grass",
            "Necklace": "Bird Skull",
            "FacialHair": "Hook"
        },
        {
            "Id": 700774135,
            "UpdatedAt": 20783949,
            "Collection": "AlgoSeas Pirates",
            "ImageUrl": "https://cdn.algoseas.io/pirates/2147-full.png",
            "Combat": 57,
            "Constitution": 53,
            "Luck": 50,
            "Plunder": 32,
            "Scenery": "Cloud Ocean",
            "Body": "Seafoam",
            "Pants": "Sea",
            "Footwear": "Bandages",
            "Tattoo": "Sailing",
            "Face": "Cigar",
            "Head": "Captain",
            "BackHand": "Darts",
            "Pet": "Salvador"
        }
    ],
    "RelatedListings": [
        {
            "assetId": 695314315,
            "listing": {
                "date": "2022-07-10T15:42:13.000Z",
                "isDutch": false,
                "marketplace": "algoseas",
                "minBidDelta": 1000000000,
                "nextPayout": "BDLXHCWENALYKQUDTOUS6TKVX4K3GBFZYKJK553BQHH3Z34WPZWGOUCA2E",
                "seller": "BDLXHCWENALYKQUDTOUS6TKVX4K3GBFZYKJK553BQHH3Z34WPZWGOUCA2E",
                "price": 200000000,
                "quantity": 1,
                "royalty": 500,
                "royaltyString": "AfQNBQ0FDQZ9DISVYGbHb3flVEzYJHJTRPDfbyEHBpc1A2jcMYo0y6qe3bRKOG64bIJUT7yL83p5hxzJrXGe4eHKpQ2YbEnTDF0juPpkscbaHkWdvr+flxXnpRKgX1j773TCegWCbQ0=",
                "verifiedRoyalty": true,
                "listingID": 803760500,
                "variableID": 803760499
            }
        },
        {
            "assetId": 846289864,
            "listing": {
                "date": "2022-09-04T15:40:20.000Z",
                "isDutch": false,
                "marketplace": "algoseas",
                "minBidDelta": 1000000000,
                "nextPayout": "7QREIXC4JCK3UMFXJZRDWE7UT7A3YQGEYGZRNIZIEBOAWACITLXZYV64TQ",
                "seller": "7QREIXC4JCK3UMFXJZRDWE7UT7A3YQGEYGZRNIZIEBOAWACITLXZYV64TQ",
                "price": 5000000,
                "quantity": 1,
                "royalty": 500,
                "royaltyString": "AfQNBQ0FDQZ9DISVYGbHb3flVEzYJHJTRPDfbyEHBpc1A2jcMYo0y6qe3bRKOG64bIJUT7yL83p5hxzJrXGe4eHKpQ2YbEnTDF0juPpkscbaHkWdvr+flxXnpRKgX1j773TCegWCbQ0=",
                "verifiedRoyalty": true,
                "listingID": 861012908,
                "variableID": 861012907
            }
        }
    ]
}
```

The `assets` route allows the MetaData to be queried for any given asset

**Request:**
`curl http://localhost:8080/assets?assetId=815577765`

**Response:**

```
{
    "Id": 815577765,
    "UpdatedAt": 22373201,
    "Collection": "AlgoSeas Pirates",
    "ImageUrl": "https://cdn.algoseas.io/pirates/14548-full.png",
    "Combat": 52,
    "Constitution": 28,
    "Luck": 37,
    "Plunder": 56,
    "Scenery": "Cloud Ocean",
    "Body": "Seafoam",
    "Pants": "Violet",
    "Footwear": "Boots",
    "HipItem": "Bomb Belt",
    "Face": "Cigar",
    "BackgroundAccent": "Seagulls",
    "Necklace": "Glitter",
    "Head": "Long Hair"
}
```

## System Design

A MySQL database is used to house all assets information (`asset` table). There is a considerable amount of data, this information is critical to the project, and the fetching of historical data can take some time depending on the speed of Algoindexer API, so it makes sense to persist this information.

The software maintains several in-memory hashmaps, the two key mappings are from Asset ID => Listing which allows the quick serving of active listings without the complexity of having to constantly alter a database table - we simple overwrite the map every time the active-listings endpoint is polled for data, and also Asset ID => Asset Object which is used to map between the outputs of similarity calculations, and to serve the `/assets` endpoint.

## Running the project

The code is written in **GoLang**, and the database used is **MySQL** - both of these dependencies **must be installed**.

### Dependencies

1. [GoLang](https://go.dev/dl/)
2. [MySQL](https://www.mysql.com/downloads/)

Once the dependencies are installed, you need to ensure that the MySQL server is running on the machine.

### Configuration

The following configuration steps must be taken prior to running the program:

1. Create a database named `algoseas`: (`CREATE DATABASE algoseas;`)
2. Create a db user: `CREATE USER 'X'@'localhost' IDENTIFIED BY 'Y';`, where X is the username, and Y is the password.
3. Grant permission to this newly created user: `GRANT ALL PRIVILEGES ON algoseas.* TO 'X'@'localhost';`
4. Create a .env file and populate the following fields with the username and password you previously set:

```
DB_USER=X
DB_PASSWORD=Y
```

### Executing the program

Once the configuration steps have been completed, simply run the following command from this directory in your terminal of choice:
`go run .`

Please note that the initial population of the assets could take a while, however once this has been completed the following message will be logged to the console: `Finished initial load`, from this point onwards the API will be exposed and the polling for new data will begin.

## Extensibility, and to the future...

If we wanted to extend this to other collections, we would need to create a new database schema as the current `assset` table is tailored to this collection, so it would probably require a new table & a new set of parsers. Once this has been wired up, you'd also need to add some more logic into the polling/ fetching of all assets so that it fetches the data for the additional collections. For each collection new KdTree would need to be instantiated and populated.

It is also worth noting that this current model can only be used to find similar assets based on numerical properties e.g. `[combat, constitution, luck, plunder]`, as categorical data cannot be used to find a euclidean distance.

A mapping of collection => KdTree, and any other supporting helper functions to go from assetId => collection etc would enable the current model to scale to several different collections.
