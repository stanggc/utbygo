#include <stdlib.h>
#include "preamble.h"

void FreeRequest(struct Request *req) {
	if (!req) return;
	if (req->ID) { free(req->ID); req->ID = NULL; }
	if (req->Method) { free(req->Method); req->Method = NULL; }
	if (req->URL) { free(req->URL); req->URL = NULL; }
	if (req->Headers) { free(req->Headers); req->Headers = NULL; }
	if (req->Content) { free(req->Content); req->Content = NULL; }
	req->ContentLength = 0;
}

void FreeResponse(struct Response *resp) {
	if (!resp) return;
	if (resp->ID) { free(resp->ID); resp->ID = NULL; }
	if (resp->Headers) { free(resp->Headers); resp->Headers = NULL; }
	if (resp->Content) { free(resp->Content); resp->Content = NULL; }
	resp->ContentLength = 0;
}
