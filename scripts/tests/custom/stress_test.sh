#!/bin/bash

# Default number of child processes to spawn
N=100000
if [ ! -z "$1" ]; then
    N=$1
fi

# Array to store child PIDs
CHILD_PIDS=()

# Setup signal handler to kill all child processes when parent receives SIGINT or SIGTERM
cleanup() {
    echo "Received termination signal. Killing all child processes..."
    for pid in "${CHILD_PIDS[@]}"; do
        if kill -0 $pid 2>/dev/null; then
            echo "Sending SIGTERM to child process $pid"
            kill -TERM $pid
        fi
    done
    
    # Wait for all children to terminate
    wait
    echo "All child processes terminated. Exiting parent."
    exit 0
}

# Register signal handlers
trap cleanup SIGINT SIGTERM

# Make sure asker.sh is executable
chmod +x ./asker.sh

# Spawn N child processes
echo "Starting $N child processes..."
for ((i=1; i<=N; i++)); do
    ./asker.sh &
    CHILD_PID=$!
    CHILD_PIDS+=($CHILD_PID)
    echo "Started child process $i with PID: $CHILD_PID"
    
    # Optional: slow down the spawning of processes to avoid overwhelming the system
    if (( i % 100 == 0 )); then
        sleep 1
    fi
done

echo "All $N child processes started. Parent PID: $$"
echo "Press Ctrl+C to terminate all processes."

# Wait for all child processes
wait
