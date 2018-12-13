package main

func sanitiseArguments(args string) []string {
	/* Split `args` on its spaces; except double quotes. These should remain intact.
	*/
	var processedArgs []string

	var wordStack []uint8
	for i := range args {
		char := args[i]

		if len(wordStack) > 0 && wordStack[0] == '"' && char == '"'{
			// If we're currently in a glob and find an ending quote
			// end the glob
			// remove leading quote from string
			wordStack = wordStack[1:]
			// Add to processedArgs and clear the wordStack
			processedArgs = append(processedArgs, string(wordStack))
			wordStack = []uint8{}
		} else if char == '"' {
			// If we're not in a glob and find a quote
			// start the glob
			wordStack = append(wordStack, char)
		} else if len(wordStack) > 0 && wordStack[0] == '"' {
			// if we're in a glob and NOT a quote, just append to the stack
			wordStack = append(wordStack, char)
		} else if char == ' ' || i == len(args){
			// Characters in wordStack must be a standalone argument -  dump
			// them into processed args
			processedArgs = append(processedArgs, string(wordStack))
			// Empty word stack
			wordStack = []uint8{}
		} else {
			// Not in a glob, not a space, not starting a glob
			wordStack = append(wordStack, char)
		}
	}

	// Left over word; dump into processedArgs
	if len(wordStack) > 0 {
		processedArgs = append(processedArgs, string(wordStack))
		wordStack = []uint8{}
	}
	return processedArgs
}