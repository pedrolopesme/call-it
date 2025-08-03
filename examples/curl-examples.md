# cURL Command Examples for Call-It

This document shows examples of cURL commands that can be pasted into Call-It using **Ctrl+P**.

## Simple GET Request
```bash
curl https://httpbin.org/get
```

## GET with Headers
```bash
curl -H "Authorization: Bearer your-token-here" -H "User-Agent: Call-It/1.0" https://httpbin.org/get
```

## POST with JSON Data
```bash
curl -X POST https://httpbin.org/post -H "Content-Type: application/json" -d '{"name": "John", "age": 30}'
```

## POST with Form Data
```bash
curl -X POST https://httpbin.org/post -H "Content-Type: application/x-www-form-urlencoded" -d "name=John&age=30"
```

## PUT Request
```bash
curl -X PUT https://httpbin.org/put -H "Content-Type: application/json" -d '{"id": 1, "name": "Updated Name"}'
```

## DELETE Request
```bash
curl -X DELETE https://httpbin.org/delete -H "Authorization: Bearer your-token"
```

## Complex Request with Multiple Headers
```bash
curl -X POST https://httpbin.org/post \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" \
  -H "X-API-Key: abc123def456" \
  -H "User-Agent: MyApp/2.0" \
  -d '{"user": {"name": "Alice", "email": "alice@example.com"}, "action": "create"}'
```

## How to Use

1. Launch Call-It: `./call-it`
2. Press **Ctrl+P** to enter curl paste mode
3. Paste any of the above examples
4. Press **Enter** to parse the command
5. The form fields will be automatically populated
6. Press **Enter** again to start the load test

## Features Supported

- ✅ HTTP methods (GET, POST, PUT, DELETE, etc.)
- ✅ Headers (-H, --header)
- ✅ Request body (-d, --data)
- ✅ URL parsing
- ✅ Content-Type detection
- ✅ Authorization headers
- ✅ Custom headers

## Notes

- All cURL commands are parsed and converted to Call-It's internal format
- You can modify the parsed values before running the test
- Complex cURL features like file uploads, cookies, and proxies may not be fully supported
- The parser is designed to handle most common API testing scenarios