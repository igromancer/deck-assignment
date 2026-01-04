# deck-assignment

Hello, Deck hiring team. 

Please find my implementation of Scrape Job API in this repository.

Entire application can be run by `docker compose up`.

Postman collection is included for your testing convenience as well as an openapi specification. Latter can be accessed on `localhost:8081`.

Suggested testing flow would be the following:

- `POST /apikey` with empty body to get an API key that should be included in `X-Api-Key` header for all subsequent requests.
- `POST /jobs` to create a scrape job. 
- `GET /jobs/:id` to get the job status
- `GET /jobs:id/result` to fetch a job with the result (when status is completed)
- `GET /jobs` to get a list of submitted jobs for the api key. Optional query parameters are `offset=:num` and `limit=:num`. If not provided those are defaulted to 0 and 25.

For local development there is an `adminer` service included so you can see what is in database. Select `Postgres` and login with credentials `user` and `password`.

The version of RabbitMQ also comes with admin interface which can be accessed at `localhost:15672`. Credentials are same as the DB.

I tried to stick with all good architecture practices and tried to structure the codebase according to community guidelines: `https://github.com/golang-standards/project-layout`. In some places I just used my best judgement.

Hope this is sufficient, looking forward to chat.
