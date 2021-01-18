apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: {{ .App.ProjectName }}
  namespace: default
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: {{ .App.ProjectName }}
    spec:
      imagePullSecrets:
      - name: ali-registry-vpc-secret
      containers:
      - name: {{ .App.ProjectName }}-app
        image: registry.cn-hangzhou.aliyuncs.com/repo/{{ .App.ProjectName }}:latest
        imagePullPolicy: Always
        env:
        - name: SERVER_PORT
          value: "{{ .K8s.Server.Port }}"
        - name: SPRING_PROFILES_ACTIVE
          value: {{ .K8s.Spring.Profiles }}
        - name: SPRING_CLOUD_CONFIG_URI
          value: {{ .K8s.Spring.CloudConfigUri }}
        - name: EUREKA_CLIENT_SERVICE_URL_DEFAULTZONE
          value: {{ .K8s.Eureka.ServiceUrl }}
        - name: SPRING_DATASOURCE_URL
          value: {{ .K8s.Spring.DataSourceUrl }}
        ports:
        - name: web
          containerPort: {{ .K8s.Server.Port }}
        readinessProbe:
          httpGet:
            path: /management/health
            port: web
        livenessProbe:
          httpGet:
            path: /management/health
            port: web
          initialDelaySeconds: 180
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .App.ProjectName }}
  namespace: default
  labels:
    app: {{ .App.ProjectName }}
spec:
  selector:
    app: {{ .App.ProjectName }}
  ports:
    - name: web
      port: {{ .K8s.Server.Port }}
