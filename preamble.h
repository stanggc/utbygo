#pragma once

#include <stdbool.h>
#include <sys/types.h>

struct ServerOptions {
	char *TLSCertFile;
	char *TLSKeyFile;
};

struct Response {
	char *ID;
	int StatusCode;
	char *Headers;
	size_t ContentLength;
	char * Content;
};

struct Request {
	char *ID;
	char *Method;
	char *URL;
	char *Headers;
	size_t ContentLength;
	char *Content;
};

struct FetchOptions {
	bool SkipVerifyTLSCert;
};

enum Status {
	E_OK = 0,
	E_START_SERVER,
	E_SERVER_STARTED,
	E_SERVER_NOT_STARTED,
	E_NO_MSG,
	E_URL_PARSE,
	E_FETCH,
	E_READ,
};

void FreeRequest(struct Request *req);
void FreeResponse(struct Response *resp);
