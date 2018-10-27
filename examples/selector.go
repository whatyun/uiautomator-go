package main

import (
	ug "uiautomator"
)

func main() {
	client := ug.New(&ug.Config{
		Host:      "10.10.20.78",
		Port:      7912,
		AutoRetry: 0,
		Timeout:   10,
	})

	ele, err := client.GetElementBySelector(
		map[string]interface{}{
			"resourceId": "com.android.chrome:id/url_bar",
		},
	)
	if err != nil {
		panic(err)
	}

	/*
		// Get child element

		ele, err = ele.ChildByText(
			"Clock",
			map[string]interface{}{
				"className": "android.widget.FrameLayout",
			},
		)
	*/

	/*
		// Get element by index

		ele, err = ele.Eq(0)
		if err != nil {
			panic(err)
		}
	*/

	/*
		// Get text

		text, err := ele.GetText()
		if err != nil {
			panic(err)
		}
		fmt.Println(text)
	*/

	/*
		// Set text

		err = ele.SetText("https://www.google.com/")
		if err != nil {
			panic(err)
		}
	*/

	/*
		// Long click

		err = ele.LongClick()
		if err != nil {
			panic(err)
		}
	*/

	/*
		// Swipe element

		err = ele.SwipeLeft()
		if err != nil {
			panic(err)
		}
	*/

	// Clear the text input
	err = ele.ClearText()
	if err != nil {
		panic(err)
	}
}
