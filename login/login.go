package login

import "fmt"

//LoginKit defines the functions needed to form and execute a login request
type LoginKit interface {
	RequestBuilder(string, string) string
	EndpointBuilder(string) string
}

var loginKits map[string]LoginKit

func init() {
	loginKits = make(map[string]LoginKit)
	loginKits["xtrac"] = &XtracLoginKit{}
}

//XtracLoginKit defines LoginKit methods for logging into XTRAC
type XtracLoginKit struct{}

//RequestBuilder builds a login request for logging into XTRAC via the loginXt service
func (xt *XtracLoginKit) RequestBuilder(username string, password string) string {
	return `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns="http://xmlns.fmr.com/systems/dev/xtrac/2004/06/" xmlns:ser="http://xmlns.fmr.com/common/headers/2005/12/ServiceProcessingDirectives" xmlns:ser1="http://xmlns.fmr.com/common/headers/2005/12/ServiceCallContext" xmlns:typ="http://xmlns.fmr.com/systems/dev/xtrac/2004/06/types">
   <soapenv:Header/>
   <soapenv:Body>
      <ns:loginXt>
         <ns:credentials>
            <typ:operatorName>` + username + `</typ:operatorName>
            <typ:password>` + password + `</typ:password>
          </ns:credentials>
      </ns:loginXt>
   </soapenv:Body>
</soapenv:Envelope>`
}

//EndpointBuilder builds a login request for logging into XTRAC via the loginXt service
func (xt *XtracLoginKit) EndpointBuilder(hostportSpec string) string {
	return fmt.Sprintf("http://%s/XtracWeb/services/Login", hostportSpec)
}

//GetLoginKit returns the kit for the given provider, if one is present in the loginKits map.
func GetLoginKit(provider string) LoginKit {
	return loginKits[provider]
}

//SupportedProvider returns true is the given provider is supported.
func SupportedProvider(provider string) bool {
	kit := loginKits[provider]
	return kit != nil
}
