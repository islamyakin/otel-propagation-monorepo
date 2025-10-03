#!/bin/bash

# Example API Usage Script
# Make sure the server is running on localhost:8080

BASE_URL="http://localhost:8080/api/v1"

echo "=== Todo App API Examples ==="
echo

# Health check
echo "1. Health Check:"
curl -s -X GET http://localhost:8080/health | jq .
echo

# Register a new user
echo "2. Register User:"
curl -s -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "password": "password123"}' | jq .
echo

# Register an admin user
echo "3. Register Admin User:"
curl -s -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}' | jq .
echo

# Login to get token
echo "4. Login User:"
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "password": "password123"}')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
echo $LOGIN_RESPONSE | jq .
echo

# Get user profile
echo "5. Get Profile:"
curl -s -X GET $BASE_URL/profile \
  -H "Authorization: Bearer $TOKEN" | jq .
echo

# Create a todo
echo "6. Create Todo:"
curl -s -X POST $BASE_URL/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Buy groceries", "description": "Milk, bread, eggs"}' | jq .
echo

# Create another todo
echo "7. Create Another Todo:"
curl -s -X POST $BASE_URL/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Complete project", "description": "Finish the todo app"}' | jq .
echo

# Get all user todos
echo "8. Get User Todos:"
curl -s -X GET $BASE_URL/todos \
  -H "Authorization: Bearer $TOKEN" | jq .
echo

# Update todo status
echo "9. Update Todo Status:"
curl -s -X PATCH $BASE_URL/todos/1/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"status": "completed"}' | jq .
echo

# Get specific todo
echo "10. Get Specific Todo:"
curl -s -X GET $BASE_URL/todos/1 \
  -H "Authorization: Bearer $TOKEN" | jq .
echo

# Update todo
echo "11. Update Todo:"
curl -s -X PUT $BASE_URL/todos/2 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Complete project ASAP", "description": "Finish the todo app by today"}' | jq .
echo

echo "=== Admin Examples (need to manually update role in DB) ==="
echo "To make 'admin' user an admin, run in PostgreSQL:"
echo "UPDATE users SET role = 'admin' WHERE username = 'admin';"
echo

# Login as admin (after manually setting role)
echo "12. Login Admin:"
ADMIN_LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}')

ADMIN_TOKEN=$(echo $ADMIN_LOGIN_RESPONSE | jq -r '.token')
echo $ADMIN_LOGIN_RESPONSE | jq .
echo

# Get all users (admin only)
echo "13. Get All Users (Admin):"
curl -s -X GET $BASE_URL/users \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
echo

# Get all todos (admin only)
echo "14. Get All Todos (Admin):"
curl -s -X GET $BASE_URL/admin/todos \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
echo

echo "=== Error Examples ==="

# Try to access without token
echo "15. Access without token (should fail):"
curl -s -X GET $BASE_URL/todos | jq .
echo

# Try admin endpoint as regular user
echo "16. Admin endpoint as regular user (should fail):"
curl -s -X GET $BASE_URL/users \
  -H "Authorization: Bearer $TOKEN" | jq .
echo

echo "=== Script Complete ==="