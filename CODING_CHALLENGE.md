# Coding Test: Purchase Cart Service

Thank you for taking the time to complete our coding test. Your challenge is to recreate a
simplified version of a purchase cart service.

## Problem
1. You must use Git. Git is mandatory in Subito, so the test must be provided by a
repository (e.g. GitHub).
2. You must use Docker. Docker is a mandatory technology to master. Each project
must contain:
    - a Dockerfile
    - a command to execute the test suite within the Docker container
    - a command to run the service within the Docker container

    We typically suggest using script/tests.sh and scripts/run.sh as wrapper scripts.
3. Subito's main backend language is Golang, but feel free to use your favourite
language.

Your goal is to create a RESTful service that, given a set of products, allows you to create an
order.

The returned data structure must include:
- the ID of the order
- the total price for the order
- the total VAT for the order
- price and VAT for each item in the order

Take some time to think about and understand the assignment. Design thoroughly every
aspect of the solution, such as:
- the endpoint
- proper pricing data
- the storage
- the orders
- potential evolutions

In the README file, describe how to run and test the project, as well as the considerations
taken.

N.B. the main evaluated aspects are:
- the readability of the code
- the structure of the application
- the quality of the tests implemented
- the quality of the documentation