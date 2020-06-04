package main

type TelegramWorker struct {
	updatesGetter UpdatesGetter
	updateHandler UpdateHandler
}

func NewTelegramWorker(getter UpdatesGetter, handler UpdateHandler) *TelegramWorker {
	return &TelegramWorker{
		updatesGetter: getter,
		updateHandler: handler,
	}
}

func (worker *TelegramWorker) Work() error {
	updates, err := worker.updatesGetter.GetUpdates()
	if err != nil {
		return err
	}

	return worker.updateHandler.HandleUpdates(updates)
}
