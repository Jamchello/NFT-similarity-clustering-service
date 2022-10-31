# AlgoSeas Hackathon

This project indexes every single AlgoSeas pirate which has been minted & had its stats claimed, and then continuously polls for new assets / updates to the existing metadata via on-chain activity extracted from AlgoIndexer. The software also polls for the set of active listings via the Algoseas Marketplace API.

Upon initial startup, and any successful polling, a KNN clustering algorithm will be used to produce a set of different clusters - each cluster containing a set of Assets which have been deemed similar based upon their stats.

The application also exposes a RESTFUL HTTP server on port 8080 with the following routes:

```
GET /similar?assetId=x
GET /assets?assetId=x
```

The `similar` route returns the entire set of assets which were clustered with the provided assetId, alongside all active marketplace listings which belong to assets in this cluster.

**Request:**
`curl http://localhost:8080/similar?assetId=815577765`

**Response:**

```
TODO!!
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

TODO: Finsh the explanation here, update below...
The software maintains an in-memory mapping of asset id => active listings, this allows the quick serving of active listings without the complexity of having to constantly alter a database table - we simple overwrite the map every time the active-listings endpoint is polled for data.

## Clustering Algorithm

TODO: Fill this in with detailed explanation as to how the model works, and why we chose how many clusters + how many iterations to train the model with.

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

### Methodology

We are using K-Mean clustering as it is one of the most commonly used algorithms for grouping unlabeled datasets into clusters, this allows us to group 'similar' assets. For K-Mean clustering, there are two parameters that we need to decide, the number of iterations, and number of clusters (K). Using a very high number of iterations is normally unnecessary since [K-Means converges after 20-50 iterations](https://static.googleusercontent.com/media/research.google.com/vi//pubs/archive/42853.pdf). With that said, computationally it is not expensive or time consuming to run a very large number of iterations.
