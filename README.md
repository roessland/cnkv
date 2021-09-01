# cnkv
Cloud Native distributed resilient 12-factor instrumented 
scalable hexagonal architecture scalable key-value store.

Not really though, but that's the goal.

## TODO

- [ ] Add tests
- [ ] Close transaction log when program is closed
- [ ] Make sure buffer is empty and log is closed when closing program
- [ ] Encode keys/values in transaction log so they can contain arbitrary stuff
- [ ] Add a bound to key/val size in bytes to prevent disk from filling
- [ ] Use protobuf for transaction log to save space and increase performance
- [ ] Add log compaction (offline) to remove deleted values from the file
- [ ] Add pipeline to create Docker image
- [ ] Add pipeline to create github release with Linux/Win/MacOS binaries
- [ ] Add a sweet homepage

## Run PostgreSQL server for development

    docker-compose up

DB admin accessible at http://localhost:5433/

# Build Docker image

    docker build --tag cnkv .

Run it using

    