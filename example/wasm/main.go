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
    results := make(chan map[string] string)
    //for _, method := range fs {
    //   go func() {
    //       //var out string
    //       var out map[string] string
    //       wiasm.Call(method, "input", &out)
    //       results <- out
    //   }()
    //}
    go func() {
       //var out string
       var out map[string] string
       wiasm.Call("Add", "input", &out)
       results <- out
    }()
    go func() {
       //var out string
       var out map[string] string
       wiasm.Call("Sub", "input", &out)
       results <- out
    }()
    //
    //len := len(fs)
    for i := 0; i < 2; i++ {
        select {
        case out := <-results:
            exectueDecision(out)
        }
    }
}

func exectueDecision(out map[string] string)  {
    method := out["method"]
    switch method {
    case "Add":
        params := out["params"]
        wiasm.Log("select out: " + params)
        pnum, _ := strconv.Atoi(params)
        result := Add(pnum)
        m := method + "Result"
        resultCallBack(m, result)
    case "Sub":
        params := out["params"]
        wiasm.Log("select out: " + params)
        pnum, _ := strconv.Atoi(params)
        result := Sub(pnum)
        m := method + "Result"
        resultCallBack(m, result)
    }
}

func resultCallBack(method string, result int)  {
    c := make(chan map[string] string)
    go func() {
        var output map[string] string
        wiasm.Call(method, strconv.Itoa(result), &output)
        c <- output
    }()
    select {
    case <- c:
        wiasm.Log("result return finished ")
    }
}