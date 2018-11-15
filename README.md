# Go Marshaller
A simple library for marshalling / unmarshalling objects to/from byte arrays in go

## Example
~~~go
import (
  ...
  reqmarshaller "github.com/c3systems/c3-utils-go-marshaller/http/req"
  ...
)

func foo(r *http.Req) error {
  tr, err := reqmarshal.TransformRequest(r)
  if err != nil {
    return err
  }

  b, err := tr.Marshal()
  if err != nil {
    return err
   }

   // do something with the bytes array, b
   // ...

}

func bar(b []byte) error {
  tr := new(reqmarshal.TransformedRequest)
  if err := tr.Unmarshal(b); err != nil {
    return err
   }

   r, err := reqmarshal.UntransformRequest(tr)
   if err != nil {
    return err
   }

   // do something with the request, r
   // ...
}
~~~
