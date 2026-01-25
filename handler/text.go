package handler

const (
	TextUnknownCommand = "Unknown command. Use /start to begin."

	TextGreeting        = "Hello, "
	TextWelcome         = "Welcome to the XLS File Reader Bot!\n"
	TextWelcomeDesc     = "I'm here to help you process Excel files.\n\n"
	TextFunctionsHeader = "📋 Available functions:\n"
	TextFunction1       = "• Send me an Excel file (.xls, .xlsx) to read and process\n"
	TextFunction2       = "• I will extract and display the data for you\n"
	TextFunction3       = "• Use /start to see this message again"

	TextInstructionsHeader = "📖 Bot Instructions\n\n"
	TextInstructionsDesc   = "This bot helps you read and process Excel files.\n\n"
	TextInstructionsFuncs  = "📋 Functions:\n"
	TextInstructionsFunc1  = "• Send Excel files (.xls, .xlsx) - I will read and display the data\n"
	TextInstructionsFunc2  = "• File processing - Extract information from your spreadsheets\n"
	TextInstructionsFunc3  = "• Data display - View your Excel data in a readable format\n\n"
	TextInstructionsTip    = "💡 To get started, use /start command or simply send me an Excel file!"
)

func GetWelcomeText(username string) string {
	if username == "" {
		username = "there"
	}

	text := TextGreeting + username + "! 👋\n\n"
	text += TextWelcome
	text += TextWelcomeDesc
	text += TextFunctionsHeader
	text += TextFunction1
	text += TextFunction2
	text += TextFunction3

	return text
}

func GetInstructionsText() string {
	text := TextInstructionsHeader
	text += TextInstructionsDesc
	text += TextInstructionsFuncs
	text += TextInstructionsFunc1
	text += TextInstructionsFunc2
	text += TextInstructionsFunc3
	text += TextInstructionsTip

	return text
}
