#!/bin/bash

# Number of total requests
TOTAL_REQUESTS=160

# Number of concurrent requests
CONCURRENCY=110

# Function to make a single API call
make_request() {
    curl --silent --location 'http://localhost:8080/calculate' \
    --header 'Content-Type: application/json' \
    --data '{
        "revenue": 101,
        "expenses": 1,
        "taxRate": 0.3
    }'
    echo ""  # Add newline after each response
}

# Function to run concurrent requests
run_concurrent_requests() {
    for ((i=1; i<=TOTAL_REQUESTS; i++))
    do
        # Run CONCURRENCY number of requests in background
        ((j=i%CONCURRENCY))
        if ((j==0)); then
            # Every CONCURRENCY requests, wait for all to complete
            make_request
            wait
            echo "Completed $i requests"
        else
            make_request &
        fi
    done
    # Wait for any remaining background requests
    wait
}

echo "Starting $TOTAL_REQUESTS requests with concurrency of $CONCURRENCY"
time run_concurrent_requests
echo "All requests completed"
