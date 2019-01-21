package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"github.com/perlin-network/life/exec"
	"github.com/hjw6160602/wiasm/resolv"
    "encoding/json"
)


type Function struct  {
    MethodName string
    Params interface{}
    Result interface{}
}

func readWasm() []byte {
    b, err := ioutil.ReadFile("example/app.wasm")
    if err != nil {
        log.Fatal("file read:", err)
    }
    return b
}

func setupVmAndResolv(b []byte) (r *resolv.Resolver, vm *exec.VirtualMachine) {
    r = resolv.New()
    vm, err := exec.NewVirtualMachine(b, exec.VMConfig{}, r, nil)
    if err != nil { // if the wasm bytecode is invalid
        log.Fatal("vm create:", err)
    }
    entryID, ok := vm.GetFunctionExport("run") // can change to whatever exported function name you want
    if !ok {
        panic("entry function not found")
    }
    _, err = vm.Run(entryID, 0, 0) // start vm
    if err != nil {
        vm.PrintStackTrace()
        log.Fatal("vm run:", err)
    }
    return r, vm
}

func generateInput(f Function) []byte {
    //Input := make(map[string] interface{})
    //Input["method"] = f.methodName
    //
    //Input["params"] = f.params
    //bytes, err:= json.Marshal(Input)
    bytes, err:= json.Marshal(f)
    if err != nil {
        fmt.Println(err.Error())
    }
    return bytes
}

func resumeCallFunc(vm *exec.VirtualMachine, r *resolv.Resolver, input []byte, index int)  {
    ret, err := r.Resume(vm, resolv.FCall{ // resume vm execution with callback result
        CB:     r.BlockedCalls[index].CB,
        Output: input,
    })

    if err != nil {
        vm.PrintStackTrace()
        log.Fatal("vm run:", err)
    }
    log.Printf("ret: %v, log:%v", ret, r.Stderr.String())
}

func resumeReturnCall(vm *exec.VirtualMachine, r *resolv.Resolver, methodName string) string {
    for index, v := range r.BlockedCalls {
        fmt.Println("resumeReturnCall:" + v.Method)
        if v.Method == methodName {
            result := string(r.BlockedCalls[index].Input)
            input, _ := json.Marshal(make(map[string] string))
            resumeCallFunc(vm, r, input, index)
            return result
        }
    }
    return ""
}

func callFunc(vm *exec.VirtualMachine, r *resolv.Resolver, f Function) string {
    if len(r.BlockedCalls) > 0 {
        isExist := false
        for index, v := range r.BlockedCalls {
            fmt.Println("callFunc:" + v.Method)
            if v.Method == f.MethodName {
                isExist = true
                input := generateInput(f)
                //input := []byte("\"22\"")
                fmt.Println("executing :"  + string(input) )
                resumeCallFunc(vm, r, input, index)
            }
        }

        return resumeReturnCall(vm, r, f.MethodName + "Result")

        if !isExist {
            fmt.Println("there's no function named " + f.MethodName + " exported in wasm file!")
        }
    }
    return ""
}

func callAdd(params string) Function {
    f := Function{"Add",params,""}
    return f
}

func callSub(params string) Function {
    //params := "11"
    f := Function{"Sub",params,""}
    return f
}

func main() {
    b := readWasm()
    r, vm := setupVmAndResolv(b)
    result1 := callFunc(vm, r, callAdd("10"))
    fmt.Println("result1:" + result1)
    //result2 := callFunc(vm, r, callSub("11"))
    //fmt.Println("result2:" + result2)
    //result3 := callFunc(vm, r, callSub("12"))
    //fmt.Println("result2:" + result3)
    //result4 := callFunc(vm, r, callSub("13"))
    //fmt.Println("result2:" + result4)
    //result5 := callFunc(vm, r, callSub("14"))
    //fmt.Println("result2:" + result5)
}
