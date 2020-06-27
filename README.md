# conf

Minimalistic go configuration package.

##  How it works 

Reads configurations from env vars.
If an env var is not defined will use default values from `default:` tags.

Data types are inferred from configuration struct fields.

Each app can have an optional prefix, if defined all env vars starts with that prefix.

Your struct fields should be public (InitCase) and use idiomatic go naming (camelCase). Env vars are mapped to field names by convention.
Field names are transformed to snake case , upper cased and preceeded by prefix.

e.g.:
```
App. prefix : MYAPP
Field name : ApiPortNumber
Env var : MYAPP_API_PORT_NUMBER
```

## Supported types

### Integer and float types, string, bool

```
Count           int     `default: "10"`       // MYAPP_COUNT 
MaxBytes        int64   `default: "4194304"`  // MYAPP_MAX_BYTES
MinAmount       float64 `default: "50.0"`     // MYAPP_MIN_AMOUNT
Hostname        string  `default: "localhost" // MYAPP_HOSTNAME
EnableCoolMode  bool    `default: "true"`     //  MYAPP_ENABLE_COOL_MODE
```

### Nested structs

```
type GeoApi struct {
  Url     string  `default: "http://192.168.0.100"`
  Token   string  `default: "abcdefghijklmnopq"
}

cfg := Config {
  ApiPort       int       // MYAPP_API_PORT
  ApiHostname   string    // MYAPP_HOSTNAME
  GeoApi        GeoApi    // MYAPP_GEOAPI_URL
                          // MYAPP_GEOAPI_TOKEN
}{}
```

### Slices

```
IntSlice  []int     `default: "5,6,7,8"` // MYAPP_INT_SLICE=7,8,9
StrSlice  []string  `default: "a,b,c"`   // MYAPP_STR_SLICE=a,b,c
```

### Maps

```
Countries map[string]string // MYAPP_COUNTRIES="ar:Argentina,es:Spain,it:Italy"
Prefixes  map[string]int    // MYAPP_PREFIXES="ar:54,es:34,it:39"
```

