//go:build ignore

#include <stdio.h>
#include "utbygo.h"

int main() {
	enum Status r = E_OK;
	struct FetchOptions fetchOpts = {
		.SkipVerifyTLSCert = true,
	};
	struct Response resp = { 0 };

	// Do fetch.
	r = Fetch(
		"GET",
		"https://localhost:8080/hello",
		NULL,
		NULL,
		0,
		fetchOpts,
		&resp
	);
	if (r) {
		printf("Error doing fetch. Error code: %d\n", r);
		FreeResponse(&resp);
		goto finish;
	}

	// Signal server to shut down.
	r = Fetch(
		"GET",
		"https://localhost:8080/shutdown",
		NULL,
		NULL,
		0,
		fetchOpts,
		&resp
	);
	if (r) {
		printf("Error signalling server to shut down. Error code: %d\n", r);
		FreeResponse(&resp);
		goto finish;
	}

finish:
	if (!r) printf("OK: test_client\n");
	return r;
}
