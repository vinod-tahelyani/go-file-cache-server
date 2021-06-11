# **File Cache Manager**

A simpel file cache server which will cache remote files locally and server it through HTTP end point.

## **APIs**

##### `/, /healthy -GET`
    Check the health of the service. 200, if OK.
##### `/cache-file -POST`
    Caches files locally, requires JSON payload. Example:
    {
        "fileURL" "https://www.google.com/robots.txt":
        "username": "foo",
        "password": "bar"
    }
    Response:
    {
        id: "sd2333cdde33344dddd",
        "status" "DOWNLOADED", 
        "message", "Downloaded Successfully"
    }
    
##### `/get-file -GET`
    Get the file through http end point. Requires 'fileURL' or 'ID' as query params.
    Example- http://localhost:8080/get-file?ID=sd2333cdde33344dddd

##### `/invalidate-cache`
    Deletes all the files and meta-data.
    