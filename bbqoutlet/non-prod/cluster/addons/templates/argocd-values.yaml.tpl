controller:
  enableStatefulSet: true
  logformat: json
  nodeSelector:
    fusion_node_type: infra
  metrics:
    enabled: true

redis:
  nodeSelector:
    fusion_node_type: infra

server:
  nodeSelector:
    fusion_node_type: infra
  metrics:
    enabled: true
  additionalApplications:
    - name: managed-fusion-applications
      namespace: argocd
      project: default
      source:
        repoURL: ${repoURL}
        targetRevision: HEAD
        path: ${clusterType}/applications
        directory:
          recurse: true
      destination:
        name: in-cluster
        namespace: argocd
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
  ingress: 
    enabled: true
    https: true
    annotations:
      kubernetes.io/ingress.class: nginx
      cert-manager.io/cluster-issuer: letsencrypt-production
      nginx.ingress.kubernetes.io/backend-protocol: "HTTPS"
      nginx.ingress.kubernetes.io/ssl-passthrough: "true"
    hosts:
      - ${argocdHost}
    tls:
      - secretName: ${argocdHostTlsName}
        hosts:
          - ${argocdHost}
  config:
    admin.enabled: "false"
    url: https://${argocdHost}
    oidc.config: |
      name: Okta
      issuer: https://lucidworks.okta.com/oauth2/default
      clientID: ${oktaClientID}
      clientSecret: ${oktaClientSecret}
      requestedIDTokenClaims:
        groups:
          essential: true
      requestedScopes:
        - "openid"
        - "profile"
        - "email"
        - "groups"
      logoutURL: https://${argocdHost}
  rbacConfig:
    policy.default: role:readonly
    policy.csv: |
      g, cc_ops, role:admin
    scopes: "[groups, email]"
repoServer:
  nodeSelector:
    fusion_node_type: infra
  metrics:
    enabled: true

dex:
  enabled: false

configs:
  credentialTemplates:
    ${repoName}:
      url: ${repoURL}
      sshPrivateKey: ${repoSshPrivateKey}
      sshPublicKey: ${repoSshPublicKey}
    config-sync-config:
      url: ${configSyncRepoURL}
      sshPrivateKey: ${configSyncSshPrivateKey}
      sshPublicKey: ${configSyncSshPublicKey}
    fusion-dev-helm-v2:
      url: https://artifactory.lucidworks.com:443/artifactory/api/helm/fusion-dev-helm
      username: ${fusionDevHelmRepoUsername}
      password: ${fusionDevHelmRepoPassword}
  repositories:
    ${repoName}:
      url: ${repoURL}
      name: ${repoName}
      type: git
    config-sync-config:
      url: ${configSyncRepoURL}
      name: config-sync-config
      type: git
    fusion-dev-helm-v2:
      url: https://artifactory.lucidworks.com:443/artifactory/api/helm/fusion-dev-helm
      name: fusion-dev-helm-v2
      type: helm
  secret:
    githubSecret: "${githubWebhookSecret}"