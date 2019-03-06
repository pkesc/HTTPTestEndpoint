# HTTPTestEndpoint

You need to implement a special API or want to test how your code acts when it receives a special HTTP reason or HTTP body? Then this tool is what you are searching for.

| URL Path                         | Description                                                  |
| -------------------------------- | ------------------------------------------------------------ |
| /status/[code]                   | Answers with the passed HTTP status code. Optionally can the received HTTP body and content-type header also be in the answer. |
| /delay/[seconds]                 | Waits the passed seconds and answers 200 OK. Optionally can the received HTTP body and content-type header also be in the answer. |
| /ip                              | Returns the remote IP address in the HTTP body. |
| /userAgent                       | Returns the user-agent string from the HTTP header in the HTTP body. |
| /basicAuth/[username]/[password] | Checks if the username and password of the authorization are the same as the one passed in the URL path. If yes 200 OK and if not 401 unauthorized will be answered. |
| /file/[filename]                 | Returns the content of the passed file in the HTTP body. The file must be in the same directory the HTTPTestEndpoint is. |

For more details start the tool and open [http://localhost:8080](http://localhost:8080/) in your browser to see the detailed help page.

![](assets/helpPage.png)