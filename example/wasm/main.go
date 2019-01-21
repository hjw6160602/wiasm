package main

import (
    "github.com/hjw6160602/wiasm"
    "strconv"
)

//func main() {
//   var output string
//   wiasm.Call("test", "hello", &output)
//
//   platform := platforms["a"]
//   eric := platform.users[ericID]
//   wiasm.Log(eric.repTracer.Rep.String())
//   LiquidateRepByRepID("a", []RepID{ericID}, 4)
//   wiasm.Log(eric.repTracer.Rep.String())
//}

//func main() {
//	var wg sync.WaitGroup
//	for i := 0; i < 5; i++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			var output string
//			wiasm.Call("myModule.say", "hello", &output)
//			wiasm.Log("worker result:" + output)
//		}()
//	}
//	wg.Wait()
//	wiasm.Log("wiasm program finished")
//}


type WasmFunc func(params interface{}) interface{}

var methodsMap = map[string] WasmFunc {
    "Add": Add,
    "Sub": Sub,
}

type Function struct  {
    MethodName string
    Params interface{}
    Result interface{}
}

func Add(a interface{}) interface{} {
    var result string
    switch a.(type) {
    case string:
        x, _ := strconv.Atoi(a.(string))
        result = strconv.Itoa(x+1)
    }
    return result
}

func Sub(a interface{}) interface{} {
    var result string
    switch a.(type) {
    case string:
        x, _ := strconv.Atoi(a.(string))
        result = strconv.Itoa(x-1)
    }
    return result
}

func main() {
    //fs := []string{"Add", "Sub"}
    channel := make(chan *Function)
    //for key, _ := range methodsMap {
    //    wiasm.Log("72 for loop:" + key)
    //    go func() {
    //        call := &Function{key,"",""}
    //        wiasm.Call(key, "", &call)
    //        channel <- call
    //    }()
    //}
    go func() {
       call := &Function{"Add","",""}
       wiasm.Call("Add", "", &call)
       channel <- call
    }()
    go func() {
       call := &Function{"Sub","",""}
       wiasm.Call("Sub", "", &call)
       channel <- call
    }()

    for i := 0; i < 100; i++ {
        select {
        case call := <-channel:
            exectueDecision(channel, call)
        }
    }
}


func resultCallBack(channel chan *Function, method string, result interface{})  {
    c := make(chan *Function)
    go func() {
        var callback *Function
        wiasm.Call(method + "Result", result, &callback)
        c <- callback

        var call *Function
        wiasm.Call(method, "", &call)
        channel <- call
        select {
        case call := <-channel:
            exectueDecision(channel, call)
        }
    }()
    select {
    case <- c:
        wiasm.Log("result return finished ")
    }
}

func exectueDecision(channel chan *Function, call *Function)  {
    for pattern, handleFunc := range (methodsMap) {
        if pattern == call.MethodName {
            resultCallBack(channel, pattern, handleFunc(call.Params))
        }
    }
}
