# Ingestion

This piece of software is responsible for fetching every Algoseas NFT which has been minted, and had the stats claimed.
Additionally, whilst the software is being executed, it will continuously poll for new mints & updates to metadata via contract events extracted via AlgoIndexer, and also poll for new listings of the pirates via the Algoseas marketplace API.

Additionally, the server exposes the API.... TODO HERE!!

## Design

A MySQL database is used to house all assets information (`asset` table). There is a considerable amount of data, this information is critical to the project, and the fetching of historical data can take some time depending on the speed of Algoindexer API, so it makes sense to persist this information.

The software maintains an in-memory mapping of asset id => active listings, this allows the quick serving of active listings without the complexity of having to constantly alter a database table - we simple overwrite the map every time the active-listings endpoint is polled for data.

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
