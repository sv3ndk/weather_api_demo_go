# Creation of the certificate to be used by the REST client + corresponding trustore to be registered on the API Gateway

echo "creating fake root CA key and certificate"
openssl req \
    -new \
    -newkey ec \
    -pkeyopt ec_paramgen_curve:secp384r1 \
    -sha384 \
    -x509 \
    -nodes \
    -days 365 \
    -out rootCACert.pem \
    -keyout rootCAKey.pem \
    -subj "/C=BE/ST=Test state/L=Brusels/O=Test organization/OU=Test unit/CN=Test root CA"

echo "creating REST client private key"
openssl genpkey \
    -algorithm EC \
    -pkeyopt ec_paramgen_curve:secp384r1 \
    -out clientKey.pem

echo "creating Certificate Signing Request for the client"
openssl req \
    -new \
    -key clientKey.pem \
    -sha384 \
    -out client.csr \
    -subj "/C=BE/ST=Client state/L=Brussels/O=Client organization/OU=Test client unit/CN=Test client"

echo "creating REST client certificate signed by fake root CA"
openssl x509 \
    -req \
    -in client.csr \
    -set_serial 01 \
    -sha384 \
    -CA rootCACert.pem \
    -CAkey rootCAKey.pem \
    -days 365 \
    -out clientCert.pem 

echo "using Root CA certificate as trustore for the REST service"
cp rootCACert.pem weather-rest-service-truststore.pem
