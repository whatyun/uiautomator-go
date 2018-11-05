/**
Selector is a handy mechanism to identify a specific UI object in the current window.
https://github.com/openatx/uiautomator2#selector
*/
package uiautomator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	Selector map[string]interface{}

	Element struct {
		ua       *UIAutomator
		position *Position
		selector Selector
		original Selector
	}

	ElementRect struct {
		Bottom int `json:"bottom"`
		Left   int `json:"left"`
		Right  int `json:"right"`
		Top    int `json:"top"`
	}

	ElementInfo struct {
		ContentDescription string       `json:"contentDescription"`
		Checked            bool         `json:"checked"`
		Scrollable         bool         `json:"scrollable"`
		Text               string       `json:"text"`
		PackageName        string       `json:"packageName"`
		Selected           bool         `json:"selected"`
		Enabled            bool         `json:"enabled"`
		ClassName          string       `json:"className"`
		Focused            bool         `json:"focused"`
		Focusable          bool         `json:"focusable"`
		Clickable          bool         `json:"clickable"`
		ChileCount         int          `json:"chileCount"`
		LongClickable      bool         `json:"longClickable"`
		Checkable          bool         `json:"checkable"`
		Bounds             *ElementRect `json:"bounds"`
		VisibleBounds      *ElementRect `json:"visibleBounds"`
	}
)

var _MASK = map[string]int{
	"text":                  0x01,       // MASK_TEXT,
	"textContains":          0x02,       // MASK_TEXTCONTAINS,
	"textMatches":           0x04,       // MASK_TEXTMATCHES,
	"textStartsWith":        0x08,       // MASK_TEXTSTARTSWITH,
	"className":             0x10,       // MASK_CLASSNAME
	"classNameMatches":      0x20,       // MASK_CLASSNAMEMATCHES
	"description":           0x40,       // MASK_DESCRIPTION
	"descriptionContains":   0x80,       // MASK_DESCRIPTIONCONTAINS
	"descriptionMatches":    0x0100,     // MASK_DESCRIPTIONMATCHES
	"descriptionStartsWith": 0x0200,     // MASK_DESCRIPTIONSTARTSWITH
	"checkable":             0x0400,     // MASK_CHECKABLE
	"checked":               0x0800,     // MASK_CHECKED
	"clickable":             0x1000,     // MASK_CLICKABLE
	"longClickable":         0x2000,     // MASK_LONGCLICKABLE,
	"scrollable":            0x4000,     // MASK_SCROLLABLE,
	"enabled":               0x8000,     // MASK_ENABLED,
	"focusable":             0x010000,   // MASK_FOCUSABLE,
	"focused":               0x020000,   // MASK_FOCUSED,
	"selected":              0x040000,   // MASK_SELECTED,
	"packageName":           0x080000,   // MASK_PACKAGENAME,
	"packageNameMatches":    0x100000,   // MASK_PACKAGENAMEMATCHES,
	"resourceId":            0x200000,   // MASK_RESOURCEID,
	"resourceIdMatches":     0x400000,   // MASK_RESOURCEIDMATCHES,
	"index":                 0x800000,   // MASK_INDEX,
	"instance":              0x01000000, // MASK_INSTANCE,
}

/*
Get element info
*/
func (ele Element) GetInfo() (*ElementInfo, error) {
	var RPCReturned ElementInfo

	if err := ele.ua.post(
		&RPCOptions{
			Method: "objInfo",
			Params: []interface{}{getParams(ele.selector)},
		},
		&RPCReturned,
		nil,
	); err != nil {
		return nil, err
	}

	return &RPCReturned, nil
}

/*
Get Widget rect bounds
*/
func (ele Element) GetRect() (rect *ElementRect, err error) {
	info, err := ele.GetInfo()
	if err != nil {
		return
	}

	rect = info.Bounds
	if rect == nil {
		rect = info.VisibleBounds
	}
	return
}

/*
Get Widget center point
*/
func (ele Element) Center(offset *Position) (*Position, error) {
	rect, err := ele.GetRect()
	if err != nil {
		return nil, err
	}

	lx, ly, rx, ry := rect.Left, rect.Top, rect.Right, rect.Bottom
	width, height := rx-lx, ry-ly

	if offset == nil {
		offset = &Position{0.5, 0.5}
	}

	abs := &Position{}
	abs.X = float32(lx) + float32(width)*offset.X
	abs.Y = float32(ly) + float32(height)*offset.Y
	return abs, nil
}

/*
Get the count
*/
func (ele Element) Count() (int, error) {
	var RPCReturned struct {
		Result int `json:"result"`
	}
	transform := func(response *http.Response) error {
		err := json.NewDecoder(response.Body).Decode(&RPCReturned)
		if err != nil {
			return err
		}
		return nil
	}

	return RPCReturned.Result, ele.ua.post(
		&RPCOptions{
			Method: "count",
			Params: []interface{}{getParams(ele.selector)},
		},
		nil,
		transform,
	)
}

/*
Get the instance via index
*/
func (ele *Element) Eq(index int) (*Element, error) {
	// Update the selector
	ele.original["instance"] = index

	// Recompile the selector
	recompile, err := parseSelector(ele.original)
	if err != nil {
		return nil, err
	}

	ele.selector = recompile
	return ele, nil
}

/*
Check if the specific UI object exists
*/
func (ele Element) WaitForExists() error {
	var RPCReturned struct {
		Result bool `json:"result"`
	}
	transform := func(response *http.Response) error {
		err := json.NewDecoder(response.Body).Decode(&RPCReturned)
		if err != nil {
			return err
		}
		return nil
	}

	err := ele.ua.post(
		&RPCOptions{
			Method: "waitForExists",
			Params: []interface{}{getParams(ele.selector), 20000},
		},
		nil,
		transform,
	)
	if err != nil || RPCReturned.Result == false {
		return &UiaError{
			Code:    -32002,
			Message: "Element not found",
		}
	}

	return nil
}

/*
Swipe the element
*/
func (ele Element) swipe(direction string) error {
	if err := ele.WaitForExists(); err != nil {
		return err
	}
	rect, err := ele.GetRect()
	if err != nil {
		return err
	}

	lx, ly, rx, ry := rect.Left, rect.Top, rect.Right, rect.Bottom
	cx, cy := (lx+rx)/2, (ly+ry)/2

	switch direction {
	case "up":
		return ele.ua.Swipe(
			&Position{X: float32(cx), Y: float32(cy)},
			&Position{X: float32(cx), Y: float32(ly)},
			0.1,
		)
	case "down":
		return ele.ua.Swipe(
			&Position{X: float32(cx), Y: float32(cy)},
			&Position{X: float32(cx), Y: float32(ry - 1)},
			0.1,
		)
	case "left":
		return ele.ua.Swipe(
			&Position{X: float32(cx), Y: float32(cy)},
			&Position{X: float32(lx), Y: float32(cy)},
			0.1,
		)
	case "right":
		return ele.ua.Swipe(
			&Position{X: float32(cx), Y: float32(cy)},
			&Position{X: float32(rx - 1), Y: float32(cy)},
			0.1,
		)
	}

	return nil
}

/*
Swipe to up
*/
func (ele *Element) SwipeUp() error {
	return ele.swipe("up")
}

/*
Swipe to down
*/
func (ele *Element) SwipeDown() error {
	return ele.swipe("down")
}

/*
Swipe to left
*/
func (ele *Element) SwipeLeft() error {
	return ele.swipe("left")
}

/*
Swipe to right
*/
func (ele *Element) SwipeRight() error {
	return ele.swipe("right")
}

/*
Click on the screen
*/
func (ele *Element) Click(offset *Position) error {
	if err := ele.WaitForExists(); err != nil {
		return err
	}

	abs, err := ele.Center(offset)
	if err != nil {
		return err
	}

	return ele.ua.Click(abs)
}

/*
Long click on the element
*/
func (ele *Element) LongClick() error {
	if err := ele.WaitForExists(); err != nil {
		return err
	}

	abs, err := ele.Center(nil)
	if err != nil {
		return err
	}

	return ele.ua.LongClick(abs, 0)
}

/*
Get the children or grandchildren
*/
func (ele *Element) Child(selector Selector) (*Element, error) {
	selector, err := parseSelector(selector)
	if err != nil {
		return nil, err
	}

	ele.selector["childOrSibling"] = []interface{}{"child"}
	ele.selector["childOrSiblingSelector"] = []interface{}{selector}

	return ele, nil
}

func (ele *Element) childByMethod(keywords string, method string, selector Selector) (*Element, error) {
	var RPCReturned struct {
		Result string `json:"result"`
	}
	transform := func(response *http.Response) error {
		err := json.NewDecoder(response.Body).Decode(&RPCReturned)
		if err != nil {
			return err
		}
		return nil
	}

	selector, err := parseSelector(selector)
	if err != nil {
		return nil, err
	}

	if err := ele.ua.post(
		&RPCOptions{
			Method: method,
			Params: []interface{}{ele.selector, selector, keywords, true},
		},
		nil,
		transform,
	); err != nil {
		return nil, err
	}

	ele.selector = map[string]interface{}{"__UID": RPCReturned.Result}
	return ele, nil
}

func (ele *Element) ChildByText(keywords string, selector Selector) (*Element, error) {
	return ele.childByMethod(keywords, "childByText", selector)
}

func (ele *Element) ChildByDescription(keywords string, selector Selector) (*Element, error) {
	return ele.childByMethod(keywords, "childByDescription", selector)
}

/*
Get the sibling
*/
func (ele *Element) Sibling(selector Selector) (*Element, error) {
	selector, err := parseSelector(selector)
	if err != nil {
		return nil, err
	}

	ele.selector["childOrSibling"] = []interface{}{"sibling"}
	ele.selector["childOrSiblingSelector"] = []interface{}{selector}

	return ele, nil
}

/*
Get widget text
*/
func (ele Element) GetText() (string, error) {
	if err := ele.WaitForExists(); err != nil {
		return "", err
	}

	var RPCReturned struct {
		Result string `json:"result"`
	}
	transform := func(response *http.Response) error {
		err := json.NewDecoder(response.Body).Decode(&RPCReturned)
		if err != nil {
			return err
		}
		return nil
	}

	return RPCReturned.Result, ele.ua.post(
		&RPCOptions{
			Method: "getText",
			Params: []interface{}{getParams(ele.selector)},
		},
		nil,
		transform,
	)
}

/*
Set widget text
*/
func (ele Element) SetText(text string) error {
	if err := ele.WaitForExists(); err != nil {
		return err
	}

	return ele.ua.post(
		&RPCOptions{
			Method: "setText",
			Params: []interface{}{getParams(ele.selector), text},
		},
		nil,
		nil,
	)
}

/*
Clear the widget text
*/
func (ele Element) ClearText() error {
	if err := ele.WaitForExists(); err != nil {
		return err
	}

	return ele.ua.post(
		&RPCOptions{
			Method: "clearTextField",
			Params: []interface{}{getParams(ele.selector)},
		},
		nil,
		nil,
	)
}

/*
Query the UI element by selector
*/
func (ua *UIAutomator) GetElementBySelector(selector Selector) (ele *Element, err error) {
	ele = &Element{ua: ua, original: selector}

	selector, err = parseSelector(selector)
	if err != nil {
		return nil, err
	}

	ele.selector = selector
	return
}

func parseSelector(selector Selector) (Selector, error) {
	res := make(Selector)

	// Params initalization
	res["mask"] = 0
	res["childOrSibling"] = []interface{}{}         // Unknow
	res["childOrSiblingSelector"] = []interface{}{} // Unknow

	for k, v := range selector {
		if selectorMask, ok := _MASK[k]; ok {
			res[k] = v
			res["mask"] = res["mask"].(int) | selectorMask
		} else {
			return res, fmt.Errorf("Invalid selector: %v", selector)
		}
	}

	return res, nil
}

func getParams(selector Selector) interface{} {
	if uid, ok := selector["__UID"]; ok {
		return uid
	}
	return selector
}
