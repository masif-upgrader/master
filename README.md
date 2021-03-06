## About

The **Masif Upgrader master** is a component of *Masif Upgrader*.

Consult Masif Upgrader's [manual] on its purpose
and the master's role in its architecture
and [demo] for a full stack live demonstration.

## Configuration

The configuration file (usually `/etc/masif-upgrader-master/config.ini`)
looks like this:

```ini
[api]
listen=0.0.0.0:8150

[tls]
cert=/var/lib/puppet/ssl/certs/infra-mgmt.intern.example.com.pem
key=/var/lib/puppet/ssl/private_keys/infra-mgmt.intern.example.com.pem
ca=/var/lib/puppet/ssl/certs/ca.pem
crl=/var/lib/puppet/ssl/ca/ca_crl.pem

[db]
type=mysql
dsn=masif_upgrader_master:123456@/masif_upgrader

[log]
level=info
```

*api.listen* is the address (HOST:PORT) to listen on for requests from agents.

The *tls* section describes the X.509 PKI:

 option | description
 -------|------------------------------------------------------------
 cert   | TLS server certificate chain (may include root CA)
 key    | TLS server private key
 ca     | TLS client root CA certificate
 crl    | TLS client root CA's certificate revocation list (optional)

The *db* section describes the database the master shares with the UI:

 option | description
 -------|-----------------------------------
 type   | The database's type (only "mysql")
 dsn    | The database's [DSN]

*log.level* defines the logging verbosity and is one of:

* error
* warning
* info
* debug

## Docker

```bash
docker run --rm -d \
  -v /var/lib/puppet/ssl:/pki:ro \
  -e MASIF_MASTER_API_LISTEN=0.0.0.0:8150 \
  -e MASIF_MASTER_TLS_CERT=/pki/certs/infra-mgmt.intern.example.com.pem \
  -e MASIF_MASTER_TLS_KEY=/pki/private_keys/infra-mgmt.intern.example.com.pem \
  -e MASIF_MASTER_TLS_CA=/pki/certs/ca.pem \
  -e MASIF_MASTER_TLS_CRL=/pki/ca/ca_crl.pem \
  -e MASIF_MASTER_DB_TYPE=mysql \
  -e MASIF_MASTER_DB_DSN=masif_upgrader_master:123456@192.0.2.2/masif_upgrader \
  -e MASIF_MASTER_LOG_LEVEL=info \
  masifupgrade/master
```

[manual]: https://github.com/masif-upgrader/manual
[demo]: https://github.com/masif-upgrader/demo
[DSN]: https://github.com/go-sql-driver/mysql#dsn-data-source-name
