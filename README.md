# Gaston Sanchez - Coding Challenge

The application implementation uses the following libraries:
* gofiber/fiber as a web framework
* jackc/pgx as PostgreSQL driver
* hashicorp/go-memdb as in-memory database to store the common used values


Application endpoints:
* To get a block (get): /api/v1/block/network_code/block_hash
* To get a transaction (get): /api/v1/tx/network_code/transaction_id
* To reload in-memory database (post): /api/v1/admin/memdb/reload


Missing things:
* Unit testing artifacts
* Modular testing using Postman or Katalon
* In-memory caching for blocks and transactions calls using allegro/bigcache
