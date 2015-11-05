package roll

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFinalScopeEmptyInitialScope(t *testing.T) {
	t.Log("Given an empty scope")
	t.Log("When scopeStringWithoutAuthcodeScope is called")
	t.Log("An empty string is returned")
	assert.Equal(t, "", scopeStringWithoutAuthcodeScope(""))
}

func TestFinalScopeWithOnlyAuthcodeScope(t *testing.T) {
	t.Log("Given a scope containing only XtAuthCode")
	t.Log("When scopeStringWithoutAuthcodeScope is called")
	t.Log("An empty string is returned")
	assert.Equal(t, "", scopeStringWithoutAuthcodeScope(XtAuthCodeScope))
}

func TestFinalScopeWithSeveralScopesIncludeingAuthcodeScope(t *testing.T) {
	t.Log("Given a scope containing several scopes including XtAuthCode")
	t.Log("When scopeStringWithoutAuthcodeScope is called")
	t.Log("A string with all the scopes except xtAuthCodeScope is returned")
	assert.Equal(t, "scope1 scope2 scope3", scopeStringWithoutAuthcodeScope("scope1     scope2 "+XtAuthCodeScope+"    scope3"))
}
