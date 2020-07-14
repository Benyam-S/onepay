package app

// DeleteOnePayAccount is a method that deletes a certain onepay account.
// Also it returns all the needed data for user closing his/her account
// func (onepay *OnePay) DeleteOnePayAccount(userID string) (*entity.User,
// 	[]*entity.UserHistory, []*entity.LinkedAccount, error) {

// 	opWallet, err := onepay.WalletService.FindWallet(userID)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}

// 	// checking first if the user wallet is empty
// 	if opWallet.Amount > 0 {
// 		return nil, nil, nil, errors.New("please empty your wallet before deleting account")
// 	}

// 	// checking first if the user have any money token's that hasn't been reclaim
// 	moneyTokens := onepay.MoneyTokenService.SearchMoneyToken(userID)
// 	if len(moneyTokens) > 0 {
// 		return nil, nil, nil, errors.New("please delete or reclaim all money tokens that has not been received before deleting account")
// 	}

// 	// Reading all the user's linked accounts
// 	linkedAccounts := onepay.LinkedAccountService.SearchLinkedAccounts("user_id", userID)

// 	// Reading all the user's histories
// 	userHistories := onepay.HistoryService.AllUserHistories(userID)

// 	opUser, err := onepay.userService.DeleteUser(userID)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}

// 	// Adding deleted user to trash
// 	onepay.deletedService.AddUserToTrash(opUser)

// 	// Deleting the user's wallet
// 	onepay.WalletService.DeleteWallet(userID)

// 	// Removing all the user's linked accounts
// 	onepay.LinkedAccountService.DeleteLinkedAccounts(userID)

// 	for _, linkedAccount := range linkedAccounts {
// 		// Adding linked account to trash
// 		onepay.deletedService.AddLinkedAccountToTrash(linkedAccount)
// 	}

// 	return opUser, userHistories, linkedAccounts, nil
// }
