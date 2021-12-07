# mutual tls example with golang

## useful artical
https://youngkin.github.io/post/gohttpsclientserver/

## caution
Starting in Go 1.15 certificates must contain a SAN entry or the https request will fail. Certificates with only a CN will not be accepted. If Go 1.15 or higher is used, and --common-name is used to generate the CSR, you will likely see the following error from the client:

Get "https://localhost": x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0
As noted in the error message, this problem can be overcome by prefixing the client command with GODEBUG=x509ignoreCN=0.
