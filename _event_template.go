package bot

import (
	ddp "ddp"
	"reflect"
)

//  EXAMPLE in main.go
// --------------------
// err = client.OnTemplate(func(info string) {
// 		fmt.Println(info)
// 	}
// })
// if err != nil {
// 	panic(err)
// }

type templateListener func(info string)

// OnTemplate in order to receive Template changes.
func (c *Client) OnTemplate(listener templateListener) error {
	if _, found := c.events["OnTemplate"]; !found {

		// Subscribe to the template collection
		if err := c.ddpSubscribe(bbb.template, c.updateTemplate); err != nil {
			return err
		}
	}

	c.events["OnTemplate"] = append(c.events["OnStatus"], listener)
	return nil
}

// informs all listeners with the new info
func (c *Client) updateTemplate(collection string, operation string, id string, doc ddp.Update) {
	//Read data from doc
	info := convert.String(doc["info"], "")
	// Inform all listeners
	for _, event := range c.events["OnTemplate"] {
		// call event(info)
		f := reflect.TypeOf(event)
		if f.Kind() == reflect.Func { //is function
			if f.NumIn() == 1 && f.NumOut() == 0 { //inbound parameters == 1, outbound parameters == 0
				if f.In(0).Kind() == reflect.String { //parameter 0 is of type string (StatusType)
					go reflect.ValueOf(event).Call([]reflect.Value{reflect.ValueOf(info)}) // Call the function with the parameter info
				}
			}
		}
	}
}
