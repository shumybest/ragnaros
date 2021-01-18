# ===================================================================
# Spring Boot configuration.
#
# This configuration will be overridden by the Spring profile you use,
# for example application-dev.yml if you use the "dev" profile.
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
  client:
    enabled: true
    healthcheck:
      enabled: true
    fetch-registry: true
    register-with-eureka: true
    instance-info-replication-interval-seconds: 10
    registry-fetch-interval-seconds: 10
  instance:
    appname: {{ .App.ProjectName }}
    instanceId: {{ .App.ProjectName }}:${spring.application.instance-id:${random.value}}
    lease-renewal-interval-in-seconds: 5
    lease-expiration-duration-in-seconds: 10
    status-page-url-path: ${management.endpoints.web.base-path}/info
    health-check-url-path: ${management.endpoints.web.base-path}/health
    metadata-map:
      zone: primary # This is needed for the load balancer
      profile: ${spring.profiles.active}
      version: ${info.project.version:}
      git-version: ${git.commit.id.describe:}
      git-commit: ${git.commit.id.abbrev:}
      git-branch: ${git.branch:}
ribbon:
  ReadTimeout: 5000
  ConnectTimeout: 3000
feign:
  client:
    config:
      default:
        connectTimeout: 5000
        readTimeout: 20000

management:
  endpoints:
    web:
      base-path: /management
  endpoint:
    health:
      show-details: when-authorized

spring:
  application:
    name: {{ .App.ProjectName }}

security:
  basic:
    enabled: false

# ===================================================================
# Application specific properties
# Add your own application properties here, see the ApplicationProperties class
# to have type-safe configuration, like in the JHipsterProperties above
#
# More documentation is available at:
# https://www.jhipster.tech/common-application-properties/
# ===================================================================

application:
