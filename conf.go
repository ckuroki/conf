package conf

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// default values tag name
const tagName = "default"

var (
	ErrInvalidValue = errors.New("invalid value")
	ErrUnsupported  = errors.New("unsupported type")
	ErrUnexported   = errors.New("unexported field")
	matchFirstCap   = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap     = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// Unmarshal extracts defined fields from environment into an struct field
// , if a field is not defined at env uses value defined in "default" tag
// Supported data types are all integer types, float,  string, bool, maps and slices
// Nested structs can be used too.
// env vars are generated by convention using a prefix + field names transformed to snake case and uppercased
//
// e.g:
// cfg := struct {
//	apiPort				int			`default:"8080"`
//	srvShortName	string	`default:"local"`
// }{}
// err := conf.Unmarshal(&cfg, "MYAPP")
// cfg.apiPort will take MYAPP_API_PORT value or 8080 if it not defined
// cfg.srvShortName will take MIAPP_SRV_SHORT_NAME value o "local"
//
// Check conf_test.go to see more examples
func Unmarshal(cfg interface{}, prefix string) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return ErrInvalidValue
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrInvalidValue
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)
		tf := t.Field(i)
		switch vf.Kind() {
		case reflect.Struct:
			if !vf.Addr().CanInterface() {
				continue
			}
			iface := vf.Addr().Interface()
			err := Unmarshal(iface, toEnvVarName(tf.Name, prefix))
			if err != nil {
				return err
			}
		}
		tagVal := tf.Tag.Get(tagName)
		if tagVal == "" {
			continue
		}
		if !vf.CanSet() {
			return ErrUnexported
		}
		// Get Environment Value
		val := os.Getenv(toEnvVarName(tf.Name, prefix))
		// If env var not set then use default value
		if val == "" {
			val = tagVal
		}
		err := set(tf.Type, vf, val)
		if err != nil {
			return err
		}
	}
	return nil
}

// set updates a field with value defined in val
func set(t reflect.Type, f reflect.Value, val string) error {
	switch t.Kind() {
	case reflect.Ptr:
		ptr := reflect.New(t.Elem())
		err := set(t.Elem(), ptr.Elem(), val)
		if err != nil {
			return err
		}
		f.Set(ptr)
	case reflect.String:
		f.SetString(val)
	case reflect.Bool:
		v, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		f.SetBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		f.SetInt(int64(v))
	case reflect.Map:
		m := reflect.MakeMap(t)
		if len(strings.TrimSpace(val)) != 0 {
			pairs := strings.Split(val, ",")
			for _, pair := range pairs {
				kvpair := strings.Split(pair, ":")
				if len(kvpair) != 2 {
					return fmt.Errorf("invalid map item: %q", pair)
				}
				k := reflect.New(t.Key()).Elem()
				err := set(t.Key(), k, kvpair[0])
				if err != nil {
					return err
				}
				v := reflect.New(t.Elem()).Elem()
				err = set(t.Elem(), v, kvpair[1])
				if err != nil {
					return err
				}
				m.SetMapIndex(k, v)
			}
		}
		f.Set(m)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(val, t.Bits())
		if err != nil {
			return err
		}
		f.SetFloat(v)
	case reflect.Slice:
		s := reflect.MakeSlice(t, 0, 0)
		if t.Elem().Kind() == reflect.Uint8 {
			s = reflect.ValueOf([]byte(val))
		} else if len(strings.TrimSpace(val)) != 0 {
			vals := strings.Split(val, ",")
			s = reflect.MakeSlice(t, len(vals), len(vals))
			for i, v := range vals {
				err := set(t.Elem(), s.Index(i), v)
				if err != nil {
					return err
				}
			}
		}
		f.Set(s)
	default:
		return ErrUnsupported
	}

	return nil
}

// toEnvVarName transforms a camelCase name to upper snake case
func toEnvVarName(str string, prefix string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return prefix + "_" + strings.ToUpper(snake)
}