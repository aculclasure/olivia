package modules

import (
	"encoding/json"

	"github.com/olivia-ai/olivia/user"
)

var RollingRecycleTag = "rolling recycle"

func RollingRecycleRun(msg, token string, session Session) error {
	resp := ResponseMessage{
		Content:     "Ok let's do a rolling recycle!",
		Tag:         "rolling recycle",
		Information: user.GetUserInformation(token),
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	if err = session.Send(bytes); err != nil {
		return err
	}
	return nil
}
