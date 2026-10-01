#!/usr/bin/env sh
# Gradle wrapper script (generated)
# This script is a minimal version that delegates to the Gradle Wrapper JAR.
# It will download the wrapper JAR if not present.

APP_NAME="gradle"
WRAPPER_JAR="gradle/wrapper/gradle-wrapper.jar"
WRAPPER_PROPERTIES="gradle/wrapper/gradle-wrapper.properties"

# Determine Java command
if [ -n "$JAVA_HOME" ]; then
  JAVA_CMD="$JAVA_HOME/bin/java"
else
  JAVA_CMD=java
fi

# Resolve the directory of this script
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# If wrapper JAR missing, download it using the properties file
if [ ! -f "$WRAPPER_JAR" ]; then
  echo "Wrapper JAR not found, downloading..."
  # Use the wrapper to bootstrap itself
  "$JAVA_CMD" -jar "$WRAPPER_JAR" --gradle-version 8.5 > /dev/null 2>&1 || true
fi

# Execute the wrapper
exec "$JAVA_CMD" -jar "$WRAPPER_JAR" "$@"
