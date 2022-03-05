start_time=$(date +%s.%3N)
# perform a task
lessc tests/test.less
end_time=$(date +%s.%3N)

# elapsed time with millisecond resolution
# keep three digits after floating point.
elapsed=$(echo "scale=3; $end_time - $start_time" | bc)
echo $elapsed
