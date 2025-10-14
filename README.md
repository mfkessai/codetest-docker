# codetest


## Task

Implement the application such that the tests pass.

It should be implemented so that it runs when Docker Compose is executed.

Tests are implemented in Go and can be run with this command:

```
go test
```

You are free to use any method you want, as long as you do not edit the test code (main_test.go). Please discuss the submission deadline with your recruiter.

We have put a main.go file in place as a sample, but you may use any language you want.

## Project outline

Create a service that allows the registration of "transactions", consisting of an amount of money and a product description each.

There is a per-user limit of 1000 for the total transaction amount that may be registered. If registering a specific transaction would surpass the limit for that user, a certain response status (HTTP 402: payment required) should be returned, resulting in an error.

A database scheme has been placed in the db directory, under the assumption that MySQL will be used as the RDBMS. 

## Development

You can start a dummy application container and a MySQL DB using the schema definition by running `docker compose up`.

## Evaluation

90% of the evaluation depends on whether the tests pass.

Other than that, you may implement other improvements if you wish, but as these will not affect the evaluation, we recommend focusing on an implementation that passes the tests first.

## Submission method

Please compress the repository, including your implemented solution, into a ZIP file and submit it via the upload form.
