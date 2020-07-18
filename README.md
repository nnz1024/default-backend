# Custom error pages backend for ingress-nginx

Based on
[ingress-nginx's custom error pages backend](https://github.com/kubernetes/ingress-nginx/tree/master/images/custom-error-pages),
with some improvements:
* Added a new build system, based exclusively on Dockerfile, instead of old 
  Makefile-based one (see [Build section](#Build)).
* Implemented support for MIME types which resolves to more than one file 
  extension (original project uses the first one given by 
  `mime.ExtensionsByType`, i.e. random).
* Consequently, fixed bug with random `404 page not found` error instead of
  your nice HTML error pages (caused by `html` vs `htm` ambiguity).
* Added support for Go template engine, which allows you to use some special
  headers given by ingress-nginx (such as `X-Namespace`) to alter page contents.
  See in [Build section](#Build) how to enable it.

# Custom error (and redirection) pages

Place them under `rootfs/www`. Main points:
* Use lowercase in file names.
* Name pages according to HTTP status codes.
* If page for exact code (for example, `501.html`) is absent, server will try
  to use "more general" name ("3xx" for redirections, "4xx" for client errors,
  "5xx" for server errors). If it is also absent, "404" will be returned.
* If client declares `Accept` (and ingress-nginx sets `X-Format`) header
  different than `text/html`, server will try to found and serve pages with 
  appropriate extension, according to `rootfs/etc/mime.types`.
* If you want to edit `mime.types`, please note that some file extenstions
  are [hardcoded](https://github.com/golang/go/blob/master/src/mime/type.go#L60)
  in Go, so you cannot remap or get rid of them by simply editing `mime.types`.
* If content type is unknown or corresponding files are absent, `*.html`/`*.htm`
  files will be used.
* Backend server will run under non-privileged user (UID 101 by default, you
  can change it via `--build-arg=UID=XXX` parameter in `docker build` command),
  so, in the case of runtime errors, make sure that your files and directories
  are accessible to non-owner and non-group users for reading
  (`chmod -R o=rX rootfs`).

# Template engine

You can use [template engine](https://golang.org/pkg/text/template/) (disabled
by default, see in [Build section](#Build) how to enable it). Templating
settings are declared in file `src/templating_settings.go`, in `Templating` map.
For example,
```
"text/html": {
	enableVar:  "TEMPLATE_HTML",
	escapeType: escapeHTML,
},
```
means that for HTML pages (`*.html` and `*.htm`) template will be used only
if environment variable `TEMPLATE_HTML` set to any non-empty value (on/off
switch), and substituted values will be escaped with
[html/template](https://golang.org/pkg/html/template/)
semantical escaping. Another escape option is `escapeJSON`, it simply uses
[JSEscapeString](https://golang.org/pkg/text/template/#JSEscapeString) function.
And, of course, you can always disable escaping for given MIME type via
`escapeNone` value.

List of accessible values can be seen in `src/constants.go` in `Headers` map.
Key is "pretty" name, which will be used in the template as the variable name,
and value is the canonical HTTP header name. For example, the content of
`X-Original-URI` header (must be populated by ingress-nginx while proxying 
request to this backend) can be referenced in the template as
`{{.OriginalURI}}`.

# Build

Building version x.y without template engine:
```
docker build . -t "your-repo/default-backend:x.y"
```

Building version x.y with template engine:
```
docker build . -t "your-repo/default-backend:x.y" --build-arg="TAGS=template"
```

# Deployment

```
---
apiVersion: v1
kind: Service
metadata:
  name: default-backend
  namespace: ingress-nginx
  labels:
    app: default-backend
spec:
  ports:
    - name: http
      port: 80
      targetPort: http
  selector:
    app: default-backend
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: default-backend
  namespace: ingress-nginx
  labels:
    app: default-backend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: default-backend
  template:
    metadata:
      labels:
        app: default-backend
    spec:
      containers:
        - name: default-backend
          # Any image is permissible as long as:
          # 1. It serves a 404 page at /
          # 2. It serves 200 on a /healthz endpoint
          # image: gcr.io/google_containers/defaultbackend:1.4
          image: "your-repo/default-backend:x.y"
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8080
              scheme: HTTP
            initialDelaySeconds: 30
            timeoutSeconds: 5
          ports:
            - name: http
              containerPort: 8080
```

Then you can set it as your default HTTP backend for all your namespaces via
`--default-backend-service=ingress-nginx/default-backend` ingress-nginx command
line option (to be clear, form is `$NAMESPACE/$SERVICE`).

If you want to use specific error pages for some ingresses, create
default-backend deployment and service in their namespaces, and add
```
nginx.ingress.kubernetes.io/default-backend: default-backend
```
annotation to these ingresses (see [documentation](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#default-backend)).

Also don't forget that custom error pages will be used only for HTTP
statuses listed in [nginx.ingress.kubernetes.io/custom-http-errors](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#custom-http-errors)
annotation (per ingress) or in [custom-http-errors](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#custom-http-errors)
ingress-nginx main ConfigMap setting (global).
