#!/bin/bash

# Demo script to run all examples in sequence
# This script demonstrates all the backoff library examples

echo "========================================"
echo "Backoff Library Examples Demo"
echo "========================================"
echo

echo "1. Basic Usage (with jitter):"
echo "   Demonstrates default exponential backoff with randomization"
echo "   Expected: Random delays that grow exponentially"
echo
cd examples/basic && go run main.go
echo

echo "----------------------------------------"
echo

echo "2. Without Jitter:"
echo "   Demonstrates predictable exponential backoff"
echo "   Expected: Exact delays: 1s, 2s, 4s, 8s"
echo
cd ../without-jitter && go run main.go
echo

echo "----------------------------------------"
echo

echo "3. Reset Functionality:"
echo "   Demonstrates how to reset backoff state"
echo "   Expected: Multiple operations starting from initial delay"
echo
cd ../reset && go run main.go
echo

echo "----------------------------------------"
echo

echo "4. HTTP Retry Pattern:"
echo "   Demonstrates real-world HTTP retry with backoff"
echo "   Expected: HTTP requests with exponential delays between retries"
echo
cd ../http-retry && go run main.go
echo

echo "========================================"
echo "Demo completed!"
echo "========================================"
