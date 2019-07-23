curl -v -d @scripts/curl_tests/create_user.json -X POST http://localhost:3000/v1/users | python -m json.tool
