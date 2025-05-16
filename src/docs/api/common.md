# Common conventions followed by all the APIs

- All API endpoints are prefixed with version ie `/api/v1`.

- Data Response should be wrapped in a `data` object.

  ```json
  {
    "data": {
      "id": "string",
      "name": "string",
      "createdAt": "string",
      "updatedAt": "string"
    }
  }
  ```

- Message Response (acknowledgement) should be wrapped in a `message` field.

  ```json
  {
    "message": "string"
  }
  ```

- Error Response should be wrapped in an `message` field.

  ```json
  {
    "message": "string"
  }
  ```
