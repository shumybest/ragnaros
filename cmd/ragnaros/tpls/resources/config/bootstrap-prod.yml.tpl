# ===================================================================
# Spring Cloud Config bootstrap configuration for the "prod" profile
# ===================================================================

spring:
  cloud:
    config:
      fail-fast: true
      retry:
        initial-interval: 1000
        max-interval: 2000
        max-attempts: 100
      uri: http://admin:${jhipster.registry.password}@${registry-host}:8761/config
      # name of the config server's property source (file.yml) that we want to use
      name: {{ .App.ProjectName }}
      profile: prod # profile(s) of the property source
      label: master # toggle to switch to a different version of the configuration as stored in git
      # it can be set to any label, branch or commit of the configuration source Git repository