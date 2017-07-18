package awsfetch

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/wallix/awless/fetch"
)

func addManualInfraFetchFuncs(sess *session.Session, funcs map[string]fetch.Func) {}

func addManualAccessFetchFuncs(sess *session.Session, funcs map[string]fetch.Func)         {}
func addManualStorageFetchFuncs(sess *session.Session, funcs map[string]fetch.Func)        {}
func addManualMessagingFetchFuncs(sess *session.Session, funcs map[string]fetch.Func)      {}
func addManualDnsFetchFuncs(sess *session.Session, funcs map[string]fetch.Func)            {}
func addManualLambdaFetchFuncs(sess *session.Session, funcs map[string]fetch.Func)         {}
func addManualMonitoringFetchFuncs(sess *session.Session, funcs map[string]fetch.Func)     {}
func addManualCdnFetchFuncs(sess *session.Session, funcs map[string]fetch.Func)            {}
func addManualCloudformationFetchFuncs(sess *session.Session, funcs map[string]fetch.Func) {}
