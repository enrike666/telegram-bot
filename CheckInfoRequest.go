package main

type CheckInfoRequest struct {
	DateTimeString string `schema:"t"`
	Sum            string `schema:"s"`
	FN             string `schema:"fn"`
	FD             string `schema:"fd"`
	FP             string `schema:"fp"`
	N              string `schema:"n"`
}

func NewCheckInfoRequestBasedCheckInfoFromBot(checkInfo CheckInfoFromBot) (*CheckInfoRequest, error) {
	dateTime, err := checkInfo.GetTime()
	if err != nil {
		return nil, err
	}

	dateTimeStr := dateTime.Format("01.06.2006 15:04")

	return &CheckInfoRequest{
		DateTimeString: dateTimeStr,
		Sum:            checkInfo.Sum,
		FN:             checkInfo.FN,
		FD:             checkInfo.FD,
		FP:             checkInfo.FP,
		N:              checkInfo.N,
	}, nil
}
