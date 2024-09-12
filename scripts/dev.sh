#!/bin/bash

# Function to check for required dependencies
check_deps() {
  echo "Checking for required dependencies..."

  # Check if Go is installed
  if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go."
    exit 1
  fi

  # Check if Node.js is installed
  if ! command -v node &> /dev/null; then
    echo "Error: Node.js is not installed. Please install Node.js."
    exit 1
  fi

  # Check if npm is installed
  if ! command -v npm &> /dev/null; then
    echo "Error: npm is not installed. Please install npm."
    exit 1
  fi

  # Check if pnpm is installed
  if ! command -v pnpm &> /dev/null; then
    echo "Error: pnpm is not installed. Please install pnpm."
    exit 1
  fi

  # Check if Air is installed
  if ! command -v air &> /dev/null; then
    echo "Error: Air is not installed. Please install Air (https://github.com/cosmtrek/air)."
    exit 1
  fi

  echo "All dependencies are installed!"
}

# Function to run the Go server
run_go() {
  echo "Starting Go server..."
  air 2>&1 | sed "s/^/[Go] /"
}

# Function to run the Vite dev server
run_vite() {
  echo "Starting Vite dev server..."
  cd web
  pnpm run dev 2>&1 | sed "s/^/[Vite] /"
}

# Run both commands concurrently and capture their logs
check_deps
run_go &  # Run Go server in background
run_vite  # Run Vite server in foreground

# Wait for both processes to finish
wait
