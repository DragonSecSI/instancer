package helpers

import (
	"github.com/DragonSecSI/instancer/backend/pkg/helpers/api"
	"github.com/DragonSecSI/instancer/backend/pkg/helpers/auth"
	"github.com/DragonSecSI/instancer/backend/pkg/helpers/flag"
)

var Auth = auth.NewAuthHelper()
var Flag = flag.NewFlagHelper()
var Api = api.NewApiHelper()
