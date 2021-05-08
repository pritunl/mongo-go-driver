<p align="center"><img src="etc/assets/mongo-gopher.png" width="250"></p>
<p align="center">
  <a href="https://goreportcard.com/report/github.com/pritunl/mongo-go-driver"><img src="https://goreportcard.com/badge/github.com/pritunl/mongo-go-driver"></a>
  <a href="https://pkg.go.dev/github.com/pritunl/mongo-go-driver/mongo"><img src="etc/assets/godev-mongo-blue.svg" alt="docs"></a>
  <a href="https://pkg.go.dev/github.com/pritunl/mongo-go-driver/bson"><img src="etc/assets/godev-bson-blue.svg" alt="docs"></a>
  <a href="https://docs.mongodb.com/drivers/go/"><img src="etc/assets/docs-mongodb-green.svg"></a>
</p>

# MongoDB Go Driver

The MongoDB supported driver for Go.

-------------------------
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Bugs / Feature Reporting](#bugs--feature-reporting)
- [Testing / Development](#testing--development)
- [Continuous Integration](#continuous-integration)
- [License](#license)

-------------------------
## Requirements

- Go 1.10 or higher. We aim to support the latest supported versions of go.
- MongoDB 2.6 and higher.

-------------------------
## Installation

The recommended way to get started using the MongoDB Go driver is by using go modules to install the dependency in
your project. This can be done either by importing packages from `github.com/pritunl/mongo-go-driver` and having the build
step install the dependency or by explicitly running

```bash
go get github.com/pritunl/mongo-go-driver/mongo
```

When using a version of Go that does not support modules, the driver can be installed using `dep` by running

```bash
dep ensure -add "github.com/pritunl/mongo-go-driver/mongo"
```

-------------------------
## Usage

To get started with the driver, import the `mongo` package and create a `mongo.Client` with the `Connect` function:

```go
import (
    "github.com/pritunl/mongo-go-driver/mongo"
    "github.com/pritunl/mongo-go-driver/mongo/options"
    "github.com/pritunl/mongo-go-driver/mongo/readpref"
)

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
```

Make sure to defer a call to `Disconnect` after instantiating your client:

```go
defer func() {
    if err = client.Disconnect(ctx); err != nil {
        panic(err)
    }
}()
```

For more advanced configuration and authentication, see the [documentation for mongo.Connect](https://pkg.go.dev/github.com/pritunl/mongo-go-driver/mongo#Connect).

Calling `Connect` does not block for server discovery. If you wish to know if a MongoDB server has been found and connected to,
use the `Ping` method:

```go
ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
err = client.Ping(ctx, readpref.Primary())
```

To insert a document into a collection, first retrieve a `Database` and then `Collection` instance from the `Client`:

```go
collection := client.Database("testing").Collection("numbers")
```

The `Collection` instance can then be used to insert documents:

```go
ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
id := res.InsertedID
```

To use `bson.D`, you will need to add `"github.com/pritunl/mongo-go-driver/bson"` to your imports.

Your import statement should now look like this:

```go
import (
    "github.com/pritunl/mongo-go-driver/bson"
    "github.com/pritunl/mongo-go-driver/mongo"
    "github.com/pritunl/mongo-go-driver/mongo/options"
    "github.com/pritunl/mongo-go-driver/mongo/readpref"
)
```

Several query methods return a cursor, which can be used like this:

```go
ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
cur, err := collection.Find(ctx, bson.D{})
if err != nil { log.Fatal(err) }
defer cur.Close(ctx)
for cur.Next(ctx) {
   var result bson.D
   err := cur.Decode(&result)
   if err != nil { log.Fatal(err) }
   // do something with result....
}
if err := cur.Err(); err != nil {
  log.Fatal(err)
}
```

For methods that return a single item, a `SingleResult` instance is returned:

```go
var result struct {
    Value float64
}
filter := bson.D{{"name", "pi"}}
ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
err = collection.FindOne(ctx, filter).Decode(&result)
if err != nil {
    log.Fatal(err)
}
// Do something with result...
```

Additional examples and documentation can be found under the examples directory and [on the MongoDB Documentation website](https://docs.mongodb.com/drivers/go/).

-------------------------
## Feedback

For help with the driver, please post in the [MongoDB Community Forums](https://developer.mongodb.com/community/forums/tag/golang/).

New features and bugs can be reported on jira: https://jira.mongodb.org/browse/GODRIVER

-------------------------
## Testing / Development

The driver tests can be run against several database configurations. The most simple configuration is a standalone mongod with no auth, no ssl, and no compression. To run these basic driver tests, make sure a standalone MongoDB server instance is running at localhost:27017. To run the tests, you can run `make` (on Windows, run `nmake`). This will run coverage, run go-lint, run go-vet, and build the examples.

### Testing Different Topologies

To test a **replica set** or **sharded cluster**, set `MONGODB_URI="<connection-string>"` for the `make` command.
For example, for a local replica set named `rs1` comprised of three nodes on ports 27017, 27018, and 27019:

```
MONGODB_URI="mongodb://localhost:27017,localhost:27018,localhost:27018/?replicaSet=rs1" make
```

### Testing Auth and TLS

To test authentication and TLS, first set up a MongoDB cluster with auth and TLS configured. Testing authentication requires a user with the `root` role on the `admin` database. The Go Driver repository comes with example certificates in the `data/certificates` directory. These certs can be used for testing. Here is an example command that would run a mongod with TLS correctly configured for tests:

```
mongod \
--auth \
--tlsMode requireTLS \
--tlsCertificateKeyFile $(pwd)/data/certificates/server.pem \
--tlsCAFile $(pwd)/data/certificates/ca.pem \
--tlsAllowInvalidCertificates
```

To run the tests with `make`, set `MONGO_GO_DRIVER_CA_FILE` to the location of the CA file used by the database, set `MONGODB_URI` to the connection string of the server, set `AUTH=auth`, and set `SSL=ssl`. For example:

```
AUTH=auth SSL=ssl MONGO_GO_DRIVER_CA_FILE=$(pwd)/data/certificates/ca.pem  MONGODB_URI="mongodb://user:password@localhost:27017/?authSource=admin" make
```

Notes:
- The `--tlsAllowInvalidCertificates` flag is required on the server for the test suite to work correctly.
- The test suite requires the auth database to be set with `?authSource=admin`, not `/admin`.

### Testing Compression

The MongoDB Go Driver supports wire protocol compression using Snappy, zLib, or zstd. To run tests with wire protocol compression, set `MONGO_GO_DRIVER_COMPRESSOR` to `snappy`, `zlib`, or `zstd`.  For example:

```
MONGO_GO_DRIVER_COMPRESSOR=snappy make
```

Ensure the [`--networkMessageCompressors` flag](https://docs.mongodb.com/manual/reference/program/mongod/#cmdoption-mongod-networkmessagecompressors) on mongod or mongos includes `zlib` if testing zLib compression.

-------------------------
## Contribution

Check out the [project page](https://jira.mongodb.org/browse/GODRIVER) for tickets that need completing. See our [contribution guidelines](CONTRIBUTING.md) for details.

-------------------------
## Continuous Integration

Commits to master are run automatically on [evergreen](https://evergreen.mongodb.com/waterfall/mongo-go-driver).

-------------------------
## Thanks and Acknowledgement 

<a href="https://github.com/ashleymcnamara">@ashleymcnamara</a> - Mongo Gopher Artwork

-------------------------
## License

The MongoDB Go Driver is licensed under the [Apache License](LICENSE).
