# Explorer451

This is a simple file explorer that allows you to navigate through the files and directories of S3.

## Using the API

List root level

```shell
curl http://localhost:8080/api/buckets
{
   "buckets":[
      "nb-bucket-eu-central-1"
   ]
}
```

List contents of 'documents' folder:

```shell
curl http://localhost:8080/api/buckets/nb-bucket-eu-central-1/objects
{
   "items":[
      {
         "key":"folder1",
         "size":0,
         "isFolder":true,
         "type":"folder"
      },
      {
         "key":"folder2",
         "size":0,
         "isFolder":true,
         "type":"folder"
      }
   ],
   "isTruncated":false,
   "totalItems":2,
   "pageSize":100
}
```

List contents of nested folder

```shell
curl http://localhost:8080/api/buckets/nb-bucket-eu-central-1/objects?prefix=folder1/
{
   "items":[
      {
         "key":"folder1/file1.txt",
         "size":0,
         "lastModified":"2024-11-12T23:48:36Z",
         "isFolder":false,
         "type":"file",
         "contentType":"text/plain"
      }
   ],
   "isTruncated":false,
   "totalItems":2,
   "pageSize":100
}
```
