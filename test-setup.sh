# Copyright 2025 William Marchesi

# Author: William Marchesi
# Email: will@marchesi.io
# Website: https://marchesi.io/

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/bash
# Test script to verify spool-scanner setup

echo "🔍 Spool Scanner Setup Test"
echo "=========================="

# Check environment variables
echo -e "\n📋 Checking environment variables..."
required_vars=("SPOOLMAN_URL" "PRINTER_1_NAME" "PRINTER_1_URL" "PRINTER_1_KEY")
missing_vars=()

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        missing_vars+=("$var")
        echo "❌ $var is not set"
    else
        if [[ "$var" == *"KEY"* ]]; then
            echo "✅ $var is set (hidden)"
        else
            echo "✅ $var = ${!var}"
        fi
    fi
done

if [ ${#missing_vars[@]} -gt 0 ]; then
    echo -e "\n⚠️  Missing required environment variables!"
    echo "Please set: ${missing_vars[*]}"
    exit 1
fi

# Test Spoolman connection
echo -e "\n🔗 Testing Spoolman connection..."
if curl -s -f "${SPOOLMAN_URL}/api/v1/info" > /dev/null 2>&1; then
    echo "✅ Spoolman is accessible at ${SPOOLMAN_URL}"
    
    # Get spool count
    spool_count=$(curl -s "${SPOOLMAN_URL}/api/v1/spool" | grep -o '"id"' | wc -l)
    echo "📦 Found ${spool_count} spools in Spoolman"
else
    echo "❌ Cannot connect to Spoolman at ${SPOOLMAN_URL}"
    echo "   Make sure Spoolman is running and accessible"
fi

# Test OctoPrint connections
echo -e "\n🖨️  Testing printer connections..."
for i in {1..10}; do
    name_var="PRINTER_${i}_NAME"
    url_var="PRINTER_${i}_URL"
    key_var="PRINTER_${i}_KEY"
    
    if [ -n "${!name_var}" ]; then
        echo -e "\nPrinter ${i}: ${!name_var}"
        
        # Test OctoPrint API
        response=$(curl -s -w "\n%{http_code}" -H "X-Api-Key: ${!key_var}" "${!url_var}/api/version")
        http_code=$(echo "$response" | tail -n1)
        
        if [ "$http_code" = "200" ]; then
            echo "✅ OctoPrint API is accessible"
            
            # Check for OctoPrint-Spoolman-API plugin
            plugin_response=$(curl -s -H "X-Api-Key: ${!key_var}" "${!url_var}/api/plugin/spoolman_api" -d '{"command":"get_current_spool","tool":0}' -H "Content-Type: application/json")
            if [[ "$plugin_response" == *"success"* ]] || [[ "$plugin_response" == *"spool_id"* ]]; then
                echo "✅ OctoPrint-Spoolman-API plugin is installed"
            else
                echo "❌ OctoPrint-Spoolman-API plugin is NOT installed!"
                echo "   Install it manually: ~/oprint/bin/pip install https://github.com/wmarchesi123/octoprint-spoolman-api/archive/main.zip"
                echo "   This plugin is REQUIRED for spool-scanner to work"
            fi
        elif [ "$http_code" = "403" ]; then
            echo "❌ API key is invalid or has insufficient permissions"
        else
            echo "❌ Cannot connect to OctoPrint at ${!url_var} (HTTP $http_code)"
        fi
    fi
done

# Test local spool-scanner if running
echo -e "\n🌐 Testing spool-scanner..."
port="${PORT:-8080}"
if curl -s -f "http://localhost:${port}/" > /dev/null 2>&1; then
    echo "✅ Spool-scanner is running on port ${port}"
else
    echo "ℹ️  Spool-scanner is not running locally on port ${port}"
    echo "   Run 'go run cmd/server/main.go' to start it"
fi

echo -e "\n✨ Setup test complete!"