# Authetication On Kubernetes

https://kubernetes.io/docs/reference/access-authn-authz/authentication/#x509-client-certs

https://kubernetes.io/docs/reference/access-authn-authz/certificate-signing-requests/#kubernetes-signers

https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/#create-a-certificate-signing-request


## Pre-Req

- **openssl**

## Let's Go

The first step is create a CSR (Certificate Signing Request)

```bash
openssl req -nodes \
	  -days 365 \
	  -newkey rsa:2048 \
	  -keyout estagiario.key \
	  -out estagiario.csr \
	  -subj '/CN=estagiario/O=<org-name>' \
	  -addext 'subjectAltName = DNS:<>'
```

After create a CSR, we need send this CSR for k8s cluster. The k8s cluster is who will assigned this.

1. We need create a yaml file for that

Ex.

```yaml
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: estagiario-csr
spec:
  request: <csr aqui>
  signerName: kubernetes.io/kube-apiserver-client
  usages:
  - client auth
```

We need to inform the csr in one line, but how do this? We need to use the command

```bash
cat estagiario.csr | base64 -w 0
```

> -w 0 told for the base64 command: 'Hey your output needs to be in one line'

2. After created the **estagiario.csr** we need to approve the CSR

```bash
kubenetes get certificate
kubectl certificate approve estagiario-csr
kubectl get csr estagiario-csr -ojson |\
		jq -r .status.certificate |\
		base64 -d > estagiario.crt
openssl x509 -in estagiario.crt -text
```

3. Now, we have the cert and the key

4. We need to create a context, user and cluster in own pc

```bash
$ kubectl config set-credentials estagiario \
		--client-certificate estagiario.crt \
		--client-key estagiario.key
	
$ kubectl config set-context estagiario \
		--cluster kind-kind \
		--user estagiario
 
$ kubectl config use-context estagiario
 
$ kubectl get pods
```

Ban: Error happened

```txt
Error from server (Forbidden): pods is forbidden: User "estagiario" cannot list resource "pods" in API group "" in the namespace "default"
```

Now, We need to give permission for the estagiario


Acess the path [rbac](../rbac/cluster-role-pod-reader.yaml) to see how it is