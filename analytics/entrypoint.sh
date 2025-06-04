#!/bin/sh
set -e

# Fallbacks if YEAR and MONTH are unset
YEAR=${YEAR:-$(date +%Y)}
MONTH=${MONTH:-$(date +%-m)}

# Print env for debugging
echo "Running in-out-analytics pipeline with:"
echo "  YEAR: $YEAR"
echo "  MONTH: $MONTH"

# Run main with env vars
python main.py "$YEAR" "$MONTH"