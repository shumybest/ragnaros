# ===================================================================
# Spring Boot configuration for the "prod" profile.
#
# This configuration overrides the application.yml file.
#
# More information on profiles: https://www.jhipster.tech/profiles/
# More information on configuration properties: https://www.jhipster.tech/common-application-properties/
# ===================================================================

# ===================================================================
# Standard Spring Boot properties.
# Full reference is available at:
# http://docs.spring.io/spring-boot/docs/current/reference/html/common-application-properties.html
# ===================================================================
eureka:
  instance:
    prefer-ip-address: true
  client:
    service-url:
      defaultZone: http://admin:${jhipster.registry.password}@${registry-host}:8761/eureka/

spring:
  datasource:
    url: jdbc:mysql://${registry-host}:3306/ragnaros?useUnicode=true&characterEncoding=utf8&useSSL=false&useLegacyDatetimeCode=false&serverTimezone=Asia/Shanghai
    username: mysqlUserHere
    password: mysqlPasswordHere

# ===================================================================
# To enable TLS in production, generate a certificate using:
# keytool -genkey -alias ragnaros -storetype PKCS12 -keyalg RSA -keysize 2048 -keystore keystore.p12 -validity 3650
#
# You can also use Let's Encrypt:
# https://maximilian-boehm.com/hp2121/Create-a-Java-Keystore-JKS-from-Let-s-Encrypt-Certificates.htm
#
# Then, modify the server.ssl properties so your "server" configuration looks like:
#
# server:
#    port: 443
#    ssl:
#        key-store: classpath:config/tls/keystore.p12
#        key-store-password: password
#        key-store-type: PKCS12
#        key-alias: ragnaros
#        # The ciphers suite enforce the security by deactivating some old and deprecated SSL cipher, this list was tested against SSL Labs (https://www.ssllabs.com/ssltest/)
#        ciphers: TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384 ,TLS_DHE_RSA_WITH_AES_128_GCM_SHA256 ,TLS_DHE_RSA_WITH_AES_256_GCM_SHA384 ,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384,TLS_DHE_RSA_WITH_AES_128_CBC_SHA256,TLS_DHE_RSA_WITH_AES_128_CBC_SHA,TLS_DHE_RSA_WITH_AES_256_CBC_SHA256,TLS_DHE_RSA_WITH_AES_256_CBC_SHA,TLS_RSA_WITH_AES_128_GCM_SHA256,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_CBC_SHA256,TLS_RSA_WITH_AES_256_CBC_SHA256,TLS_RSA_WITH_AES_128_CBC_SHA,TLS_RSA_WITH_AES_256_CBC_SHA,TLS_DHE_RSA_WITH_CAMELLIA_256_CBC_SHA,TLS_RSA_WITH_CAMELLIA_256_CBC_SHA,TLS_DHE_RSA_WITH_CAMELLIA_128_CBC_SHA,TLS_RSA_WITH_CAMELLIA_128_CBC_SHA
# ===================================================================
server:
  port: {{ .K8s.Server.Port }}

# ===================================================================
# JHipster specific properties
#
# Full reference is available at: https://www.jhipster.tech/common-application-properties/
# ===================================================================

jhipster:
  security:
    authentication:
      jwt:
        # This token must be encoded using Base64 and be at least 256 bits long (you can type `openssl rand -base64 64` on your command line to generate a 512 bits one)
        # As this is the PRODUCTION configuration, you MUST change the default key, and store it securely:
        # - In the JHipster Registry (which includes a Spring Cloud Config server)
        # - In a separate `application-prod.yml` file, in the same folder as your executable WAR file
        # - In the `JHIPSTER_SECURITY_AUTHENTICATION_JWT_BASE64_SECRET` environment variable
        base64-secret: ZTViYTRhNDZkMWQ2N2ViMzkzZmVhNWRkYTI0YzU2NTE3YzFmNmQxNDZkMWQ2N2E0NmQxZDY3ZWIzOWJhNGE0OTNmZWE1ZAo=
        # Token is valid 24 hours
        token-validity-in-seconds: 86400
        token-validity-in-seconds-for-remember-me: 2592000

# ===================================================================
# Application specific properties
# Add your own application properties here, see the ApplicationProperties class
# to have type-safe configuration, like in the JHipsterProperties above
#
# More documentation is available at:
# https://www.jhipster.tech/common-application-properties/
# ===================================================================

ragnaros:
  elasticsearch:
    url: http://localhost:9200
    host:
    port:
    username:
    password:

application:
