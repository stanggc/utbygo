//go:build ignore

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "utbygo.h"

int main(int argc, char *argv[]) {
	int port = 0;
	char *err = NULL;
	struct ServerOptions serverOpts = { 0 };
	struct Request req = { 0 };
	enum Status r = E_OK;

	port = 8080;
	serverOpts.TLSCertFile = "cert.crt";
	serverOpts.TLSKeyFile = "cert.key";

	if (port <= 0) {
		printf("Specify a valid number as port number.\n");
		goto finish;
	}

	if ((r = StartServer(port, serverOpts))) {
		printf("Unable to start server at port %d.\n", port);
		goto finish;
	}

	err = ServerError();
	if (err) {
		printf("Error starting server: %s\n", err);
		r = E_START_SERVER;
		goto finish;
	}

	if (serverOpts.TLSCertFile && serverOpts.TLSKeyFile) {
		printf("Listening securely on port %d.\n", port);
	} else {
		printf("Listening on port %d.\n", port);
	}

	do {
		if ((r = ReadRequest(&req))) {
			// Reply error.
			if (r == E_NO_MSG) {
				Sleep(1);
				goto nextIter;
			} else {
				printf("Error occurred while read request. Error code: %d\n", r);
			}
		} else {
			// Has message.
			if (!strcmp(req.ID, "SHUTDOWN")) {
				// Server has shut down. Can stop loop.
				printf("Server has shut down.\n");
				FreeRequest(&req);
				break;
			} else if (!strcmp(req.URL, "/shutdown")) {
				SendResponse((struct Response){
					.ID = req.ID,
					.StatusCode = 200,
					.Headers =
						"Content-Type: text/plain\r\n",
					.Content = "OK shutting down.",
					.ContentLength = strlen("OK shutting down."),
				});
				StopServer();
			} else {
				SendResponse((struct Response){
					.ID = req.ID,
					.StatusCode = 200,
					.Headers =
						"Content-Type: text/plain\r\n",
					.Content = "Hello!",
					.ContentLength = strlen("Hello!"),
				});
			}
		}

	nextIter:
		FreeRequest(&req);
	} while (1);

finish:
	if (err) free(err);
	if (!r) printf("OK: test_server.\n");
	return r;
}
