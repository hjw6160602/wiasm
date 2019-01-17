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

type Function struct  {
    methodName string
    params []string
    result string
}

func Add(a int) int {
    return a + 1
}

func Sub(a int) int {
    return a - 1
}

func main() {
    //fs := []string{"Add", "Sub"}
    channel := make(chan map[string] string)

    //for _, method := range fs {
    //    go func() {
    //        //var out string
    //        var call map[string] string
    //        wiasm.Call(method, "input", &call)
    //        results <- call
    //    }()
    //}

    go func() {
        //var out string
        var call map[string] string
        wiasm.Call("Add", "input", &call)
        channel <- call
    }()
    go func() {
        //var out string
        var call map[string] string
        wiasm.Call("Sub", "input", &call)
        channel <- call
    }()

    //
    //len := len(fs)
    for i := 0; i < 100; i++ {
        select {
        case call := <-channel:
            exectueDecision(channel, call)
        }
    }
}

func exectueDecision(channel chan map[string] string, call map[string] string)  {
    method := call["method"]
    switch method {
    case "Add":
        params := call["params"]
        wiasm.Log("select out: " + params)
        pnum, _ := strconv.Atoi(params)
        result := Add(pnum)
        resultCallBack(channel, method, result)
    case "Sub":
        params := call["params"]
        wiasm.Log("select out: " + params)
        pnum, _ := strconv.Atoi(params)
        result := Sub(pnum)
        resultCallBack(channel, method, result)
    }
}

func resultCallBack(channel chan map[string] string, method string, result int)  {
    c := make(chan map[string] string)
    go func() {
        var output map[string] string
        wiasm.Call(method + "Result", strconv.Itoa(result), &output)
        c <- output

        var call map[string] string
        wiasm.Call(method, "input", &call)
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